package flag

import (
	"flag"

	"github.com/ghophp/buildbot-dashboard/config"
)

type Loader struct{}

func NewFlagLoader() *Loader {
	return new(Loader)
}

func (f *Loader) Load(cfg *config.Config) {
	flag.StringVar(&cfg.StaticPath, "templates", "", "path to static folder (default ./static/)")
	flag.StringVar(&cfg.BuildBotUrl, "buildbot", "", "buildbot url eg. http://10.0.0.1/")
	flag.IntVar(&cfg.RefreshSec, "refresh", config.MinRefreshRate, "refresh rate in seconds (default and min 30 seconds)")
	flag.IntVar(&cfg.CacheInvalidate, "invalidate", config.CacheInvalidate, "cache invalidate in seconds (default and min 300 seconds)")
	flag.StringVar(&cfg.FilterStr, "filter", "", "regex applied over the builder name")

	flag.Parse()
}
