package handler

import (
	"github.com/ghophp/buildbot-dashing/container"
	"github.com/martini-contrib/render"
)

type IndexHandler struct {
	c *container.ContainerBag
}

func NewIndexHandler(c *container.ContainerBag) *IndexHandler {
	return &IndexHandler{
		c: c,
	}
}

func (h IndexHandler) ServeHTTP(r render.Render) {
	r.HTML(200, "index", "")
}
