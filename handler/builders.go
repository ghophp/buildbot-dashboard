package handler

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"io/ioutil"

	"github.com/ghophp/buildbot-dashing/container"
	"github.com/martini-contrib/render"
)

const BuildersCache string = "builders.json"

type (
	BuildersHandler struct {
		c *container.ContainerBag
	}

	Builder struct {
		Id           string `json:"id"`
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

var timeout = time.Duration(2 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

func GetBuilder(c *container.ContainerBag, id string) (Builder, error) {
	var b Builder

	req, err := http.Get(c.BuildBotUrl + "json/builders/" + id + "?as_text=1")
	if err != nil {
		return b, err
	}

	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&b); err != nil {
		return b, err
	}

	return b, nil
}

func GetBuilders(c *container.ContainerBag) (map[string]Builder, error) {
	var data map[string]Builder

	dataBytes, err := c.Cache.GetCache(BuildersCache)
	if err != nil {
		req, err := http.Get(c.BuildBotUrl + "json/builders/?as_text=1")
		if err != nil {
			return nil, err
		}

		defer req.Body.Close()

		dataBytes, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
	}

	if err := c.Cache.SetCache(BuildersCache, dataBytes); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(dataBytes, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (h BuildersHandler) ServeHTTP(r render.Render) {
	if builders, err := GetBuilders(h.c); err == nil {
		r.JSON(200, builders)
	} else {
		r.Error(500)
	}
}
