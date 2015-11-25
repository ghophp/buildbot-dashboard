package container

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"

	"github.com/ghophp/buildbot-dashboard/buildbot"
	"github.com/ghophp/buildbot-dashboard/cache"
	"github.com/ghophp/buildbot-dashboard/config"
)

// ContainerBag carries all the instantiated dependencies necessary to the handlers work
type ContainerBag struct {
	HashedUrl   string
	GenericSize string
	FilterRegex *regexp.Regexp
	RefreshSec  int
	Cache       *cache.Cache
	Buildbot    buildbot.Buildbot
}

// NewContainerBag return a new instance of the ContainerBag with the instantiated dependencies for the given config
func NewContainerBag(c *config.Config) *ContainerBag {
	hasher := md5.New()
	hasher.Write([]byte(c.BuildBotUrl + c.Filter))

	var filter *regexp.Regexp = nil
	if len(c.Filter) > 0 {
		if r, err := regexp.Compile(c.Filter); err == nil {
			filter = r
		}
	}

	return &ContainerBag{
		HashedUrl:   hex.EncodeToString(hasher.Sum(nil)),
		GenericSize: c.GenericSize,
		RefreshSec:  c.RefreshSec,
		Cache:       cache.NewCache(c.CacheInvalidate),
		FilterRegex: filter,
		Buildbot:    buildbot.NewBuildbotApi(c.BuildBotUrl),
	}
}
