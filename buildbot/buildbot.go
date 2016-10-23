package buildbot

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/op/go-logging"
)

type Buildbot interface {
	GetUrl() string
	FetchBuilder(id string) ([]byte, error)
	FetchBuilders() ([]byte, error)
}

type BuildbotApi struct {
	url    string
	logger *logging.Logger
}

func NewBuildbotApi(buildbotUrl string, logger *logging.Logger) *BuildbotApi {
	return &BuildbotApi{url: buildbotUrl, logger: logger}
}

func (api *BuildbotApi) GetUrl() string {
	return api.url
}

func (api *BuildbotApi) FetchBuilder(id string) ([]byte, error) {
	builderUrl := api.url + "json/builders/" + id + "/builds?select=-1&select=-1&as_text=1"
	api.logger.Debug(fmt.Sprintf("ready to fetch %s", builderUrl))

	req, err := http.Get(builderUrl)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	req.Body.Close()
	return b, nil
}

func (api *BuildbotApi) FetchBuilders() ([]byte, error) {
	buidersUrl := api.url + "json/builders/?as_text=1"
	api.logger.Debug(fmt.Sprintf("ready to fetch %s", buidersUrl))

	req, err := http.Get(buidersUrl)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	req.Body.Close()
	return b, nil
}
