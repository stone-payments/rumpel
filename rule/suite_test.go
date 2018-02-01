package rule

import (
	"testing"

	check "gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type RuleSuite struct{}

var _ = check.Suite(&RuleSuite{})
