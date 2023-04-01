package client

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/jerempy/brang/config"
)

func mockRSet(url string) *RequestSet {
	return &RequestSet{
		Method:      "GET",
		URL:         url,
		AuthType:    "Bearer",
		Cred:        "ABC-456",
		HeaderSlice: []string{},
	}
}

func mockReq() *http.Request {
	want, _ := http.NewRequest("GET", "https://mysite.com/posts/1", http.NoBody)
	want.Header.Add("Accept", "text/html,text/plain,application/json")
	want.Header.Add("Content-Type", "application/json; charset=UTF-8")
	want.Header.Add("Authorization", "Bearer ABC-456")
	return want
}

func compareReqs(r1, r2 *http.Request) bool {
	if hok := compareHeaders(r1.Header, r2.Header); !hok {
		return false
	}
	return *r1.URL == *r2.URL && r1.Body == r2.Body && r1.Method == r2.Method
}

func TestBuildRequest(t *testing.T) {
	rset := mockRSet("https://mysite.com/posts/1")
	got, err := rset.BuildRequest()
	if err != nil {
		t.Errorf("got err: %v", err)
	}
	want := mockReq()
	if !compareReqs(got, want) {
		t.Errorf("got %v - want %v", got, want)
	}
}

func TestBufBodyBuilder(t *testing.T) {
	tests := map[string]struct {
		in   *RequestSet
		want string
	}{
		"no body":   {in: mockRSet("https://mysite.com/posts/1"), want: ""},
		"with body": {in: &RequestSet{Body: "test body"}, want: "test body"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			reader := tc.in.bufBodyBuilder()
			got := make([]byte, len(tc.want))
			reader.Read(got)
			if string(got) != tc.want {
				t.Errorf("got %v - want %v", got, tc.want)
			}
		})
	}
}

func TestRequestSetSend(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello")
	}))
	defer ts.Close()
	testDoReqYml := []byte(fmt.Sprintf(`
testspace:
  requests:
    site: %v
`, ts.URL))
	testConfigyml := []byte(fmt.Sprintf(`
outWriter: file
outWriterFormat: basic
deleteTempFileOnClose: false
outWriterFileName: testfile
outWriterFilePath: %v
`, os.TempDir()))
	config.Requests.SetConfigType("yaml")
	config.Brang.SetConfigType("yaml")
	config.Requests.ReadConfig(bytes.NewBuffer(testDoReqYml))
	config.Brang.ReadConfig(bytes.NewBuffer(testConfigyml))
	rset1 := mockRSet("testspace.site")
	rset2 := mockRSet(ts.URL)
	tests := map[string]struct {
		in     *RequestSet
		format string
		want   string
	}{
		"from cmd": {in: rset1, format: "basic", want: "Status Code: 200\n"},
		"loadsave": {in: rset2, format: "pretty", want: fmt.Sprintf("---| Request: GET --- url=%v\n", ts.URL)},
		"raw":      {in: rset1, format: "raw", want: "Content-Length: 6\r\n"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			config.Brang.Set("outWriterFormat", tc.format)
			tc.in.Send()
			fname := path.Join(os.TempDir(), "testfile.txt")
			f, err := os.Open(fname)
			if err != nil {
				t.Errorf("%v: %v", name, err)
			}
			s := bufio.NewReader(f)
			got, _ := s.ReadString('\n')
			defer os.Remove(fname)
			defer f.Close()
			if got != tc.want {
				t.Errorf("%v: got %v  - want %s", name, got, tc.want)
			}
		})
	}

}

func TestBodyFile(t *testing.T) {
	f, _ := os.CreateTemp("", "testBodyFile_*.json")
	b := []byte(`{"test": "value"}`)
	f.Write(b)
	r := mockRSet("mysite.site.test")
	r.BodyFile(f.Name())
	defer os.Remove(f.Name())
	defer f.Close()
	want := string(b)
	if r.Body != want {
		t.Errorf("testBodyFile: got %v - want %v", r.Body, want)
	}
}
