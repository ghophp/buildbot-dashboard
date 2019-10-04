package config

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"regexp"
)

const (
	MinRefreshRate  int = 20
	CacheInvalidate int = 300
)

type Config struct {
	StaticPath      string
	BuildBotUrl     string
	HashedUrl       string
	FilterStr       string
	FilterRegex     *regexp.Regexp
	RefreshSec      int
	CacheInvalidate int
}

// Loader defines the behaviour to load the values from
// multiple sources that will populate the config
type Loader interface {
	Load(*Config)
}

func NewConfig(loaders []Loader) (*Config, error) {
	if loaders == nil || len(loaders) <= 0 {
		return nil, fmt.Errorf("NewConfig %v", "no loader provided")
	}

	var (
		cfg    = new(Config)
		loaded = false
	)

	for _, l := range loaders {
		l.Load(cfg)
		if err := cfg.Validate(); err == nil {
			loaded = true
			break
		}
	}

	if !loaded {
		return nil, fmt.Errorf("NewConfig %v", "invalid configuration provided")
	}

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

	if cfg.BuildBotUrl[len(cfg.BuildBotUrl)-1:] != "/" {
		cfg.BuildBotUrl = cfg.BuildBotUrl + "/"
	}
	if cfg.RefreshSec < MinRefreshRate {
		cfg.RefreshSec = MinRefreshRate
	}
	if cfg.CacheInvalidate < CacheInvalidate {
		cfg.CacheInvalidate = CacheInvalidate
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if len(c.BuildBotUrl) <= 0 {
		return fmt.Errorf("Config.Validate %s", "no buildbot url informed")
	}
	return nil
}
