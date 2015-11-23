package handler

import (
	"github.com/ghophp/buildbot-dashboard/container"
	"github.com/ghophp/render"
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
