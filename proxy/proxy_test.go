package proxy

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	check "gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type ProxySuite struct{}

var _ = check.Suite(&ProxySuite{})

func (s *ProxySuite) TestProxy(c *check.C) {
	w := httptest.NewRecorder()

	header := http.Header{}
	header.Set("Content-Type", "application/json; charset=UTF-8")

	pr := &Request{
		Method: "POST",
		Header: header,
		Body:   bytes.NewBuffer([]byte(`{"test":"hello test"}`)),
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Method, check.Equals, pr.Method)
	}))
	defer server.Close()

	pr.URL = server.URL
	c.Assert(Do(w, pr), check.IsNil)
}

func (s *ProxySuite) TestCopyHeader(c *check.C) {
	dst := http.Header{}
	src := http.Header{"Gnr": []string{"axl", "slash", "duff", "steve", "dizzy"}}

	copyHeader(dst, src)
	c.Assert(src, check.DeepEquals, dst)
}
