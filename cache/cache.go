package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const InternalCacheFolder string = ".bbd"

type (
	Cache interface {
		SetCache(string, []byte, int) error
		GetCache(string) ([]byte, error)
	}

	FileCache struct {
		path string
	}
)

func NewFileCache() *FileCache {
	cc := &FileCache{
		path: "/tmp/",
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

func (c *FileCache) GetPath() string {
	return c.path
}

func (c *FileCache) SetCache(name string, data []byte, ttl int) error {
	if len(name) <= 0 {
		return fmt.Errorf("invalid cache name")
	}

	data = []byte(string(data) + "|" + strconv.Itoa(ttl))
	return ioutil.WriteFile(c.path+name, data, 0777)
}

func (c *FileCache) GetCache(name string) ([]byte, error) {
	if len(name) <= 0 {
		return nil, fmt.Errorf("invalid cache name")
	}

	f, err := os.Stat(c.path + name)
	if os.IsNotExist(err) {
		return nil, err
	}

	dat, err := ioutil.ReadFile(c.path + name)
	if err != nil {
		return nil, err
	}

	content := strings.Split(string(dat), "|")
	if len(content) < 2 {
		return nil, fmt.Errorf("[GetCache] %s", "invalid cache content")
	}

	ttl, err := strconv.Atoi(content[1])
	if err != nil {
		return nil, err
	}

	if time.Since(f.ModTime()).Seconds() > float64(ttl) {
		return nil, fmt.Errorf("[GetCache] %s", "cache invalidate")
	}

	return []byte(content[0]), nil
}
