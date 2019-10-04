package env

import (
	"os"
	"testing"

	gc "gopkg.in/check.v1"

	"github.com/ghophp/buildbot-dashboard/config"
)

var _ = gc.Suite(&EnvLoaderSuite{})

type EnvLoaderSuite struct{}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *EnvLoaderSuite) TestNewEnvLoaderShouldLoadValuesIntoConfig(c *gc.C) {
	os.Setenv("STATIC_PATH", "/etc/bb/")
	os.Setenv("BUILDBOT_URL", "http://10.0.0.1/")
	os.Setenv("FILTER_REGEX", "*")
	os.Setenv("REFRESH_INTERVAL", "10")
	os.Setenv("INVALIDATE_INTERVAL", "15")

	var (
		cfg    = &config.Config{}
		loader = NewEnvLoader()
	)

	loader.Load(cfg)
	c.Check(cfg.StaticPath, gc.Equals, "/etc/bb/")
	c.Check(cfg.BuildBotUrl, gc.Equals, "http://10.0.0.1/")
	c.Check(cfg.RefreshSec, gc.Equals, 10)
	c.Check(cfg.CacheInvalidate, gc.Equals, 15)
	c.Check(cfg.FilterStr, gc.Equals, "*")
}
