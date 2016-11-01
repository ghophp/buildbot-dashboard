package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	cc "github.com/ghophp/buildbot-dashboard/cache"
	"github.com/ghophp/buildbot-dashboard/config"
	"github.com/op/go-logging"

	gc "gopkg.in/check.v1"
)

func SendGetBuildersRequest(cfg *config.Config) *httptest.ResponseRecorder {
	var (
		buildbot = &MockBuildbotApi{url: cfg.BuildBotUrl}
		router   = GetNewTestRouter(cfg, buildbot)
		handler  = NewBuildersHandler(cfg, buildbot, cc.NewFileCache(), logging.MustGetLogger("test"))
	)

	router.Get("/builders", handler.GetBuilders)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/builders", nil)

	router.ServeHTTP(res, req)
	return res
}

func SendGetBuilder(cfg *config.Config, id string) *httptest.ResponseRecorder {
	var (
		buildbot = &MockBuildbotApi{url: cfg.BuildBotUrl}
		router   = GetNewTestRouter(cfg, buildbot)
		handler  = NewBuildersHandler(cfg, buildbot, cc.NewFileCache(), logging.MustGetLogger("test"))
	)

	router.Get("/builder/:id", handler.GetBuilder)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/builder/"+id, nil)

	router.ServeHTTP(res, req)
	return res
}

func (s *HandlerSuite) TestGetBuilderWithWrongUrlMustReturn500(c *gc.C) {
	res := SendGetBuildersRequest(GetNewTestConfig(c, "test_url", ""))
	c.Check(res.Code, gc.Equals, http.StatusInternalServerError)
}

func (s *HandlerSuite) TestGetBuildersMustReturnValidJson(c *gc.C) {
	res := SendGetBuildersRequest(GetNewTestConfig(c, "http://10.0.0.1", ""))
	c.Check(res.Code, gc.Equals, http.StatusOK)

	var data map[string]Builder
	err := json.Unmarshal(res.Body.Bytes(), &data)

	c.Check(err, gc.IsNil)
	c.Check(len(data), gc.Equals, 1)
}

func (s *HandlerSuite) TestGetBuildersMustFilterBuilders(c *gc.C) {
	res := SendGetBuildersRequest(GetNewTestConfig(c, "http://10.0.0.1", "not-select-current"))
	c.Check(res.Code, gc.Equals, http.StatusOK)

	var data map[string]Builder
	err := json.Unmarshal(res.Body.Bytes(), &data)

	c.Check(err, gc.IsNil)
	c.Check(len(data), gc.Equals, 0)
}

func (s *HandlerSuite) TestGetSingleBuilder(c *gc.C) {
	res := SendGetBuilder(GetNewTestConfig(c, "http://10.0.0.1", ""), "buildbot-dashboard")
	c.Check(res.Code, gc.Equals, http.StatusOK)

	fmt.Println(res.Body.String())

	var data Builder
	err := json.Unmarshal(res.Body.Bytes(), &data)

	c.Check(err, gc.IsNil)
	c.Check(data.State, gc.Equals, failedState)
	c.Check(data.Number, gc.Equals, 8)
}

func (s *HandlerSuite) TestGetNonExistentBuilderShouldReturnError(c *gc.C) {
	res := SendGetBuilder(GetNewTestConfig(c, "http://10.0.0.1", ""), "not-exist-builder")
	c.Check(res.Code, gc.Equals, http.StatusInternalServerError)
}
