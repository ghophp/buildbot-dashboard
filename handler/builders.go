package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	bb "github.com/ghophp/buildbot-dashboard/buildbot"
	cc "github.com/ghophp/buildbot-dashboard/cache"
	"github.com/ghophp/buildbot-dashboard/config"
	"github.com/op/go-logging"

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
		cfg      *config.Config
		buildbot bb.Buildbot
		cache    cc.Cache
		logger   *logging.Logger
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

func NewBuildersHandler(cfg *config.Config, buildbot bb.Buildbot, cache cc.Cache, logger *logging.Logger) *BuildersHandler {
	return &BuildersHandler{cfg, buildbot, cache, logger}
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
	var data map[string]DetailedBuilder

	dataBytes, err := h.cache.GetCache(h.cfg.HashedUrl + "_" + id)
	if err != nil {
		h.logger.Errorf("failed to retrieve from cache %s", err)

		dataBytes, err = h.buildbot.FetchBuilder(id)
		if err != nil {
			return nil, err
		}

		err = h.cache.SetCache(h.cfg.HashedUrl+"_"+id, dataBytes, 10)
		if err != nil {
			return nil, err
		}
	}

	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return nil, err
	}

	if current, ok := data["-1"]; ok && current.Error == "" {
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

	dataBytes, err := h.cache.GetCache(h.cfg.HashedUrl)
	if fresh || err != nil {
		if err != nil {
			h.logger.Errorf("failed to retrieve from cache %s", err)
		}

		dataBytes, err = h.buildbot.FetchBuilders()
		if err != nil {
			return nil, err
		}

		err = h.cache.SetCache(h.cfg.HashedUrl, dataBytes, h.cfg.CacheInvalidate)
		if err != nil {
			return nil, err
		}
	}

	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return nil, err
	}

	if h.cfg.FilterRegex != nil {
		for key, _ := range data {
			if !h.cfg.FilterRegex.MatchString(key) {
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

	builders, err := h.fetchBuilders(fresh)
	if err != nil {
		h.logger.Error(err)
		r.Error(500)
	} else {
		r.JSON(200, builders)
	}
}

func (h BuildersHandler) GetBuilder(params martini.Params, r render.Render) {
	builder, err := h.fetchBuilder(params["id"])
	if err != nil {
		h.logger.Error(err)
		r.Error(500)
	} else {
		r.JSON(200, builder)
	}
}
