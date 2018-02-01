package rule

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// ErrCannotReadRulesFile to report error when cannot read a rules config file
type ErrCannotReadRulesFile struct {
	Reason error
}

func (e *ErrCannotReadRulesFile) Error() string {
	return fmt.Sprintf("cannot read rules file, reason: %v", e.Reason)
}

// ErrInvalidRulesFile to report error when the rules config file is invalid
type ErrInvalidRulesFile struct {
	Reason error
}

func (e *ErrInvalidRulesFile) Error() string {
	return fmt.Sprintf("invalid rule file, reason: %v", e.Reason)
}

// Config get rules from config file
func Config(path string) (Rules, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, &ErrCannotReadRulesFile{err}
	}

	rules := make([]Rule, 0)
	if err := yaml.Unmarshal(content, &rules); err != nil {
		return nil, &ErrInvalidRulesFile{err}
	}
	return rules, nil
}
