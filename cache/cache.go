package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

const InternalCacheFolder string = ".bbd"

type Cache struct {
	refreshTime int
	path        string
}

func NewCache(t int) *Cache {
	cc := &Cache{
		refreshTime: t,
		path:        "/tmp/",
	}

	usr, err := user.Current()
	if err != nil {
		return cc
	}

	cc.path = usr.HomeDir +
		string(filepath.Separator) +
		InternalCacheFolder +
		string(filepath.Separator)

	_, err = os.Stat(cc.path)
	if os.IsNotExist(err) {
		if err = os.Mkdir(cc.path, 0777); err != nil {
			cc.path = "/tmp/"
		}
	}

	return cc
}

func (c *Cache) GetPath() string {
	return c.path
}

func (c *Cache) SetCache(name string, data []byte) error {
	if len(name) <= 0 {
		return fmt.Errorf("invalid cache name")
	}
	return ioutil.WriteFile(c.path+name, data, 0777)
}

func (c *Cache) GetCache(name string) ([]byte, error) {
	if len(name) <= 0 {
		return nil, fmt.Errorf("invalid cache name")
	}

	f, err := os.Stat(c.path + name)
	if os.IsNotExist(err) {
		return nil, err
	}

	if time.Since(f.ModTime()).Minutes() > float64(c.refreshTime) {
		return nil, fmt.Errorf("[GetCache] %s", "cache invalidate")
	}

	dat, err := ioutil.ReadFile(c.path + name)
	if err != nil {
		return nil, err
	}
	return dat, nil
}
