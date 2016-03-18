package main

import (
	"html/template"
	"strconv"

	"github.com/ghophp/buildbot-dashboard/config"
	"github.com/ghophp/buildbot-dashboard/container"
	"github.com/ghophp/buildbot-dashboard/handler"

	"github.com/ghophp/render"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/staticbin"
)

func NewRouter(c *container.ContainerBag) *martini.ClassicMartini {
	var (
		indexHandler    = handler.NewIndexHandler(c)
		buildersHandler = handler.NewBuildersHandler(c)
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
					return strconv.Itoa(c.RefreshSec)
				},
				"buildbotUrl": func() string {
					return c.Buildbot.GetUrl()
				},
				"hashedUrl": func() string {
					return c.HashedUrl
				},
			},
		},
	}))

	router.Get("/", indexHandler.ServeHTTP)
	router.Get("/builders", buildersHandler.GetBuilders)
	router.Get("/builder/:id", buildersHandler.GetBuilder)

	return router
}

func main() {
	cfg, err := config.NewConfig(&config.FlagLoader{})
	if err != nil {
		panic(err)
	}

	NewRouter(container.NewContainerBag(cfg)).Run()
}
