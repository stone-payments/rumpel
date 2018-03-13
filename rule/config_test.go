package rule

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRulesFileWithEmptyPath(t *testing.T) {
	rls, err := Config("")
	assert.Nil(t, rls)
	assert.IsType(t, &ErrCannotReadRulesFile{}, err)
}

func TestReadRulesFileNoExistsFile(t *testing.T) {
	rls, err := Config("./test.yaml")
	assert.Nil(t, rls)
	assert.IsType(t, &ErrCannotReadRulesFile{}, err)
}

func TestInvalidFormatRulesFile(t *testing.T) {
	path := "./"
	content := []byte(`routes test`)

	pathFile := fmt.Sprintf("%v/rules.yaml", path)

	file, err := os.Create(pathFile)
	assert.Nil(t, err)

	_, err = file.Write(content)
	assert.Nil(t, err)

	defer func(t *testing.T, pathFile string) {
		assert.Nil(t, os.Remove(pathFile))
	}(t, pathFile)

	rls, err := Config(path)
	assert.Nil(t, rls)
	assert.IsType(t, &ErrInvalidRulesFile{}, err)
}

func TestReadRulesFileSuccessful(t *testing.T) {
	path := "./"
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

	pathFile := fmt.Sprintf("%v/rules.yaml", path)

	file, err := os.Create(pathFile)
	assert.Nil(t, err)

	_, err = file.Write(content)
	assert.Nil(t, err)

	defer func(t *testing.T, pathFile string) {
		assert.Nil(t, os.Remove(pathFile))
	}(t, pathFile)

	result, err := Config(path)
	assert.Nil(t, err)
	assert.Len(t, result, 2)
}

func TestErrorFromErrCannotReadRulesFile(t *testing.T) {
	err := &ErrCannotReadRulesFile{Reason: errors.New("B")}
	expected := fmt.Sprintf("cannot read rules file, reason: %v", err.Reason)
	assert.Equal(t, err.Error(), expected)
}

func TestErrorFromErrInvalidRulesFile(t *testing.T) {
	err := &ErrInvalidRulesFile{Reason: errors.New("B")}
	expected := fmt.Sprintf("invalid rule file, reason: %v", err.Reason)
	assert.Equal(t, err.Error(), expected)
}
