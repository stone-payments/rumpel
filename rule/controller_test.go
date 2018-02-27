package rule

import (
	"net/http"
	"net/http/httptest"

	check "gopkg.in/check.v1"
)

func (s *RuleSuite) TestControllerNotFoundRule(c *check.C) {
	rls := Rules{
		{Name: "A", URL: "a.com.br/publish/split_queue", Claims: []Claim{{Path: "/a"}, {Path: "/ab", Method: "PUT"}}},
		{Name: "C", URL: "c.com.br/publish/split_queue", Claims: []Claim{{Path: "/c", Headers: map[string]string{"X-Test": "true"}}}},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/x", nil)
	rls.Proxy(false).ServeHTTP(w, r)

	c.Check(w.Code, check.Equals, http.StatusNotFound)
}

func (s *RuleSuite) TestControllerProxy(c *check.C) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusGone)
	}))
	defer server.Close()

	rls := Rules{
		{Name: "A", URL: server.URL, Claims: []Claim{{Path: "/a"}, {Path: "/ab", Method: "PUT"}}},
		{Name: "C", URL: "c.com.br/publish/split_queue", Claims: []Claim{{Path: "/c", Headers: map[string]string{"X-Test": "true"}}}},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/a", nil)
	rls.Proxy(false).ServeHTTP(w, r)

	c.Check(w.Code, check.Equals, http.StatusGone)
}
