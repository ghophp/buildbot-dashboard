package cache

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"

	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(&CacheSuite{})

type CacheSuite struct{}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *CacheSuite) SetUpTest(c *gc.C) {
	usr, err := user.Current()
	if err != nil {
		c.Error(err)
	}

	os.Remove(usr.HomeDir + string(filepath.Separator) + InternalCacheFolder)
}

func (s *CacheSuite) TestSetCacheEmptyNameShouldReturnError(c *gc.C) {
	cache := NewFileCache()

	err := cache.SetCache("", []byte("test content"), 10)
	c.Check(err, gc.NotNil)

	_, err = cache.GetCache("")
	c.Check(err, gc.NotNil)
}

func (s *CacheSuite) TestSetCacheShouldCreateFile(c *gc.C) {
	cache := NewFileCache()
	cache.SetCache("test", []byte("test content"), 10)

	_, err := os.Stat(cache.GetPath() + "test")
	c.Check(err, gc.IsNil)
}

func (s *CacheSuite) TestNonExistingKeyShouldReturnError(c *gc.C) {
	cache := NewFileCache()
	cache.SetCache("test", []byte("test content"), 10)

	data, err := cache.GetCache("test")

	c.Check(err, gc.IsNil)
	c.Check(string(data), gc.Equals, "test content")
}

func (s *CacheSuite) TestExpiredTimeShouldReturnError(c *gc.C) {
	cache := NewFileCache()
	cache.SetCache("test", []byte("test content"), -1)

	_, err := cache.GetCache("test")
	c.Check(err, gc.NotNil)
}
