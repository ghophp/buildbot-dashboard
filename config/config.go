package config

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"regexp"
)

const (
	minRefreshRate  int = 20
	cacheInvalidate int = 5
)

type Config struct {
	BuildBotUrl     string
	HashedUrl       string
	FilterStr       string
	FilterRegex     *regexp.Regexp
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
	flag.StringVar(&cfg.FilterStr, "filter", "", "regex applied over the builder name")

	flag.Parse()
}

func NewConfig(loader ConfigLoader) (*Config, error) {
	cfg := &Config{}
	loader.Load(cfg)

	hasher := md5.New()
	hasher.Write([]byte(cfg.BuildBotUrl + cfg.FilterStr))

	var filter *regexp.Regexp = nil
	if len(cfg.FilterStr) > 0 {
		if r, err := regexp.Compile(cfg.FilterStr); err == nil {
			filter = r
		}
	}

	cfg.HashedUrl = hex.EncodeToString(hasher.Sum(nil))
	cfg.FilterRegex = filter

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
