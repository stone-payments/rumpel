package environment

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Constants for work with alias for environment variables
const (
	rulesConfigPathName = "RUMPEL_RULES_FILE"
	applicationPortName = "RUMPEL_APPLICATION_PORT"
)

// Alias for environments mode name
const (
	ProductionMode  = "PRODUCTION"
	DevelopmentMode = "DEVELOPMENT"
	TestMode        = "TEST"
)

// Environment as recipient for all environment parameters
type Environment struct {
	Name            string
	RulesConfigPath string
	ApplicationPort string
}

// ErrCannotReadEnvironment to report error when cannot read environment
type ErrCannotReadEnvironment struct {
	Reason error
}

func (e *ErrCannotReadEnvironment) Error() string {
	return fmt.Sprintf("Cannot read environment, reason: %v", e.Reason)
}

// Read environments
func Read(name string, args []string) (*Environment, error) {
	env := &Environment{Name: strings.ToUpper(name)}
	if env.Name == "" {
		env.Name = DevelopmentMode
	}

	cmd := flag.NewFlagSet(args[0], flag.ContinueOnError)

	env.RulesConfigPath = os.Getenv(rulesConfigPathName)
	if env.RulesConfigPath == "" {
		cmd.StringVar(&env.RulesConfigPath, "rules", ".rumpel.yaml", "parameter for rule configurations")
	}

	env.ApplicationPort = os.Getenv(applicationPortName)
	if env.ApplicationPort == "" {
		cmd.StringVar(&env.ApplicationPort, "port", ":28080", "parameter for set application port")
	}

	if err := cmd.Parse(args[1:]); err != nil {
		return nil, &ErrCannotReadEnvironment{err}
	}
	return env, nil
}
