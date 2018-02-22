package rule

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/stone-payments/rumpel/proxy"
)

func (rls Rules) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	claim := newClaim(r.Host, r.URL.EscapedPath(), r.Method, r.Header)
	rl, err := getMatchRuleByClaim(rls, claim)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	u := rl.URL
	if !rl.AbsolutePath {
		u = fmt.Sprintf("%v%v", rl.URL, r.URL.Path)
	}

	pr := &proxy.Request{
		Method:  r.Method,
		URL:     u,
		Header:  r.Header,
		Body:    r.Body,
		Timeout: time.Duration(rl.Timeout),
	}
	if err = proxy.Do(w, pr); err != nil {
		if ce, ok := err.(*url.Error); ok {
			if ce.Timeout() {
				w.WriteHeader(http.StatusGatewayTimeout)
				return
			}
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

// ErrRuleNotFound is an structure to response not found error
type ErrRuleNotFound struct {
	message string
}

func (e *ErrRuleNotFound) Error() string {
	return e.message
}

func getMatchRuleByClaim(rls Rules, claim *Claim) (*Rule, error) {
	for _, rl := range rls {
		if found := rl.MatchByClaim(claim); found {
			return &rl, nil
		}
	}
	return nil, &ErrRuleNotFound{"rule not found"}
}
