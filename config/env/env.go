package env

import (
	"os"
	"strconv"

	"github.com/ghophp/buildbot-dashboard/config"
)

type Loader struct{}

func NewEnvLoader() *Loader {
	return new(Loader)
}

func (l *Loader) Load(cfg *config.Config) {
	cfg.StaticPath = os.Getenv("STATIC_PATH")
	cfg.BuildBotUrl = os.Getenv("BUILDBOT_URL")
	cfg.FilterStr = os.Getenv("FILTER_REGEX")

	var (
		refreshString    = os.Getenv("REFRESH_INTERVAL")
		invalidateString = os.Getenv("INVALIDATE_INTERVAL")
	)

	if v, err := strconv.Atoi(refreshString); err == nil {
		cfg.RefreshSec = v
	}
	if v, err := strconv.Atoi(invalidateString); err == nil {
		cfg.CacheInvalidate = v
	}
}
