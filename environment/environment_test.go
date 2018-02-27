package environment

import (
	"errors"
	"fmt"
	"os"
	"testing"

	check "gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type EnvironmentSuite struct{}

var _ = check.Suite(&EnvironmentSuite{})

func (s *EnvironmentSuite) TestReadEnvironmentByVariables(c *check.C) {
	c.Assert(os.Setenv(applicationPortName, "test-port"), check.IsNil)
	c.Assert(os.Setenv(rulesConfigPathName, "test-file"), check.IsNil)

	defer func(c *check.C) {
		c.Assert(os.Unsetenv(applicationPortName), check.IsNil)
		c.Assert(os.Unsetenv(rulesConfigPathName), check.IsNil)
	}(c)

	result, _ := Read(TestMode, []string{"rumpel"})
	expected := Environment{TestMode, "test-file", "test-port", false}

	c.Check(expected, check.DeepEquals, *result)
}

func (s *EnvironmentSuite) TestReadEnvironmentByFlags(c *check.C) {
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
		c.Check(test.Expected, check.DeepEquals, *result)
	}
}

func (s *EnvironmentSuite) TestReadEnvironmentWithEmptyName(c *check.C) {
	result, err := Read("", []string{"rumpel"})
	expected := DevelopmentMode
	c.Assert(err, check.IsNil)
	c.Check(expected, check.Equals, result.Name)
}

func (s *EnvironmentSuite) TestReadEnvironmentFailure(c *check.C) {
	_, err := Read("", []string{"rumpel", "-test.v", "testing"})
	c.Assert(err, check.NotNil)
}

func (s *EnvironmentSuite) TestErrorFromErrCannotReadEnvironment(c *check.C) {
	err := &ErrCannotReadEnvironment{Reason: errors.New("B")}
	expected := fmt.Sprintf("Cannot read environment, reason: %v", err.Reason)
	c.Check(err.Error(), check.Equals, expected)
}
