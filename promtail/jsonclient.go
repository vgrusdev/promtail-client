package promtail

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
	"net/http"
)

// =================================================
// Logs format to be used for Marshal() to LOKI JSON api
// =================================================
// =======================================
// JSON msgs format that are sent to LOKI
// =======================================
//	See: https://grafana.com/docs/loki/latest/reference/loki-http-api/#ingest-logs
//	{
//		"streams": [
//		  {
//			"stream": {
//			  "<label>": "<value>"
//			},
//			"values": [
//				[ "<unix epoch in nanoseconds>", "<log line>" ],
//				[ "<unix epoch in nanoseconds>", "<log line>" ]
//			]
//		  }
//		]
//  }

type jsonEntry [2]string     // [ "<unix epoch in nanoseconds>", "<log line>" ]

type jsonStream struct {
	Labels  map[string]string  `json:"stream"`
	Entries []*jsonEntry       `json:"values"`
}

type jsonStreams struct {
	Streams []*jsonStream      `json:"streams"`
}
// -------------------

type clientJson struct {
	config    *ClientConfig
	quit      chan struct{}
	entries   chan *PromtailStream
	single    chan *SingleEntry
	waitGroup sync.WaitGroup
	client    myHttpClient

}

func NewClientJson(conf *ClientConfig) (Client, error) {
	n := conf.Name
	if n == "" {
		conf.Name = "unknown_name"
	}
	url, err := validateUrl(conf.PushURL)
	if err != nil {
		fmt.Printf("promtail.ClientJson: incorrect PushURL: %s\n", conf.PushURL)
		return nil, err
	}
	conf.PushURL = url

	if conf.TenantID == "" {
		conf.TenantID = "fake"
	}
	client := clientJson {
		config:  conf,
		quit:    make(chan struct{}),
		entries: make(chan *PromtailStream, LOG_ENTRIES_CHAN_SIZE),
		single:   make(chan *SingleEntry, LOG_ENTRIES_CHAN_SIZE),
		client:  myHttpClient{
					parent: http.Client {
						Timeout: conf.Timeout,
					},
		},
	}

	client.waitGroup.Add(1)
	go client.run()

	return &client, nil
}

func (c *clientJson) Chan() chan<- *PromtailStream {
	return c.entries
}

func (c *clientJson) Single() chan<- *SingleEntry {
	return c.single
}

func (c *clientJson) Shutdown() {
	close(c.quit)
	c.waitGroup.Wait()
}

func (c *clientJson) GetLocation() *time.Location {
	return c.config.Location
}

func (c *clientJson) run() {
	var batch []*PromtailStream
	batchSize := 0
	maxWait := time.NewTimer(c.config.BatchWait)

	defer func() {
		if batchSize > 0 {
			c.send(batch)
		}
		c.waitGroup.Done()
	}()

	for {
		select {
		case <-c.quit:
			return
		case entry := <-c.entries:
			batch = append(batch, entry)
			batchSize++
			if batchSize >= c.config.BatchEntriesNumber {
				c.send(batch)
				batch = []*PromtailStream{}
				batchSize = 0
				maxWait.Reset(c.config.BatchWait)
			}
		case sentry := <-c.single:
/*
		// Promtail common Logs entry format accepted by Chan() chan<- *PromtailStream
type PromtailEntry struct {
	Ts    time.Time
	Line  string
}
type PromtailStream struct {
	Labels  map[string]string
	Entries []*PromtailEntry
}
type SingleEntry struct {
	Labels  map[string]string
	Ts    	time.Time
	Line  	string
}
*/
			e := PromtailEntry {
				Ts: sentry.Ts,
				Line: sentry.Line,
			}
			s := PromtailStream {
				Labels: sentry.Labels,
				Entries: []*PromtailEntry { &e, },
			}
			batch = append(batch, &s)
			batchSize++
			if batchSize >= c.config.BatchEntriesNumber {
				c.send(batch)
				batch = []*PromtailStream{}
				batchSize = 0
				maxWait.Reset(c.config.BatchWait)
			}
		case <-maxWait.C:
			if batchSize > 0 {
				c.send(batch)
				batch = []*PromtailStream{}
				batchSize = 0
			}
			maxWait.Reset(c.config.BatchWait)
		}
	}
}

func (c *clientJson) send(batch []*PromtailStream) {

	//entries := []*jsonEntry{}
	streams := []*jsonStream{}

	for _, pStream := range batch {
		entries := []*jsonEntry{}
		for _, pEntry := range pStream.Entries {
			jEntry := jsonEntry { fmt.Sprint(pEntry.Ts.UnixNano()), pEntry.Line, }
			entries = append(entries, &jEntry)
		}
		pStream.Labels["service_name"] = c.config.Name
		jStream := jsonStream {
			Labels: pStream.Labels,
			Entries: entries,
		}
		streams = append(streams, &jStream)
	}

	msg := jsonStreams{
		Streams: streams,
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("promtail.ClientJson: unable to marshal a JSON document: %s\n", err)
		return
	}
	// Debug  fmt.Println(string(jsonMsg))

	resp, body, err := c.client.sendReq("POST", c.config.PushURL, "application/json", c.config.TenantID, jsonMsg)
	if err != nil {
		log.Printf("promtail.ClientJson: unable to send an HTTP request: %s\n", err)
		return
	}

	if resp.StatusCode != 204 {
		if resp.StatusCode == 400 {
			// Message too old
			// Debug log.Printf("promtail.ClientJson: %s'n", resp.StatusCode, body)
			return
		}
		log.Printf("promtail.ClientJson: Unexpected HTTP status code: %d, message: %s\n", resp.StatusCode, body)
		return
	}
}
