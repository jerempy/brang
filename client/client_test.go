package client

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jerempy/brang/config"
)

func TestNewClient(t *testing.T) {
	got := NewClient()
	want := &brangClient{
		&http.Client{Timeout: time.Second * 10},
	}
	if got.Timeout != want.Timeout {
		t.Errorf("got %v - want %v", got, want)
	}
}

func TestCheckIfNotHttp(t *testing.T) {
	tests := map[string]struct {
		in   string
		want bool
	}{
		"true":    {in: "https://mysite.com", want: true},
		"true 2":  {in: "http://mysite.com", want: true},
		"true 3":  {in: "www.mysite.com", want: true},
		"false":   {in: "mysite.posts.all", want: false},
		"false 2": {in: "mysite.posts.1", want: false},
		"false 3": {in: "mysite.com", want: false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := isHttp(tc.in)
			if got != tc.want {
				t.Errorf("%v: got %v - want %v", tc.in, got, tc.want)
			}
		})
	}
}

type mockBResponse struct{ code int }

func (br *mockBResponse) CaptureResponse(r *http.Response, e error) {
	br.code = r.StatusCode
}

func (br *mockBResponse) WriteResponse() {}

func TestDoRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()
	testDoReqYml := []byte(fmt.Sprintf(`
testspace:
  requests:
    site: %v
`, ts.URL))
	config.Requests.SetConfigType("yaml")
	config.Requests.ReadConfig(bytes.NewBuffer(testDoReqYml))
	sreq, err := LoadSavedRequest(mockRSet("testspace.site"))
	if err != nil {
		t.Error(err)
	}
	cmdReq, err := mockRSet(ts.URL).BuildRequest()
	if err != nil {
		t.Error(err)
	}
	tests := map[string]struct {
		in   *http.Request
		want int
	}{
		"from cmd": {in: cmdReq, want: 200},
		"loadsave": {in: sreq, want: 200},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rw := &mockBResponse{}
			NewClient().DoRequest(tc.in, rw)
			if rw.code != tc.want {
				t.Errorf("got %d  - want %d", rw.code, tc.want)
			}
		})
	}

}
