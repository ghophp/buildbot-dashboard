package handler

import (
	"fmt"
	"html/template"
	"strconv"
	"testing"

	bb "github.com/ghophp/buildbot-dashboard/buildbot"
	"github.com/ghophp/buildbot-dashboard/config"

	"github.com/ghophp/render"
	"github.com/go-martini/martini"
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
	cfg.FilterStr = f.filter
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
	fmt.Println(fmt.Sprintf("try to fetch %s", api.url))
	if api.url == "test_url/" {
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

func GetNewTestConfig(c *gc.C, url string, filter string) *config.Config {
	cfg, err := config.NewConfig(&MockLoader{
		url:    url,
		filter: filter,
	})

	if err != nil {
		c.Error(err)
	}

	return cfg
}

func GetNewTestRouter(cfg *config.Config, buildbot bb.Buildbot) *martini.ClassicMartini {
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
					return strconv.Itoa(cfg.RefreshSec)
				},
				"buildbotUrl": func() string {
					return buildbot.GetUrl()
				},
				"hashedUrl": func() string {
					return cfg.HashedUrl
				},
			},
		},
	}))

	return m
}
