package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/jerempy/brang/config"
)

type BResponse struct {
	*http.Response
	OutBody bytes.Buffer
	errs    []error
}

type BResponseWriter interface {
	WriteResponse()
}

type BResponseCapturer interface {
	CaptureResponse(*http.Response, error)
}

type BResponseHandler interface {
	BResponseWriter
	BResponseCapturer
}

func (br *BResponse) Write(p []byte) (n int, err error) {
	return br.OutBody.Write(p)
}

func (br *BResponse) StringResponseBody() string {
	_, err := io.Copy(&br.OutBody, br.Body)
	if err != nil {
		br.AddError(fmt.Errorf("error reading the body: %v", err))
	}
	return br.OutBody.String()
}

func NewBResponse() *BResponse {
	return &BResponse{&http.Response{}, bytes.Buffer{}, []error{}}
}

func (br *BResponse) AddError(e error) {
	br.errs = append(br.errs, e)
}

func (br *BResponse) writeOutErrors() string {
	var s string
	for i, e := range br.errs {
		s += fmt.Sprintf("%d: %v\n", i+1, e)
	}
	return s
}

func (br *BResponse) CaptureResponse(r *http.Response, e error) {
	br.Response = r
	if e != nil {
		br.AddError(e)
	}
}

func (br *BResponse) WriteResponse() {
	owr := config.OutputWriter()
	w := owr.Init()
	if w.Err != nil {
		fmt.Println(w.Err)
		return
	}
	if w.Fn != nil {
		defer w.Fn()
	}
	switch w.Format {
	case "raw":
		br.Header.Write(w.Writer)
		w.Writer.WriteString(br.StringResponseBody())
	case "basic":
		t, err := template.New("basic").Parse(basicTmpl)
		if err != nil {
			br.AddError(err)
			w.Writer.WriteString(br.writeOutErrors())
			break
		}
		t.Execute(w.Writer, br)
	default:
		t, err := template.New("pretty").Funcs(template.FuncMap{
			"headerToStringForPrint": headerToStringForPrint,
			"writeErrors":            br.writeOutErrors,
		}).Parse(prettyTmpl)
		if err != nil {
			br.AddError(err)
			w.Writer.WriteString(br.writeOutErrors())
			break
		}
		t.Execute(w.Writer, br)
	}
}

const (
	prettyTmpl = `---| Request: {{.Request.Method}} --- url={{.Request.URL}}
   | Request Header:  {{headerToStringForPrint .Request.Header}} |---
---| Response --- Status Code: {{.StatusCode}} |---
{{ .StringResponseBody }}
 ---| End Response |---
{{ $length := len .errs }} {{ if gt $length 0 }}
Errors:
{{ .writeOutErrors }}
{{ end }}
`
	basicTmpl = `Status Code: {{.StatusCode}}
{{.StringResponseBody}}	
`
)

// Make pretty the header for printing as part of request-response output
func headerToStringForPrint(h *http.Header) string {
	var s string
	for k, v := range *h {
		if k == "Authorization" {
			s += fmt.Sprintf(`- %s: [******] -`, k)
		} else {
			s += fmt.Sprintf(`- %s: %v -`, k, v)
		}
	}
	return s
}
