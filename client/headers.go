package client

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// Auth holds the authentication data for building HTTP Auth Headers.
// Allowed AuthTypes are: Password, Bearer, Token, or empty string
// If AuthType is passed value of 'Basic' it will be handled as 'Password'
type Auth struct {
	AuthType string `yaml:"authType,omitempty" json:"authType,omitempty"`
	Token    string `yaml:"token,omitempty" json:"token,omitempty"`
	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
}

// Returns value for Authorization Header as "Basic username:pasword"
func (a *Auth) PasswordAuth() string {
	auth := a.Username + ":" + a.Password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// Returns value for Authorization Header as "Bearer|Token the-token-123"
func (a *Auth) TokenAuth() string {
	return a.AuthType + " " + a.Token
}

// Makes a *http.Header from data provided to attach to a http.Request
func (rset *RequestSet) BuildHeader() (*http.Header, error) {
	a, h := &Auth{}, &http.Header{}
	a.AuthType = rset.AuthType
	if rset.AuthType == "Basic" || rset.AuthType == "Password" {
		s := strings.Split(rset.Cred, ":")
		a.Username, a.Password = s[0], s[1]
	} else {
		a.Token = rset.Cred
	}
	if err := mapHeaderSliceToHeader(rset.HeaderSlice, h); err != nil {
		return nil, err
	}
	addDefaultHeaders(h)
	if h.Get("Authorization") != "" {
		return h, nil
	}
	switch a.AuthType {
	case "":
		break
	case "Bearer", "Token":
		h.Set("Authorization", a.TokenAuth())
	case "Password", "Basic":
		h.Set("Authorization", a.PasswordAuth())
	default:
		return nil, fmt.Errorf("wrong auth type. accepts: Password|Bearer|Token. was given: %v", a.AuthType)
	}
	return h, nil
}

func addDefaultHeaders(h *http.Header) {
	if h.Get("Content-Type") == "" {
		h.Add("Content-Type", "application/json; charset=UTF-8")
	}
	if h.Get("Accept") == "" {
		h.Add("Accept", "text/html,text/plain,application/json")
	}
}

// Takes a []string of "key:value" from -H flag and adds them to *http.Header
func mapHeaderSliceToHeader(slice []string, h *http.Header) error {
	for _, v := range slice {
		s := strings.Split(v, ":")
		if len(s) != 2 {
			return fmt.Errorf("headers need to be <key>:<value> - ex: Content-Type:application/json")
		}
		h.Add(s[0], s[1])
	}
	return nil
}
