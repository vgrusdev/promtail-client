package promtail

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
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
	Streams []*jsonStream     `json:"streams"`
}
// -------------------

type clientJson struct {
	config    *ClientConfig
	quit      chan struct{}
	entries   chan *promtailStream
	waitGroup sync.WaitGroup
	client    myHttpClient

}

func NewClientJson(conf ClientConfig) (Client, error) {
	client := clientJson {
		config:  &conf,
		quit:    make(chan struct{}),
		entries: make(chan *promtailStream, LOG_ENTRIES_CHAN_SIZE),
		client:  myHttpClient{},
	}

	client.waitGroup.Add(1)
	go client.run()

	return &client, nil
}

func (c *clientJson) Chan() chan<- *promtailStream {
	return c.entries
}

func (c *clientJson) Shutdown() {
	close(c.quit)
	c.waitGroup.Wait()
}

func (c *clientJson) run() {
	var batch []*promtailStream
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
				batch = []*promtailStream{}
				batchSize = 0
				maxWait.Reset(c.config.BatchWait)
			}
		case <-maxWait.C:
			if batchSize > 0 {
				c.send(batch)
				batch = []*promtailStream{}
				batchSize = 0
			}
			maxWait.Reset(c.config.BatchWait)
		}
	}
}

func (c *clientJson) send(batch []*promtailStream) {

	entries := []*jsonEntry{}
	streams := []*jsonStream{}

	
	for _, pStream := range batch {
		for pEntry := range pStream.Entries {
			jEntry := []jsonEntry { fmt.Sprint(*pEntry.Ts.UnixNano()), *pEntry.Line, }
			entries = append(entries, &jEntry)
		}
		jStream := jsonStream {
			Labels: pStream.Labels,
			Entries: entries,
		}
		streams = append(streams, &jStream)
	}

	msg := jsonStreams{streams}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("promtail.ClientJson: unable to marshal a JSON document: %s\n", err)
		return
	}
	log.Println(string(msg))
	/*
	resp, body, err := c.client.sendReq("POST", c.config.PushURL, "application/json", jsonMsg)
	if err != nil {
		log.Printf("promtail.ClientJson: unable to send an HTTP request: %s\n", err)
		return
	}

	if resp.StatusCode != 204 {
		log.Printf("promtail.ClientJson: Unexpected HTTP status code: %d, message: %s\n", resp.StatusCode, body)
		return
	}
	*/
}
