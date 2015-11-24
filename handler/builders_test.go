package handler

import (
	"net/http"
	"net/http/httptest"

	"encoding/json"

	"github.com/ghophp/buildbot-dashboard/container"
	gc "gopkg.in/check.v1"
)

func ServeHTTPForContainer(ctn *container.ContainerBag) *httptest.ResponseRecorder {
	router := GetNewTestRouter(ctn)

	handler := NewBuildersHandler(ctn)
	router.Get("/foobar", handler.ServeHTTP)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	router.ServeHTTP(res, req)

	return res
}

func (s *HandlerSuite) TestGetBuilderWithWrongUrlMustReturn500(c *gc.C) {
	res := ServeHTTPForContainer(GetNewContainerBag(c, "test_url", ""))
	c.Check(res.Code, gc.Equals, http.StatusInternalServerError)
}

func (s *HandlerSuite) TestGetBuildersMustReturnValidJson(c *gc.C) {
	res := ServeHTTPForContainer(GetNewContainerBag(c, "http://10.0.0.1", ""))
	c.Check(res.Code, gc.Equals, http.StatusOK)

	var data map[string]Builder
	err := json.Unmarshal(res.Body.Bytes(), &data)

	c.Check(err, gc.IsNil)
	c.Check(len(data), gc.Equals, 1)
}

func (s *HandlerSuite) TestGetBuildersMustFilterBuilders(c *gc.C) {
	res := ServeHTTPForContainer(GetNewContainerBag(c, "http://10.0.0.1", "not-select-current"))
	c.Check(res.Code, gc.Equals, http.StatusOK)

	var data map[string]Builder
	err := json.Unmarshal(res.Body.Bytes(), &data)

	c.Check(err, gc.IsNil)
	c.Check(len(data), gc.Equals, 0)
}

func (s *HandlerSuite) TestGetBuilderMustUpdateBuilder(c *gc.C) {
	b := Builder{
		Id:         "buildbot-dashboard",
		State:      "",
		Reason:     "",
		Blame:      []string{},
		Number:     0,
		Slave:      "",
		LastUpdate: "",
	}

	b, err := GetBuilder(GetNewContainerBag(c, "http://10.0.0.1", ""), b.Id, b)
	c.Check(err, gc.IsNil)
	c.Check(b.State, gc.Equals, failedState)
	c.Check(b.Number, gc.Equals, 8)
}

func (s *HandlerSuite) TestGetBuilderMustReturnErrorForWrongId(c *gc.C) {
	b := Builder{
		Id:         "xxxx",
		State:      "",
		Reason:     "",
		Blame:      []string{},
		Number:     0,
		Slave:      "",
		LastUpdate: "",
	}

	b, err := GetBuilder(GetNewContainerBag(c, "http://10.0.0.1", ""), b.Id, b)
	c.Check(err, gc.NotNil)
}
