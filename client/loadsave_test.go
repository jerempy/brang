package client

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"github.com/jerempy/brang/config"
)

var mockYml = []byte(`
testspace:
  auth:
    authtype: Bearer
    token: ABC-456
  requests:
    posts:
      1: https://mysite.com/posts/1

testspace2:
  auth:
    authtype: Password
    username: joe
    password: $TESTLOADSECRETTEST
  requests:
    users: https://mysite.com/users/	  
`)

func TestLoad(t *testing.T) {
	tests := map[string]struct {
		in      *RequestSet
		want    *http.Request
		wantErr error
	}{
		"success token":    {in: mockRSet("testspace.posts.1"), want: mockReq(), wantErr: nil},
		"success password": {in: mockRSet("testspace2.users"), want: mockReq(), wantErr: nil},
		"fail":             {in: mockRSet("not.a.real.request"), want: nil},
	}
	config.Requests.SetConfigType("yaml")
	config.Requests.ReadConfig(bytes.NewBuffer(mockYml))
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := LoadSavedRequest(tc.in)
			if err != tc.wantErr && got != tc.want {
				t.Errorf("%v: got %v - want %v", name, tc.in, tc.want)
			}
		})
	}
}

func TestMapLoadedValsToHeader(t *testing.T) {
	r, sMap := mockRSet("testspace.posts.1"), map[string]string{
		"test":  "fest",
		"test2": "festroo",
	}
	mapLoadedValsToHeaderSlice(r, sMap)
	if r.HeaderSlice[0] != "Test2:festroo" && r.HeaderSlice[0] != "Test:fest" {
		t.Errorf("values didn't get mapped to header %v, %v", r.HeaderSlice[0], r.HeaderSlice[1])
	}
}

func TestCheckEnv(t *testing.T) {
	want := "password"
	os.Setenv("TESTLOADSECRETTEST", want)
	defer os.Unsetenv("TESTLOADSECRETTEST")
	tests := map[string]struct {
		in   string
		want string
	}{
		"find env":      {"$TESTLOADSECRETTEST", "password"},
		"no env":        {"password", "password"},
		"cant find env": {"$WILLNOTFINDTHIS", ""},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := checkEnv(tc.in)
			if got != tc.want {
				t.Errorf("got %v - want %v", got, tc.want)
			}
		})
	}

}
