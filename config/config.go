package config

import "flag"

const (
	buildBotUrl     string = "http://10.0.0.5/"
	genericSize     string = "large"
	minRefreshRate  int    = 10
	cacheInvalidate int    = 10
)

type Config struct {
	BuildBotUrl     string
	GenericSize     string //small|large
	RefreshSec      int
	CacheInvalidate int
}

func NewConfig() *Config {
	buildbot := flag.String("buildbot", buildBotUrl, "buildbot url eg. http://10.0.0.1/")
	size := flag.String("size", genericSize, "generic ui size (small|large)")
	refresh := flag.Int("refresh", minRefreshRate, "refresh rate in seconds (min 10 seconds)")
	cache := flag.Int("invalidate", cacheInvalidate, "cache invalidate in seconds (min 5 minutes)")

	flag.Parse()

	cfg := &Config{
		BuildBotUrl:     *buildbot,
		GenericSize:     *size,
		RefreshSec:      *refresh,
		CacheInvalidate: *cache,
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
