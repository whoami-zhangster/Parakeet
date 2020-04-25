package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	Endpoint struct {
		Response Response `yaml:"response"`
		Port     int      `yaml:"port"`
		Path     string   `yaml:"path"`
	}

	Response struct {
		Status     string            `yaml:"status"`     // e.g. "200 OK"
		StatusCode int               `yaml:"statusCode"` // e.g. 200
		Header     map[string]string `yaml:"headers"`    // kv rep
		Body       string            `yaml:"body"`
	}
)

// CreateHandleFunc Creates a http.Handler function from an endpoint definition
func (ep *Endpoint) CreateHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ep.copyValues(w.Header())
		w.WriteHeader(ep.Response.StatusCode)
		fmt.Fprintf(w, ep.Response.Body)
	}
}

func (ep *Endpoint) copyValues(header http.Header) {
	for k, v := range ep.Response.Header {
		header.Set(k, v)
	}
}

func stringIsJSON(s string) bool {
	var js string
	return json.Unmarshal([]byte(s), &js) == nil

}
