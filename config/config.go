package config

import "flag"

const (
	genericSize     string = "large"
	minRefreshRate  int    = 10
	cacheInvalidate int    = 10
)

type Config struct {
	BuildBotUrl     string
	GenericSize     string //small|large
	Filter          string
	RefreshSec      int
	CacheInvalidate int
	EmptyBuilders   bool
}

func NewConfig() *Config {
	buildbot := flag.String("buildbot", "", "buildbot url eg. http://10.0.0.1/")
	size := flag.String("size", genericSize, "generic ui size (small|large default large)")
	refresh := flag.Int("refresh", minRefreshRate, "refresh rate in seconds (default and min 10 seconds)")
	cache := flag.Int("invalidate", cacheInvalidate, "cache invalidate in seconds (default and min 5 minutes)")
	empty := flag.Bool("empty", false, "show builders with no builds (default false)")
	filter := flag.String("filter", "", "regex applied over the builder name")

	flag.Parse()

	cfg := &Config{
		BuildBotUrl:     *buildbot,
		GenericSize:     *size,
		RefreshSec:      *refresh,
		CacheInvalidate: *cache,
		EmptyBuilders:   *empty,
		Filter:          *filter,
	}

	if len(cfg.BuildBotUrl) <= 0 {
		panic("buildbot url cannot be empty")
	}

	if cfg.BuildBotUrl[len(cfg.BuildBotUrl)-1:] != "/" {
		cfg.BuildBotUrl = cfg.BuildBotUrl + "/"
	}
	if cfg.GenericSize != "small" && cfg.GenericSize != "large" {
		cfg.GenericSize = genericSize
	}
	if cfg.RefreshSec < minRefreshRate {
		cfg.RefreshSec = minRefreshRate
	}
	if cfg.CacheInvalidate < cacheInvalidate {
		cfg.CacheInvalidate = cacheInvalidate
	}

	return cfg
}
