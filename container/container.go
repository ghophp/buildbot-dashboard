package container

import (
	"github.com/ghophp/buildbot-dashing/cache"
	"github.com/ghophp/buildbot-dashing/config"
)

// ContainerBag carries all the instantiated dependencies necessary to the handlers work
type ContainerBag struct {
	BuildBotUrl string
	GenericSize string
	RefreshSec  int
	Cache       *cache.Cache
}

// NewContainerBag return a new instance of the ContainerBag with the instantiated dependencies for the given config
func NewContainerBag(c *config.Config) *ContainerBag {
	return &ContainerBag{
		BuildBotUrl: c.BuildBotUrl,
		GenericSize: c.GenericSize,
		RefreshSec:  c.RefreshSec,
		Cache:       cache.NewCache(c.CacheInvalidate),
	}
}
