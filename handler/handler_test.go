package handler

import (
	"html/template"
	"strconv"
	"testing"

	"github.com/ghophp/buildbot-dashboard/config"
	"github.com/ghophp/buildbot-dashboard/container"
	"github.com/ghophp/render"
	"github.com/go-martini/martini"

	"fmt"

	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(&HandlerSuite{})

type HandlerSuite struct{}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

type MockLoader struct {
	url    string
	filter string
}

type MockBuildbotApi struct {
	url string
}

func (f *MockLoader) Load(cfg *config.Config) {
	cfg.BuildBotUrl = f.url
	cfg.RefreshSec = 10
	cfg.CacheInvalidate = 10
	cfg.Filter = f.filter
}

func (api *MockBuildbotApi) GetUrl() string {
	return api.url
}

func (api *MockBuildbotApi) FetchBuilder(id string) ([]byte, error) {
	if id != "buildbot-dashboard" {
		return nil, fmt.Errorf("test error")
	}

	json := `{
		"-1": {
			"builderName": "buildbot-dashboard",
			"number": 8,
			"reason": "A build was forced by 'buildbot <buildbot@localhost>': force build",
			"results": 2,
			"slave": "i-d48e3278",
			"text": [
				"failed",
				"shell_1",
				"shell_2",
				"shell_3"
			],
			"times": [
				1448379239.908776,
				1448379252.193469
			]
		}
	}`

	return []byte(json), nil
}

func (api *MockBuildbotApi) FetchBuilders() ([]byte, error) {
	if api.url == "test_url" {
		return nil, fmt.Errorf("test error")
	}

	json := `{
		"buildbot-dashboard": {
			"basedir": "buildbot-dashboard",
			"cachedBuilds": [
				0,
				1,
				2,
				3,
				4,
				5,
				6,
				7,
				8
			],
			"schedulers": [
				"buildbot-dashboard",
				"force"
			],
			"slaves": [
				"i-d48e3278"
			],
			"state": "idle"
		}
	}`

	return []byte(json), nil
}

func GetNewContainerBag(c *gc.C, url string, filter string) *container.ContainerBag {
	cfg, err := config.NewConfig(&MockLoader{
		url:    url,
		filter: filter,
	})

	if err != nil {
		c.Error(err)
	}

	ctn := container.NewContainerBag(cfg)
	ctn.Buildbot = &MockBuildbotApi{url: url}

	return ctn
}

func GetNewTestRouter(ctx *container.ContainerBag) *martini.ClassicMartini {
	m := martini.Classic()
	m.Use(martini.Static("../static/assets"))
	m.Use(render.Renderer(render.Options{
		Directory:  "../static/templates",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".html"},
		Charset:    "UTF-8",
		IndentJSON: true,
		Funcs: []template.FuncMap{
			{
				"refreshSec": func() string {
					return strconv.Itoa(ctx.RefreshSec)
				},
				"buildbotUrl": func() string {
					return ctx.Buildbot.GetUrl()
				},
				"hashedUrl": func() string {
					return ctx.HashedUrl
				},
			},
		},
	}))

	return m
}
