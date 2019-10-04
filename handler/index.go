package handler

import "github.com/martini-contrib/render"

type IndexHandler struct{}

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

func (h IndexHandler) ServeHTTP(r render.Render) {
	r.HTML(200, "index", "")
}
