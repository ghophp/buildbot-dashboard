package cache

import (
	"io/ioutil"
	"os"
)

type Cache struct {
	memoryCache map[string][]byte
}

func NewCache() *Cache {
	return &Cache{
		memoryCache: make(map[string][]byte),
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
	if _, err := os.Stat("/tmp/" + name); os.IsNotExist(err) {
		return nil, err
	}
	dat, err := ioutil.ReadFile("/tmp/" + name)
	if err != nil {
		return nil, err
	}
	return dat, nil
}
