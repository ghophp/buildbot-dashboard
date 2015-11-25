package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

const InternalCacheFolder string = ".bdd"

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

	cc.path = usr.HomeDir + string(filepath.Separator) + InternalCacheFolder

	_, err = os.Stat(cc.path)
	if os.IsNotExist(err) {
		_ = os.Mkdir(cc.path, 0777)
	}

	return cc
}

func (c *Cache) SetCache(name string, data []byte) error {
	return ioutil.WriteFile(c.path+name, data, 0644)
}

func (c *Cache) GetCache(name string) ([]byte, error) {
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
