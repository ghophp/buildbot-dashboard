package buildbot

import (
	"errors"
	"strings"

	"github.com/ghophp/buildbot-dashboard/pool"
	"github.com/op/go-logging"
)

type Buildbot interface {
	GetUrl() string
	FetchBuilder(id string) ([]byte, error)
	FetchBuilders() ([]byte, error)
}

type BuildbotApi struct {
	url    string
	pool   pool.Pool
	logger *logging.Logger
}

func NewBuildbotApi(buildbotUrl string, pool pool.Pool, logger *logging.Logger) *BuildbotApi {
	return &BuildbotApi{url: buildbotUrl, pool: pool, logger: logger}
}

func (api *BuildbotApi) GetUrl() string {
	return api.url
}

func (api *BuildbotApi) FetchBuilder(id string) ([]byte, error) {
	var (
		builderUrl = api.url + "json/builders/" + id + "/builds?select=-1&select=-1&as_text=1"
		listener   = make(chan string)
	)

	api.pool.Fetch(builderUrl, listener)

	select {
	case resp := <-listener:
		if strings.Contains(resp, pool.RequestError) {
			return nil, errors.New(resp)
		}
		return []byte(resp), nil
	}
}

func (api *BuildbotApi) FetchBuilders() ([]byte, error) {
	var (
		buidersUrl = api.url + "json/builders/?as_text=1"
		listener   = make(chan string)
	)

	api.pool.Fetch(buidersUrl, listener)

	select {
	case resp := <-listener:
		if strings.Contains(resp, pool.RequestError) {
			return nil, errors.New(resp)
		}
		return []byte(resp), nil
	}
}
