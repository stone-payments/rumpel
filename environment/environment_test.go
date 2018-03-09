package environment

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadEnvironmentByVariables(t *testing.T) {
	assert.Nil(t, os.Setenv(applicationPortName, "test-port"))
	assert.Nil(t, os.Setenv(rulesConfigPathName, "test-file"))

	defer func(t *testing.T) {
		assert.Nil(t, os.Unsetenv(applicationPortName))
		assert.Nil(t, os.Unsetenv(rulesConfigPathName))
	}(t)

	result, _ := Read(TestMode, []string{"rumpel"})
	expected := Environment{TestMode, "test-file", "test-port", false}

	assert.Equal(t, expected, *result)
}

func TestReadEnvironmentByFlags(t *testing.T) {
	cases := []struct {
		Name      string
		Arguments []string
		Expected  Environment
	}{
		{
			TestMode,
			[]string{"rumpel", "-port", ":4040"},
			Environment{
				Name:            TestMode,
				ApplicationPort: ":4040",
				RulesConfigPath: "./",
			},
		},
		{
			TestMode,
			[]string{"rumpel", "-rules", "path-test"},
			Environment{
				Name:            TestMode,
				ApplicationPort: ":28080",
				RulesConfigPath: "path-test",
			},
		},
		{
			TestMode,
			[]string{"rumpel", "-rules", "path-test", "-port", ":9080"},
			Environment{
				Name:            TestMode,
				ApplicationPort: ":9080",
				RulesConfigPath: "path-test",
			},
		},
	}

	for _, test := range cases {
		result, _ := Read(test.Name, test.Arguments)
		assert.Equal(t, test.Expected, *result)
	}
}

func TestReadEnvironmentWithEmptyName(t *testing.T) {
	result, err := Read("", []string{"rumpel"})
	expected := DevelopmentMode
	assert.NoError(t, err)
	assert.Equal(t, expected, result.Name)
}

func TestReadEnvironmentFailure(t *testing.T) {
	errTypeExpected := &ErrCannotReadEnvironment{}
	env, err := Read("", []string{"rumpel", "-test.v", "testing"})
	assert.Error(t, err)
	assert.IsType(t, errTypeExpected, err)
	assert.Nil(t, env)
}

func TestErrorFromErrCannotReadEnvironment(t *testing.T) {
	err := &ErrCannotReadEnvironment{Reason: errors.New("B")}
	expected := fmt.Sprintf("Cannot read environment, reason: %v", err.Reason)
	assert.Equal(t, err.Error(), expected)
}
