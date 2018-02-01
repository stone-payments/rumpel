package proxy

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Request is a structure for config your request proxy
type Request struct {
	Method  string
	URL     string
	Header  http.Header
	Body    io.Reader
	Timeout time.Duration
}

// Do is func to prepare and run the proxy request
func Do(w http.ResponseWriter, r *Request) error {
	u, err := url.Parse(r.URL)
	if err != nil {
		return err
	}
	nr := &http.Request{
		Method: r.Method,
		URL:    u,
		Body:   ioutil.NopCloser(r.Body),
		Header: r.Header,
	}
	c := http.Client{Timeout: r.Timeout * time.Millisecond}
	resp, err := c.Do(nr)
	if err != nil {
		return err
	}
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	return err
}

// copyHeader is an util function to copy headers from request to response
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
