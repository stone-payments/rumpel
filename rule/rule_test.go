package rule

import (
	"net/http/httptest"

	check "gopkg.in/check.v1"
)

func (s *RuleSuite) TestContainsByHeaders(c *check.C) {
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
		c.Check(test.Expected, check.Equals, result)
	}
}

func (s *RuleSuite) TestMatchWithValidRule(c *check.C) {
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
			claim := newClaim(r.URL.EscapedPath(), r.Method, r.Header)
			if found := rl.MatchByClaim(claim); found {
				c.Assert(rl.Name, check.Equals, test.Expected)
			}
		}
	}
}

func (s *RuleSuite) TestMatchWithNotFoundRule(c *check.C) {
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
			claim := newClaim(r.URL.EscapedPath(), r.Method, nil)
			c.Check(rl.MatchByClaim(claim), check.Equals, false)
		}
	}
}
