package config

import "flag"

const buildBotUrl string = "http://10.0.0.5/"

type Config struct {
	BuildBotUrl string
}

func NewConfig() *Config {
	buildbot := flag.String("buildbot", buildBotUrl, "buildbot url eg. http://10.0.0.1/")
	flag.Parse()

	return &Config{
		BuildBotUrl: *buildbot,
	}
}
