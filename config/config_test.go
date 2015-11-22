package config

import (
	"testing"

	gc "github.com/motain/gocheck"
)

var _ = gc.Suite(&ConfigSuite{})

type ConfigSuite struct{}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *ConfigSuite) TestNewConfigShouldReturnErrorIfNoBuildBot(c *gc.C) {
	cfg, err := NewConfig()

	c.Check(err, gc.NotNil)
	c.Check(cfg, gc.IsNil)
}
