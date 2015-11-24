package config

import (
	"testing"

	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(&ConfigSuite{})

type ConfigSuite struct{}

type MockLoader struct{}

func (f *MockLoader) Load(cfg *Config) {
	cfg.BuildBotUrl = "http://10.0.0.1"
	cfg.GenericSize = "will_be_replaced"
	cfg.RefreshSec = 0
	cfg.CacheInvalidate = 0
	cfg.Filter = ""
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *ConfigSuite) TestNewConfigShouldReturnErrorIfNoBuildBot(c *gc.C) {
	cfg, err := NewConfig(&FlagLoader{})

	c.Check(err, gc.NotNil)
	c.Check(cfg, gc.IsNil)
}

func (s *ConfigSuite) TestNewConfigShouldHaveDefaultValues(c *gc.C) {
	cfg, err := NewConfig(&MockLoader{})

	c.Check(err, gc.IsNil)

	c.Check(cfg.BuildBotUrl, gc.Equals, "http://10.0.0.1/")
	c.Check(cfg.GenericSize, gc.Equals, "large")
	c.Check(cfg.RefreshSec, gc.Equals, 10)
	c.Check(cfg.CacheInvalidate, gc.Equals, 5)
	c.Check(cfg.Filter, gc.Equals, "")
}
