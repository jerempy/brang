package client

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHeaderToStringForPrint(t *testing.T) {
	tests := map[string]struct {
		in   *http.Header
		want string
	}{
		"easy pass": {in: &http.Header{"Test": []string{"header"}}, want: "- Test: [header] -"},
		"block Auth Value": {
			in:   &http.Header{"Authorization": []string{"Bearer ABC-123"}},
			want: "- Authorization: [******] -",
		}}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := headerToStringForPrint(tc.in)
			if got != tc.want {
				t.Errorf("%v: got %v - want %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestBRErrors(t *testing.T) {
	br := NewBResponse()
	br.AddError(fmt.Errorf("testErr"))
	if len(br.errs) != 1 {
		t.Error("should have 1 err")
	}
	e := br.writeOutErrors()
	fmt.Println("EE", e)
	if e != "1: testErr\n" {
		t.Errorf("got wrong value: %v", e)
	}
}

func TestWrite(t *testing.T) {
	br := NewBResponse()
	want := []byte{23, 24, 25}
	br.Write(want)
	got := br.OutBody.Bytes()
	if string(got) != string(want) {
		t.Errorf("got %v - want %v", got, want)
	}
}
