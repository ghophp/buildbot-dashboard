package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ghophp/buildbot-dashboard/container"
	"github.com/ghophp/render"
	"github.com/go-martini/martini"
)

const (
	buildingState   string = "building"
	failedState     string = "failed"
	successfulState string = "successful"
	warningState    string = "warnings"
	exceptionState  string = "exception"
)

var validStates []string = []string{
	buildingState,
	failedState,
	successfulState,
	warningState,
	exceptionState,
}

type (
	BuildersHandler struct {
		c *container.ContainerBag
	}

	Builder struct {
		Id         string   `json:"id"`
		State      string   `json:"state"`
		Reason     string   `json:"reason"`
		Blame      []string `json:"blame"`
		Number     int      `json:"number"`
		Slave      string   `json:"slave"`
		LastUpdate string   `json:"last_update"`
	}

	DetailedBuilder struct {
		Blame  []string  `json:"blame"`
		Number int       `json:"number"`
		Reason string    `json:"reason"`
		Slave  string    `json:"slave"`
		Times  []float64 `json:"times"`
		Text   []string  `json:"text"`
		Error  string    `json:"error"`
	}
)

func NewBuildersHandler(c *container.ContainerBag) *BuildersHandler {
	return &BuildersHandler{
		c: c,
	}
}

func isValidState(v string) bool {
	for _, s := range validStates {
		if v == s {
			return true
		}
	}
	return false
}

func (h BuildersHandler) fetchBuilder(id string) (*Builder, error) {
	var b map[string]DetailedBuilder

	data, err := h.c.Buildbot.FetchBuilder(id)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}

	if current, ok := b["-1"]; ok && current.Error == "" {

		builder := &Builder{
			Id:         id,
			Blame:      current.Blame,
			Number:     current.Number,
			Slave:      current.Slave,
			Reason:     current.Reason,
			State:      buildingState,
			LastUpdate: strconv.Itoa(int(time.Now().Unix())),
		}

		if len(current.Times) > 0 {
			builder.LastUpdate = strconv.FormatFloat(current.Times[0], 'f', 6, 64)
		}

		if len(current.Text) > 0 {
			for _, v := range current.Text {
				if isValidState(v) {
					builder.State = v
					break
				}
			}
		}

		return builder, nil
	}

	return nil, fmt.Errorf("[GetBuilder] %s", "no last build defined")
}

// fetchBuilders will fetch the builders from buildbot, if the fresh parameter is equal true
// it will not respect cache
func (h BuildersHandler) fetchBuilders(fresh bool) (map[string]Builder, error) {
	var data map[string]Builder

	dataBytes, err := h.c.Cache.GetCache(h.c.HashedUrl)
	if fresh || err != nil {
		dataBytes, err = h.c.Buildbot.FetchBuilders()
		if err != nil {
			return nil, err
		}

		if err := h.c.Cache.SetCache(h.c.HashedUrl, dataBytes); err != nil {
			return nil, err
		}
	}

	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return nil, err
	}

	if h.c.FilterRegex != nil {
		for key, _ := range data {
			if !h.c.FilterRegex.MatchString(key) {
				delete(data, key)
			}
		}
	}

	return data, nil
}

func (h BuildersHandler) GetBuilders(req *http.Request, r render.Render) {
	fresh, err := strconv.ParseBool(req.URL.Query().Get("fresh"))
	if err != nil {
		fresh = false
	}

	if builders, err := h.fetchBuilders(fresh); err == nil {
		r.JSON(200, builders)
	} else {
		r.Error(500)
	}
}

func (h BuildersHandler) GetBuilder(params martini.Params, r render.Render) {
	builder, err := h.fetchBuilder(params["id"])
	if err != nil {
		fmt.Println(err)
		r.Error(500)
	} else {
		r.JSON(200, builder)
	}
}
