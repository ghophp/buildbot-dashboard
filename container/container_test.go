package container

import (
	"testing"

	gc "github.com/motain/gocheck"
)

var _ = gc.Suite(&ContainerBagSuite{})

type ContainerBagSuite struct{}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *ContainerBagSuite) TestNewContainerBagMustInitializeComponentsIfConfigProvided(c *gc.C) {
	c.Skip("todo")
}

func (s *ContainerBagSuite) TestNewContainerShouldReturnProperHashedUrl(c *gc.C) {
	c.Skip("todo")
}

func (s *ContainerBagSuite) TestNewContainerShouldSetFilterForValidRegex(c *gc.C) {
	c.Skip("todo")
}

func (s *ContainerBagSuite) TestNewContainerShouldNotSetFilterForInvalidRegex(c *gc.C) {
	c.Skip("todo")
}
