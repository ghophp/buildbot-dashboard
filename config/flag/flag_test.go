package flag

import (
	"os"
	"testing"

	gc "gopkg.in/check.v1"

	"github.com/ghophp/buildbot-dashboard/config"
)

var _ = gc.Suite(&FlagLoaderSuite{})

type FlagLoaderSuite struct{}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *FlagLoaderSuite) TestNewFlagLoaderShouldLoadValuesIntoConfig(c *gc.C) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{
		"buildbot-dashboard",
		"-templates=/etc/bb/",
		"-buildbot=http://10.0.0.1/",
		"-refresh=10",
		"-invalidate=15",
		"-filter=*",
	}

	var (
		cfg    = &config.Config{}
		loader = NewFlagLoader()
	)

	loader.Load(cfg)
	c.Check(cfg.StaticPath, gc.Equals, "/etc/bb/")
	c.Check(cfg.BuildBotUrl, gc.Equals, "http://10.0.0.1/")
	c.Check(cfg.RefreshSec, gc.Equals, 10)
	c.Check(cfg.CacheInvalidate, gc.Equals, 15)
	c.Check(cfg.FilterStr, gc.Equals, "*")
}
