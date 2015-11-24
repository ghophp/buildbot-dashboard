package container

import (
	"crypto/md5"
	"encoding/hex"
	"testing"

	"github.com/ghophp/buildbot-dashboard/config"

	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(&ContainerBagSuite{})

type ContainerBagSuite struct{}

type MockLoader struct{}

func (f *MockLoader) Load(cfg *config.Config) {
	cfg.BuildBotUrl = "http://10.0.0.1"
	cfg.GenericSize = "small"
	cfg.RefreshSec = 10
	cfg.CacheInvalidate = 10
	cfg.Filter = ".*"
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *ContainerBagSuite) TestNewContainerBagMustInitializeComponentsIfConfigProvided(c *gc.C) {
	cfg, err := config.NewConfig(&MockLoader{})
	c.Check(err, gc.IsNil)

	hasher := md5.New()
	hasher.Write([]byte(cfg.BuildBotUrl + cfg.Filter))

	ctx := NewContainerBag(cfg)

	c.Check(ctx.HashedUrl, gc.Equals, hex.EncodeToString(hasher.Sum(nil)))
	c.Check(ctx.GenericSize, gc.Equals, cfg.GenericSize)
	c.Check(ctx.RefreshSec, gc.Equals, cfg.RefreshSec)
	c.Check(ctx.FilterRegex, gc.NotNil)
	c.Check(ctx.Cache, gc.NotNil)
	c.Check(ctx.Buildbot.GetUrl(), gc.Equals, cfg.BuildBotUrl)
	c.Check(ctx.Buildbot, gc.NotNil)

	cfg.Filter = "))))"
	ctxRegex := NewContainerBag(cfg)
	c.Check(ctxRegex.FilterRegex, gc.IsNil)
}
