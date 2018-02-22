package rule

import (
	"strings"
)

// Headers alias type for map[string]string
type Headers map[string]string

// Claim is claim structure for rules
type Claim struct {
	Host    string  `yaml:"host"`
	Path    string  `yaml:"path"`
	Method  string  `yaml:"method"`
	Headers Headers `yaml:"headers"`
}

// Rule is structure that contains parameters to proxy
type Rule struct {
	Name         string  `yaml:"name"`
	URL          string  `yaml:"url"`
	AbsolutePath bool    `yaml:"absolute_path"`
	Timeout      int64   `yaml:"timeout"`
	Claims       []Claim `yaml:"claims"`
}

// Rules is alias to []Rule
type Rules []Rule

// ContainsByHeaders verify if contains A in B as (map[string]string)
func (c *Claim) ContainsByHeaders(hs Headers) bool {
	for key, hvalue := range c.Headers {
		if value := hs[key]; value != hvalue {
			return false
		}
	}
	return true
}

func newClaim(host, path, method string, headers map[string][]string) *Claim {
	cHeaders := make(map[string]string)
	for key, value := range headers {
		cHeaders[key] = strings.Join(value, " ")
	}
	return &Claim{
		Host:    host,
		Path:    path,
		Method:  method,
		Headers: cHeaders,
	}
}

// MatchByClaim match by rule claim
func (r Rule) MatchByClaim(c *Claim) bool {
	for _, claim := range r.Claims {
		if claim.Host != "" {
			if claim.Host != c.Host {
				continue
			}
		}
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
