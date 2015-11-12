package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ghophp/buildbot-dashing/container"
	"github.com/martini-contrib/render"
)

type (
	BuildersHandler struct {
		c *container.ContainerBag
	}

	Builder struct {
		BaseDir      string `json:"basedir"`
		CachedBuilds []int  `json:"cachedBuilds"`
		State        string `json:"state"`
	}
)

func NewBuildersHandler(c *container.ContainerBag) *BuildersHandler {
	return &BuildersHandler{
		c: c,
	}
}

func (h BuildersHandler) ServeHTTP(r render.Render) {
	req, _ := http.Get(h.c.BuildBotUrl + "json/builders/?as_text=1")
	defer req.Body.Close()

	var data map[string]Builder

	dec := json.NewDecoder(req.Body)
	dec.Decode(&data)

	r.JSON(200, data)
}
