package rule

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsByHeaders(t *testing.T) {
	cases := []struct {
		Claim     *Claim
		Parameter Headers
		Expected  bool
	}{
		{&Claim{Headers: Headers{"X-Test": "true"}}, Headers{}, false},
		{&Claim{Headers: Headers{"X-Test": "true"}}, Headers{"X-Test": "true"}, true},
		{&Claim{Headers: Headers{"X-Test": "true"}}, Headers{"X-Test": "true", "X-Test-Trace": "none"}, true},
		{&Claim{Headers: Headers{"X-Test": "true"}}, Headers{"X-Test-Trace": "none", "X-Test": "true"}, true},
		{&Claim{Headers: Headers{}}, Headers{}, true},
	}
	for _, test := range cases {
		result := test.Claim.ContainsByHeaders(test.Parameter)
		assert.Equal(t, test.Expected, result)
	}
}

func TestMatchWithValidRule(t *testing.T) {
	rls := Rules{
		{Name: "A", URL: "a.com.br/publish/split_queue", Claims: []Claim{{Path: "/a"}, {Path: "/ab", Method: "PUT"}}},
		{Name: "C", URL: "c.com.br/publish/split_queue", Claims: []Claim{{Path: "/c", Headers: map[string]string{"X-Test": "true"}}}},
		{Name: "B", URL: "b.com.br/splits", Claims: []Claim{}},
	}

	cases := []struct {
		Claim    Claim
		Expected string
	}{
		{Claim{Path: "/d", Method: "GET", Headers: nil}, "B"},
		{Claim{Path: "/a", Method: "PUT", Headers: nil}, "A"},
		{Claim{Path: "/ab", Method: "PUT", Headers: nil}, "A"},
		{Claim{Path: "/ab", Method: "DELETE", Headers: nil}, "B"},
		{Claim{Path: "/c", Method: "PATCH", Headers: nil}, "B"},
		{Claim{Path: "/c", Method: "PATCH", Headers: map[string]string{"X-Test": "true"}}, "C"},
		{Claim{Path: "/c", Method: "PUT", Headers: map[string]string{"X-Test": "true"}}, "C"},
	}
	for _, test := range cases {
		r := httptest.NewRequest(test.Claim.Method, test.Claim.Path, nil)
		if test.Claim.Headers != nil {
			for key, value := range test.Claim.Headers {
				r.Header.Add(key, value)
			}
		}
		for _, rl := range rls {
			claim := newClaim(r.Host, r.URL.EscapedPath(), r.Method, r.Header)
			if found := rl.MatchByClaim(claim); found {
				assert.Equal(t, test.Expected, rl.Name)
			}
		}
	}
}

func TestMatchWithNotFoundRule(t *testing.T) {
	rules := Rules{
		{Name: "A", URL: "a.com.br/publish/split_queue", Claims: []Claim{{Path: "/a"}, {Path: "/ab", Method: "PUT"}}},
		{Name: "C", URL: "c.com.br/publish/split_queue", Claims: []Claim{{Path: "/c", Headers: map[string]string{"X-Test": "true"}}}},
	}

	cases := []struct {
		Claim    Claim
		Expected bool
	}{
		{Claim{Path: "/c", Method: "PATCH", Headers: nil}, false},
	}
	for _, test := range cases {
		r := httptest.NewRequest(test.Claim.Method, test.Claim.Path, nil)
		if test.Claim.Headers != nil {
			for key, value := range test.Claim.Headers {
				r.Header.Add(key, value)
			}
		}
		for _, rl := range rules {
			claim := newClaim(r.Host, r.URL.EscapedPath(), r.Method, nil)
			assert.Equal(t, false, rl.MatchByClaim(claim))
		}
	}
}
