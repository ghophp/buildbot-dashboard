package main

import (
	"html/template"

	"github.com/ghophp/buildbot-dashing/config"
	"github.com/ghophp/buildbot-dashing/container"
	"github.com/ghophp/buildbot-dashing/handler"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

func NewRouter(c *container.ContainerBag) *martini.ClassicMartini {
	var (
		indexHandler    = handler.NewIndexHandler(c)
		buildersHandler = handler.NewBuildersHandler(c)
	)

	router := martini.Classic()
	router.Use(martini.Static("static/assets"))
	router.Use(render.Renderer(render.Options{
		Directory:  "static/templates",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".html"},
		Charset:    "UTF-8",
		IndentJSON: true,
		Funcs: []template.FuncMap{
			{
				"genericSize": func() string {
					return c.GenericSize
				},
				"buildbotUrl": func() string {
					return c.BuildBotUrl
				},
				"hashedUrl": func() string {
					return c.HashedUrl
				},
				"displayEmptyBuilder": func() bool {
					return c.EmptyBuilders
				},
			},
		},
	}))

	router.Get("/", indexHandler.ServeHTTP)
	router.Get("/builders", buildersHandler.ServeHTTP)

	handler.AddWs(router, c)

	return router
}

func main() {
	NewRouter(container.NewContainerBag(config.NewConfig())).Run()
}
