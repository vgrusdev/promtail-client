package promtail

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	//"google.golang.org/protobuf/types/known/timestamppb"
	"github.com/golang/snappy"
	"github.com/vgrusdev/promtail-client/logproto"
	"log"
	"sync"
	"time"
	"encoding/json"
)

type clientProto struct {
	config    *ClientConfig
	quit      chan struct{}
	entries   chan *promtailStream
	waitGroup sync.WaitGroup
	client    myHttpClient
}

func NewClientProto(conf ClientConfig) (Client, error) {
	client := clientProto{
		config:  &conf,
		quit:    make(chan struct{}),
		entries: make(chan *promtailStream, LOG_ENTRIES_CHAN_SIZE),
		client:  myHttpClient{},
	}

	client.waitGroup.Add(1)
	go client.run()

	return &client, nil
}

func (c *clientProto) Chan() chan<- *promtailStream {
	return c.entries
}

func (c *clientProto) Shutdown() {
	close(c.quit)
	c.waitGroup.Wait()
}

func (c *clientProto) run() {
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

func (c *clientProto) send(batch []*promtailStream) {

	entries := []*logproto.Entry{}
	streams := []*logproto.Stream{}

	for _, pStream := range batch {
		for _, pEntry := range pStream.Entries {
			tNano := pEntry.Ts.UnixNano()
			protoEntry := logproto.Entry { 
				Timestamp: &timestamp.Timestamp {
						Seconds: tNano / int64(time.Second),
						Nanos:   int32(tNano % int64(time.Second)),
				}
				Line:      pEntry.Line, 
			}
			entries = append(entries, &protoEntry)
		}
		jsonLabels, err := json.Marshal(pStream.Labels)
		if err != nil {
			fmt.Println(err)
			continue
		}
		protoStream := logproto.Stream {
			//Labels: pStream.Labels,
			Labels: string(jsonLabels),
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
