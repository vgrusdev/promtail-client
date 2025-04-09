package promtail

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)
const LOG_ENTRIES_CHAN_SIZE = 5000

type ClientConfig struct {
	Name               string           // Label name will be added to the stream, will be discovered as service_name in Loki
	PushURL            string			// E.g. http://localhost:3100/api/prom/push
	BatchWait          time.Duration	// Batch flush wait timeout
	BatchEntriesNumber int				// Batch buffer size
	Timeout            time.Duration    // HTTP Client Timeout !ToDo
	Location:          *time.Location	// time location provided by Stream, e.g. Europe/Moscow
}

type Client interface {
//	Debugf(format string, args ...interface{})
	Chan() chan<- *PromtailStream
	Single() chan<- *SingleEntry
	Shutdown()
}

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

// http.Client wrapper for adding new methods, particularly sendReq
type myHttpClient struct {
	parent http.Client
}

// A bit more convenient method for sending requests to the HTTP server
func (client *myHttpClient) sendReq(method, url string, ctype string, reqBody []byte) (resp *http.Response, resBody []byte, err error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", ctype)

	resp, err = client.parent.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	resBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return resp, resBody, nil
}