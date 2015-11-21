package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Cache struct {
	memoryCache map[string][]byte
	refreshTime int
}

func NewCache(t int) *Cache {
	return &Cache{
		memoryCache: make(map[string][]byte),
		refreshTime: t,
	}
}

func (c *Cache) SetCache(name string, data []byte) error {
	err := ioutil.WriteFile("/tmp/"+name, data, 0644)
	if err == nil {
		c.memoryCache[name] = data
	}
	return err
}

func (c *Cache) GetCache(name string) ([]byte, error) {
	if dat, ok := c.memoryCache[name]; ok {
		return dat, nil
	}

	f, err := os.Stat("/tmp/" + name)
	if os.IsNotExist(err) {
		return nil, err
	}

	if time.Since(f.ModTime()).Minutes() > float64(c.refreshTime) {
		return nil, fmt.Errorf("[GetCache] %s", "cache invalidate")
	}

	dat, err := ioutil.ReadFile("/tmp/" + name)
	if err != nil {
		return nil, err
	}
	return dat, nil
}
