package client

import (
	"net/http"
	"time"
)

type brangClient struct{ *http.Client }

// Returns a new *brangClient, which is a wrapper of *http.Client
func NewClient() *brangClient {
	return &brangClient{
		&http.Client{Timeout: time.Second * 10},
	}
}

// Runs http.Client.Do(*http.Request) and prints to console results
func (c *brangClient) DoRequest(r *http.Request, brw BResponseHandler) {
	brw.CaptureResponse(c.Do(r))
	brw.WriteResponse()
}
