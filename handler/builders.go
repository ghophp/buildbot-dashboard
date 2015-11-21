package handler

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"io/ioutil"

	"github.com/ghophp/buildbot-dashing/container"
	"github.com/martini-contrib/render"
)

const (
	buildingState   string = "building"
	failedState     string = "failed"
	successfulState string = "successful"
	warningState    string = "warnings"
)

type (
	BuildersHandler struct {
		c *container.ContainerBag
	}

	Builder struct {
		Id           string   `json:"id"`
		CachedBuilds []int    `json:"cachedBuilds"`
		State        string   `json:"state"`
		Reason       string   `json:"reason"`
		Blame        []string `json:"blame"`
		Number       int      `json:"number"`
		Slave        string   `json:"slave"`
		LastUpdate   string   `json:"last_update"`
	}

	DetailedBuilder struct {
		Blame  []string  `json:"blame"`
		Number int       `json:"number"`
		Reason string    `json:"reason"`
		Slave  string    `json:"slave"`
		Times  []float64 `json:"times"`
		Text   []string  `json:"text"`
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

func GetBuilder(c *container.ContainerBag, id string, builder Builder) (Builder, error) {
	var b map[string]DetailedBuilder

	req, err := http.Get(c.BuildBotUrl + "json/builders/" + id + "/builds?select=-1&select=-1&as_text=1")
	if err != nil {
		return builder, err
	}

	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&b); err != nil {
		return builder, err
	}

	if current, ok := b["-1"]; ok {

		builder.Id = id
		builder.Blame = current.Blame
		builder.Number = current.Number
		builder.Slave = current.Slave
		builder.Reason = current.Reason
		builder.State = buildingState
		builder.LastUpdate = strconv.Itoa(int(time.Now().Unix()))

		if len(current.Times) > 0 {
			builder.LastUpdate = strconv.FormatFloat(current.Times[0], 'f', 6, 64)
		}
		if len(current.Text) >= 2 {
			if current.Text[0] == failedState {
				builder.State = failedState
			} else if current.Text[1] == successfulState {
				builder.State = successfulState
			}
		} else if len(current.Text) == 1 && current.Text[0] == warningState {
			builder.State = warningState
		}

		return builder, nil
	}

	return builder, fmt.Errorf("[GetBuilder] %s", "no last build defined")
}

func GetBuilders(c *container.ContainerBag) (map[string]Builder, error) {
	var data map[string]Builder

	dataBytes, err := c.Cache.GetCache(c.HashedUrl)
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

		if err := c.Cache.SetCache(c.HashedUrl, dataBytes); err != nil {
			return nil, err
		}
	}

	if err = json.Unmarshal(dataBytes, &data); err != nil {
		return nil, err
	}

	if c.FilterRegex != nil {
		var del []string
		for key, _ := range data {
			if !c.FilterRegex.MatchString(key) {
				del = append(del, key)
			}
		}
		if len(del) > 0 {
			for _, k := range del {
				delete(data, k)
			}
		}
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
