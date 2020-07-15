package request

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Request struct {
	Req      *http.Request
	Res      *http.Response
	Client   *http.Client
	ResBytes []byte
	Err      error
}

// New - initiates http request with specified method, url and body
func New(method, url string, body io.ReadCloser) *Request {
	sr := &Request{
		Client: &http.Client{Timeout: 10 * time.Second},
	}

	sr.Req, sr.Err = http.NewRequest(method, url, body)
	return sr
}

// AddHeaders - addsd common headers (content type, authorization)
func (sr *Request) AddHeaders(key, val string) *Request {
	if sr.Err != nil {
		return sr
	}

	sr.Req.Header.Set(key, val)
	return sr
}

// Do - calls client call with previously defined request
func (sr *Request) Do() *Request {
	if sr.Err != nil {
		return sr
	}

	sr.Res, sr.Err = sr.Client.Do(sr.Req)
	return sr
}

// Read - reads response body into slice of bytes
func (sr *Request) Read() *Request {
	if sr.Err != nil {
		return sr
	}

	sr.ResBytes, sr.Err = ioutil.ReadAll(sr.Res.Body)
	return sr
}

// Decode - decodes response body into given data structure
func (sr *Request) Decode(data interface{}) *Request {
	if sr.Err != nil {
		return sr
	}

	sr.Err = json.NewDecoder(sr.Res.Body).Decode(data)

	return sr
}
