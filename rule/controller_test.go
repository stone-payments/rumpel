package rule

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestControllerNotFoundRule(t *testing.T) {
	rls := Rules{
		{Name: "A", URL: "a.com.br/publish/split_queue", Claims: []Claim{{Path: "/a"}, {Path: "/ab", Method: "PUT"}}},
		{Name: "C", URL: "c.com.br/publish/split_queue", Claims: []Claim{{Path: "/c", Headers: map[string]string{"X-Test": "true"}}}},
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/x", nil)
	rls.Proxy(false).ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestControllerProxy(t *testing.T) {
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

	assert.Equal(t, http.StatusGone, w.Code)
}
