package buildbot

import (
	"io/ioutil"
	"net/http"
)

type Buildbot interface {
	GetUrl() string
	FetchBuilder(id string) ([]byte, error)
	FetchBuilders() ([]byte, error)
}

type BuildbotApi struct {
	url string
}

func NewBuildbotApi(buildbotUrl string) *BuildbotApi {
	return &BuildbotApi{url: buildbotUrl}
}

func (api *BuildbotApi) GetUrl() string {
	return api.url
}

func (api *BuildbotApi) FetchBuilder(id string) ([]byte, error) {
	req, err := http.Get(api.url + "json/builders/" + id + "/builds?select=-1&select=-1&as_text=1")
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
	req, err := http.Get(api.url + "json/builders/?as_text=1")
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
