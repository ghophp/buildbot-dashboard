package main

import (
	"html/template"
	"os"
	"strconv"

	bb "github.com/ghophp/buildbot-dashboard/buildbot"
	cc "github.com/ghophp/buildbot-dashboard/cache"
	"github.com/ghophp/buildbot-dashboard/config"
	"github.com/ghophp/buildbot-dashboard/config/env"
	"github.com/ghophp/buildbot-dashboard/config/flag"
	"github.com/ghophp/buildbot-dashboard/handler"
	"github.com/ghophp/buildbot-dashboard/pool"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/op/go-logging"
)

const LoggerPrefix = "BUILDBOT-DASHBOARD"

var log = logging.MustGetLogger(LoggerPrefix)

// Log format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func main() {
	logging.SetFormatter(format)

	cfg, err := config.NewConfig([]config.Loader{
		flag.NewFlagLoader(),
		env.NewEnvLoader(),
	})
	if err != nil {
		panic(err)
	}

	var (
		cache    = cc.NewFileCache()
		buildbot = bb.NewBuildbotApi(cfg.BuildBotUrl, pool.NewRequestPool(), log)

		indexHandler    = handler.NewIndexHandler()
		buildersHandler = handler.NewBuildersHandler(cfg, buildbot, cache, log)
	)

	var staticPath = "./static/"
	if len(cfg.StaticPath) > 0 {
		staticPath = cfg.StaticPath
	}

	if _, err := os.Stat(staticPath); os.IsNotExist(err) {
		panic(err)
	}

	router := martini.Classic()
	router.Use(martini.Static(staticPath + "assets"))
	router.Use(render.Renderer(render.Options{
		Directory:  staticPath + "templates",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".html"},
		Charset:    "UTF-8",
		IndentJSON: true,
		Funcs: []template.FuncMap{
			{
				"refreshSec": func() string {
					return strconv.Itoa(cfg.RefreshSec)
				},
				"buildbotUrl": func() string {
					return buildbot.GetUrl()
				},
				"hashedUrl": func() string {
					return cfg.HashedUrl
				},
			},
		},
	}))

	router.Get("/", indexHandler.ServeHTTP)
	router.Get("/builders", buildersHandler.GetBuilders)
	router.Get("/builder/:id", buildersHandler.GetBuilder)

	router.Run()
}
