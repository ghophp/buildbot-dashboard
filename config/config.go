package config

import "flag"

const (
	buildBotUrl string = "http://10.0.0.5/"
	genericSize string = "large"
)

type Config struct {
	BuildBotUrl string
	GenericSize string //small|large
}

func NewConfig() *Config {
	buildbot := flag.String("buildbot", buildBotUrl, "buildbot url eg. http://10.0.0.1/")
	size := flag.String("size", genericSize, "generic ui size (small|large)")

	flag.Parse()

	cfg := &Config{
		BuildBotUrl: *buildbot,
		GenericSize: *size,
	}

	if cfg.GenericSize != "small" && cfg.GenericSize != "large" {
		cfg.GenericSize = genericSize
	}

	return cfg
}
