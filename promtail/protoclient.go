package promtail

import (
	//"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	
	"github.com/golang/snappy"
	"github.com/vgrusdev/promtail-client/logproto"
	"log"
	"sync"
	"time"
	//"encoding/json"
	"net/http"
)
// https://habr.com/ru/companies/otus/articles/784732/

type clientProto struct {
	config    *ClientConfig
	quit      chan struct{}
	entries   chan *PromtailStream
	single    chan *SingleEntry
	waitGroup sync.WaitGroup
	client    myHttpClient
}

func NewClientProto(conf *ClientConfig) (Client, error) {
	n := conf.Name
	if n == "" {
		conf.Name = "unknown_name"
	}
	client := clientProto{
		config:  conf,
		quit:    make(chan struct{}),
		entries: make(chan *PromtailStream, LOG_ENTRIES_CHAN_SIZE),
		single:  make(chan *SingleEntry, LOG_ENTRIES_CHAN_SIZE),
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

func (c *clientProto) Chan() chan<- *PromtailStream {
	return c.entries
}

func (c *clientProto) Single() chan<- *SingleEntry {
	return c.single
}

func (c *clientProto) Shutdown() {
	close(c.quit)
	c.waitGroup.Wait()
}

func (c *clientProto) run() {
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
			e := PromtailEntry {
				Ts:     sentry.Ts,
				Line:   sentry.Line,
			}
			s := PromtailStream {
				Labels:  sentry.Labels,
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

func (c *clientProto) send(batch []*PromtailStream) {

//	entries := []*logproto.Entry{}
/* from logproto.pb.go
	type Entry struct {
		Timestamp            *timestamp.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
		Line                 string               `protobuf:"bytes,2,opt,name=line,proto3" json:"line,omitempty"`
		XXX_NoUnkeyedLiteral struct{}             `json:"-"`
		XXX_unrecognized     []byte               `json:"-"`
		XXX_sizecache        int32                `json:"-"`
	}
	type Stream struct {
		Labels               string   `protobuf:"bytes,1,opt,name=labels,proto3" json:"labels,omitempty"`
		Entries              []*Entry             `protobuf:"bytes,2,rep,name=entries,proto3" json:"entries,omitempty"`
		XXX_NoUnkeyedLiteral struct{}             `json:"-"`
		XXX_unrecognized     []byte               `json:"-"`
		XXX_sizecache        int32                `json:"-"`
	}
	type PushRequest struct {
	Streams              []*Stream `protobuf:"bytes,1,rep,name=streams,proto3" json:"streams,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
	}
*/
	streams := []*logproto.Stream{}
	for _, pStream := range batch {
		entries := []*logproto.Entry{}
		for _, pEntry := range pStream.Entries {
			tNano := pEntry.Ts.UnixNano()
			protoEntry := logproto.Entry { 
							Timestamp: &timestamp.Timestamp {
											Seconds: tNano / int64(time.Second),
											Nanos:   int32(tNano % int64(time.Second)),
										},
							Line:      pEntry.Line, 
						}
			entries = append(entries, &protoEntry)
		}

		labels := mapToLabels(pStream.Labels)
		protoStream := logproto.Stream {
			//Labels: pStream.Labels,
			Labels: labels,
			Entries: entries,
		}
		streams = append(streams, &protoStream)
	}
	req := logproto.PushRequest{
		Streams: streams,
	}

	buf, err := proto.Marshal(&req)
	if err != nil {
		log.Printf("promtail.ClientProto: unable to marshal: %s\n", err)
		return
	}

	buf = snappy.Encode(nil, buf)

	resp, body, err := c.client.sendReq("POST", c.config.PushURL, "application/x-protobuf", buf)
	if err != nil {
		log.Printf("promtail.ClientProto: unable to send an HTTP request: %s\n", err)
		return
	}
	
	if resp.StatusCode != 204 {
		log.Printf("promtail.ClientProto: Unexpected HTTP status code: %d, message: %s\n", resp.StatusCode, body)
		return
	}
}
func mapToLabels(m  map[string]string) string {

	s := "{"
	f := false
	for k, v := range m {
		if f == true { 
			s = s + ","
		} else {
			f = true
		}
        s = s + k + "=\"" + v + "\""
    }
	s = s + "}"
	return s
}