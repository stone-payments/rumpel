package proxy

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxy(t *testing.T) {
	w := httptest.NewRecorder()

	header := http.Header{}
	header.Set("Content-Type", "application/json; charset=UTF-8")

	pr := &Request{
		Method: "POST",
		Header: header,
		Body:   bytes.NewBuffer([]byte(`{"test":"hello test"}`)),
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, pr.Method)
	}))
	defer server.Close()

	pr.URL = server.URL
	assert.Nil(t, Do(w, pr))
}

func TestCopyHeader(t *testing.T) {
	dst := http.Header{}
	src := http.Header{"Gnr": []string{"axl", "slash", "duff", "steve", "dizzy"}}

	copyHeader(dst, src)
	assert.Equal(t, src, dst)
}
