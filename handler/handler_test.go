package handler

import (
	"testing"

	gc "github.com/motain/gocheck"
)

var _ = gc.Suite(&HandlerSuite{})

type HandlerSuite struct{}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *HandlerSuite) TestIndexMustReturnOK(c *gc.C) {
	c.Skip("todo")
}

func (s *HandlerSuite) TestGetBuildersMustReturnListOfBuilders(c *gc.C) {
	c.Skip("todo: to make this test work, we will have to separate the 'buildbot' into a entitiy, and make it an interface, mockable, so we can mock the idea of request the list")
}

func (s *HandlerSuite) TestGetSingleBuilderMustReturnBuilder(c *gc.C) {
	c.Skip("todo")
}
