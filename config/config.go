package config

import (
	"flag"
	"fmt"
)

const (
	genericSize     string = "large"
	minRefreshRate  int    = 20
	cacheInvalidate int    = 5
)

type Config struct {
	BuildBotUrl     string
	Filter          string
	RefreshSec      int
	CacheInvalidate int
}

type ConfigLoader interface {
	Load(cfg *Config)
}

type FlagLoader struct{}

func (f *FlagLoader) Load(cfg *Config) {
	flag.StringVar(&cfg.BuildBotUrl, "buildbot", "", "buildbot url eg. http://10.0.0.1/")
	flag.IntVar(&cfg.RefreshSec, "refresh", minRefreshRate, "refresh rate in seconds (default and min 30 seconds)")
	flag.IntVar(&cfg.CacheInvalidate, "invalidate", cacheInvalidate, "cache invalidate in seconds (default and min 5 minutes)")
	flag.StringVar(&cfg.Filter, "filter", "", "regex applied over the builder name")

	flag.Parse()
}

func NewConfig(loader ConfigLoader) (*Config, error) {
	cfg := &Config{}

	loader.Load(cfg)

	if len(cfg.BuildBotUrl) <= 0 {
		return nil, fmt.Errorf("NewConfig %s", "no buildbot url informed")
	}

	if cfg.BuildBotUrl[len(cfg.BuildBotUrl)-1:] != "/" {
		cfg.BuildBotUrl = cfg.BuildBotUrl + "/"
	}
	if cfg.RefreshSec < minRefreshRate {
		cfg.RefreshSec = minRefreshRate
	}
	if cfg.CacheInvalidate < cacheInvalidate {
		cfg.CacheInvalidate = cacheInvalidate
	}

	return cfg, nil
}
