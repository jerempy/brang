package client

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/jerempy/brang/config"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type SavedRequestSet struct {
	URL    string
	Body   string
	Header map[string]string
}

// LoadSavedRequest accepts a string from arg in running command as dot.notation.
// Looks up against the requests.yaml, loads it and searches for the request.
// The saved request in requests.yaml could be <name>: <url string>.
// Can also be <name>: {url: <string>, body: <json>, header: {<key>:<value>}}
// Currently only supports 1 request file and reads whole file - this can be re-visited.
func LoadSavedRequest(rset *RequestSet) (*http.Request, error) {
	if err := config.LoadRequests(); err != nil {
		return nil, err
	}
	s := strings.Split(rset.URL, ".")
	q := s[0] + ".requests." + strings.Join(s[1:], ".")
	auth := &Auth{}
	config.Requests.UnmarshalKey(s[0]+".auth", auth)
	rset.AuthType = auth.AuthType
	rset.Cred = getCred(auth)
	v := config.Requests.Get(q)
	if v == "" || v == nil {
		return nil, fmt.Errorf("saved request not found: %s. type 'brang config -h' for help in checking config", rset.URL)
	}
	var sr SavedRequestSet
	r, ok := v.(string)
	if ok {
		sr.URL = r
	} else {
		mapstructure.Decode(v, &sr)
		mapLoadedValsToHeaderSlice(rset, sr.Header)
	}
	rset.URL = sr.URL
	if rset.Body == "" {
		rset.Body = sr.Body
	}
	req, err := rset.BuildRequest()
	if err != nil {
		return nil, fmt.Errorf("err building saved requests: %w", err)
	}
	return req, nil
}

// Converts lower-case:value to Title-Case:value and sets on header
func mapLoadedValsToHeaderSlice(r *RequestSet, sMap map[string]string) {
	if len(sMap) == 0 {
		return
	}
	for k, v := range sMap {
		r.HeaderSlice = append(r.HeaderSlice, toTitle(k)+":"+v)
	}
}

func toTitle(s string) string {
	return cases.Title(language.Und, cases.NoLower).String(s)
}

// Gets correct cred value to assign to *RequestSet
func getCred(a *Auth) string {
	var s string
	switch a.AuthType {
	case "Bearer", "Token":
		s = checkEnv(a.Token)
	case "Password", "Basic":
		s = checkEnv(a.Username) + ":" + checkEnv(a.Password)
	default:
	}
	return s
}

func checkEnv(s string) string {
	s, isEnv := strings.CutPrefix(s, "$")
	if !isEnv {
		return s
	}
	e := os.Getenv(s)
	if e == "" {
		fmt.Printf("Couldn't find env variable for: %s\n", s)
	}
	return e
}
