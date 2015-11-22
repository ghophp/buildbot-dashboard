package cache

import (
	"testing"

	gc "github.com/motain/gocheck"
)

var _ = gc.Suite(&CacheSuite{})

type CacheSuite struct{}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *CacheSuite) TestSetCacheShouldCreateFile(c *gc.C) {
	c.Skip("todo")
}

func (s *CacheSuite) TestNonExistingKeyShouldReturnError(c *gc.C) {
	c.Skip("todo")
}

func (s *CacheSuite) TestExistingCacheShouldReturnData(c *gc.C) {
	c.Skip("todo")
}

func (s *CacheSuite) TestExpiredTimeShouldReturnError(c *gc.C) {
	c.Skip("todo")
}
