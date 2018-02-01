package rule

import (
	"strings"
)

// Headers alias type for map[string]string
type Headers map[string]string

// Body is body structure for rules
type Body struct {
	Type   string            `yaml:"type"`
	Scheme map[string]string `yaml:"scheme"`
}

// Claim is claim structure for rules
type Claim struct {
	Path    string  `yaml:"path"`
	Method  string  `yaml:"method"`
	Headers Headers `yaml:"headers"`
}

// ContainsByHeaders verify if contains A in B as (map[string]string)
func (c *Claim) ContainsByHeaders(hs Headers) bool {
	for key, hvalue := range c.Headers {
		if value := hs[key]; value != hvalue {
			return false
		}
	}
	return true
}

func newClaim(path, method string, headers map[string][]string) *Claim {
	cHeaders := make(map[string]string)
	for key, value := range headers {
		cHeaders[key] = strings.Join(value, " ")
	}
	return &Claim{
		Path:    path,
		Method:  method,
		Headers: cHeaders,
	}
}

// Rule is structure that contains parameters to proxy
type Rule struct {
	Name    string  `yaml:"name"`
	URL     string  `yaml:"url"`
	Body    Body    `yaml:"body"`
	Timeout int64   `yaml:"timeout"`
	Claims  []Claim `yaml:"claims"`
}

// Rules is alias to []Rule
type Rules []Rule

// MatchByClaim match by rule claim
func (r Rule) MatchByClaim(c *Claim) bool {
	for _, claim := range r.Claims {
		if claim.Path != "" {
			if claim.Path != c.Path {
				continue
			}
		}
		if claim.Method != "" {
			if strings.ToUpper(claim.Method) != strings.ToUpper(c.Method) {
				continue
			}
		}
		if len(claim.Headers) > 0 {
			if !claim.ContainsByHeaders(c.Headers) {
				continue
			}
		}
		return true
	}
	return false
}
