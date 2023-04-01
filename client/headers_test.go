package client

import (
	"encoding/base64"
	"net/http"
	"testing"
)

func mockAuth() *Auth {
	return &Auth{
		Username: "joey",
		Password: "sandwiches",
		AuthType: "Password",
	}
}

func compareHeaders(h1, h2 http.Header) bool {
	if len(h1) != len(h2) {
		return false
	}
	for k, v := range h1 {
		v2, ok := h2[k]
		if !ok {
			return false
		}
		if v2[0] != v[0] {
			return false
		}
	}
	return true
}

func TestAuth(t *testing.T) {
	tests := map[string]struct {
		in   *Auth
		want string
	}{
		"basic":  {mockAuth(), "Basic " + base64.StdEncoding.EncodeToString([]byte("joey:sandwiches"))},
		"token":  {&Auth{AuthType: "Token", Token: "test-token"}, "Token test-token"},
		"bearer": {&Auth{AuthType: "Bearer", Token: "test-token"}, "Bearer test-token"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var got string
			if name != "basic" {
				got = tc.in.TokenAuth()
			} else {
				got = tc.in.PasswordAuth()
			}
			if got != tc.want {
				t.Errorf("%v: got %v - want %v", name, got, tc.want)
			}
		})
	}
}

func TestMapHeaderSliceToHeader(t *testing.T) {
	h := http.Header{}
	s := []string{"key:value", "test:value", "foo:roo"}
	want := http.Header{
		"Key":  []string{"value"},
		"Test": []string{"value"},
		"Foo":  []string{"roo"},
	}
	mapHeaderSliceToHeader(s, &h)
	if !compareHeaders(h, want) {
		t.Error("header didn't get mapped")
	}
}

func TestAddDefaultHeaders(t *testing.T) {
	h := http.Header{}
	addDefaultHeaders(&h)
	want := http.Header{
		"Content-Type": []string{"application/json; charset=UTF-8"},
		"Accept":       []string{"text/html,text/plain,application/json"},
	}
	if !compareHeaders(h, want) {
		t.Error("header didn't get mapped")
	}
}

func TestRequestSetBuildHeader(t *testing.T) {
	rset := mockRSet("https://mysite.com/posts/1")
	got, err := rset.BuildHeader()
	if err != nil {
		t.Error(err)
	}
	rsetB := mockRSet("https://mysite.com/posts/1")
	rsetB.AuthType = "Basic"
	rsetB.Cred = "joey:sandwiches"
	gotB, errB := rsetB.BuildHeader()
	if errB != nil {
		t.Error(errB)
	}
	want := http.Header{
		"Content-Type":  []string{"application/json; charset=UTF-8"},
		"Accept":        []string{"text/html,text/plain,application/json"},
		"Authorization": []string{"Bearer ABC-456"},
	}
	wantB := http.Header{
		"Content-Type":  []string{"application/json; charset=UTF-8"},
		"Accept":        []string{"text/html,text/plain,application/json"},
		"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte("joey:sandwiches"))},
	}
	if !compareHeaders(*got, want) {
		t.Errorf("err build header. got: %v - want: %v", got, want)
	}
	if !compareHeaders(*gotB, wantB) {
		t.Errorf("err build header. got: %v - want: %v", gotB, wantB)
	}
}
