package main

import (
	"html/template"
	"strconv"

	bb "github.com/ghophp/buildbot-dashboard/buildbot"
	cc "github.com/ghophp/buildbot-dashboard/cache"
	"github.com/ghophp/buildbot-dashboard/config"
	"github.com/ghophp/buildbot-dashboard/handler"

	"github.com/ghophp/render"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/staticbin"
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

	cfg, err := config.NewConfig(&config.FlagLoader{})
	if err != nil {
		panic(err)
	}

	var (
		cache    = cc.NewFileCache()
		buildbot = bb.NewBuildbotApi(cfg.BuildBotUrl, log)

		indexHandler    = handler.NewIndexHandler()
		buildersHandler = handler.NewBuildersHandler(cfg, buildbot, cache, log)
	)

	router := martini.Classic()
	router.Use(staticbin.Static("static/assets", Asset))
	router.Use(render.RendererBin(Asset, AssetNames(), render.Options{
		Directory:  "static/templates",
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
