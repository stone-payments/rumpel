package rule

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/stone-payments/rumpel/proxy"
)

func (rls Rules) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	claim := newClaim(r.URL.EscapedPath(), r.Method, r.Header)
	rl, err := getMatchRuleByClaim(rls, claim)
	if err != nil {
		httpErr := &ErrRuleNotFound{StatusCode: http.StatusNotFound, Message: "not could match the rule"}
		responseJSON(w, httpErr.StatusCode, httpErr)
		return
	}

	pr := &proxy.Request{
		Method:  r.Method,
		URL:     rl.URL,
		Header:  r.Header,
		Body:    r.Body,
		Timeout: time.Duration(rl.Timeout),
	}
	err = proxy.Do(w, pr)
	if err != nil {
		httpErr := &ErrInternalProxy{StatusCode: http.StatusInternalServerError, Message: "internal proxy error", Reason: err.Error()}
		responseJSON(w, httpErr.StatusCode, httpErr)
	}
}

func response(w http.ResponseWriter, contentType string, statusCode int, body io.Reader) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	if _, err := io.Copy(w, body); err != nil {
		log.Printf("error when try response request, reason: %v", err)
	}
}

func responseJSON(w http.ResponseWriter, statusCode int, rawBody interface{}) {
	body, err := json.Marshal(rawBody)
	if err != nil {
		log.Printf("error when try serealize response body, reason: %v", err)
	}
	response(w, "application/json; charset=UTF-8", statusCode, bytes.NewBuffer(body))
}

// ErrRuleNotFound is an structure to response not found error
type ErrRuleNotFound struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func (ernf *ErrRuleNotFound) Error() string {
	return ernf.Message
}

// ErrInternalProxy is an structure to report error in proxy
type ErrInternalProxy struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Reason     string `json:"reason"`
}

func (eip *ErrInternalProxy) Error() string {
	return eip.Message
}

func getMatchRuleByClaim(rls Rules, claim *Claim) (*Rule, error) {
	for _, rl := range rls {
		if found := rl.MatchByClaim(claim); found {
			return &rl, nil
		}
	}
	return nil, &ErrRuleNotFound{}
}
