package rule

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/stone-payments/rumpel/logger"
	"github.com/stone-payments/rumpel/proxy"
)

// Proxy is handle for rule proxy
func (rls Rules) Proxy(verbose bool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

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
		if verbose {
			payload := logger.Payload{
				Application:  "rumpel",
				Method:       r.Method,
				Scheme:       r.Proto,
				Origin:       fmt.Sprintf("%v%v", r.Host, r.URL.Path),
				Target:       u,
				ResponseTime: fmt.Sprintf("%vs", fmt.Sprintf("%.3f", time.Since(start).Seconds())),
			}
			if err := logger.Do(payload); err != nil {
				return
			}
		}
	})

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
