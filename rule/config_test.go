package rule

import (
	"errors"
	"fmt"
	"os"

	check "gopkg.in/check.v1"
)

func (s *RuleSuite) TestReadRulesFileWithEmptyPath(c *check.C) {
	rls, err := Config("")

	c.Assert(rls, check.IsNil)
	c.Assert(err, check.FitsTypeOf, &ErrCannotReadRulesFile{})
}

func (s *RuleSuite) TestReadRulesFileNoExistsFile(c *check.C) {
	rls, err := Config("./test.yaml")

	c.Assert(rls, check.IsNil)
	c.Assert(err, check.FitsTypeOf, &ErrCannotReadRulesFile{})
}

func (s *RuleSuite) TestInvalidFormatRulesFile(c *check.C) {
	path := "./test.yaml"
	content := []byte(`routes test`)

	file, err := os.Create(path)
	c.Assert(err, check.IsNil)

	_, err = file.Write(content)
	c.Assert(err, check.IsNil)

	defer func(c *check.C, path string) {
		c.Assert(os.Remove(path), check.IsNil)
	}(c, path)

	rls, err := Config(path)
	c.Assert(rls, check.IsNil)
	c.Assert(err, check.FitsTypeOf, &ErrInvalidRulesFile{})
}

func (s *RuleSuite) TestReadRulesFileSuccessful(c *check.C) {
	path := "./test.yaml"
	content := []byte(`
---
- name: A
  target: a.stone.com.br/splits/publish
  claims:
  - path: /a
  - path: /ab
    method: PUT

- name: B
  target: b.stone.com.br/splits
  claims:
  - path: /b
    method: POST
    headers:
      X-Test: true
`)
	file, err := os.Create(path)
	c.Assert(err, check.IsNil)

	_, err = file.Write(content)
	c.Assert(err, check.IsNil)

	defer func(c *check.C, path string) {
		c.Assert(os.Remove(path), check.IsNil)
	}(c, path)

	result, err := Config(path)
	c.Assert(err, check.IsNil)
	c.Assert(result, check.HasLen, 2)
}

func (s *RuleSuite) TestErrorFromErrCannotReadRulesFile(c *check.C) {
	err := &ErrCannotReadRulesFile{Reason: errors.New("B")}
	expected := fmt.Sprintf("cannot read rules file, reason: %v", err.Reason)
	c.Check(err.Error(), check.Equals, expected)
}

func (s *RuleSuite) TestErrorFromErrInvalidRulesFile(c *check.C) {
	err := &ErrInvalidRulesFile{Reason: errors.New("B")}
	expected := fmt.Sprintf("invalid rule file, reason: %v", err.Reason)
	c.Check(err.Error(), check.Equals, expected)
}
