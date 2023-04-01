package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

type RequestSet struct {
	URL,
	Method,
	AuthType,
	Cred,
	Params,
	Body string
	HeaderSlice []string
}

// Creates new *http.Request and attaches a *http.Header
func (rset *RequestSet) BuildRequest() (*http.Request, error) {
	rset.URL += rset.Params
	body := rset.bufBodyBuilder()
	req, err := http.NewRequest(rset.Method, rset.URL, body)
	if err != nil {
		return nil, err
	}
	header, err := rset.BuildHeader()
	if err != nil {
		return nil, fmt.Errorf("err building header: %w", err)
	}
	req.Header = *header
	return req, nil
}

func (rset *RequestSet) bufBodyBuilder() io.Reader {
	if len(rset.Body) == 0 {
		return http.NoBody
	}
	return bytes.NewBuffer([]byte(rset.Body))
}

func (rset *RequestSet) Send() {
	var req *http.Request
	if !isHttp(rset.URL) {
		r, err := LoadSavedRequest(rset)
		if err != nil {
			fmt.Println(err)
			return
		}
		req = r
	} else {
		r, err := rset.BuildRequest()
		if err != nil {
			fmt.Println("error building request from data provided: ", err)
			return
		}
		req = r
	}
	NewClient().DoRequest(req, NewBResponse())
}

// Takes path to a file with body for a request. Reads it and attaches to request.
func (rset *RequestSet) BodyFile(f string) error {
	b, err := os.ReadFile(f)
	if err != nil {
		return fmt.Errorf("err reading file: %w", err)
	}
	rset.Body = string(b)
	return nil
}

func isHttp(u string) bool {
	startsHttpOrWww, _ := regexp.Compile(`^(?:https?:\/\/|www\.)`)
	return startsHttpOrWww.MatchString(u)
}
