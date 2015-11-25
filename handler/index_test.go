package handler

import (
	"net/http"
	"net/http/httptest"

	gc "gopkg.in/check.v1"
)

func (s *HandlerSuite) TestIndexMustReturnOK(c *gc.C) {
	ctx := GetNewContainerBag(c, "http://10.0.0.1", "")
	router := GetNewTestRouter(ctx)

	handler := NewIndexHandler(ctx)
	router.Get("/foobar", handler.ServeHTTP)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	router.ServeHTTP(res, req)
	c.Check(res.Code, gc.Equals, http.StatusOK)
}
