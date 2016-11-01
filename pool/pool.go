package pool

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

const RequestError = "RequestError"

type (
	Pool interface {
		Fetch(string, chan<- string)
	}

	RequestPool struct {
		Mutex       *sync.Mutex
		CurrentPool map[string][]chan<- string
	}
)

func NewRequestPool() *RequestPool {
	return &RequestPool{
		CurrentPool: map[string][]chan<- string{},
		Mutex:       &sync.Mutex{},
	}
}

func (r *RequestPool) Fetch(url string, listener chan<- string) {
	r.Mutex.Lock()
	r.CurrentPool[url] = append(r.CurrentPool[url], listener)
	r.Mutex.Unlock()

	if len(r.CurrentPool[url]) <= 1 {
		go r.processUrl(url)
	}
}

func (r *RequestPool) broadcast(url, result string) {
	for _, c := range r.CurrentPool[url] {
		c <- result
		close(c)
	}

	r.Mutex.Lock()
	r.CurrentPool[url] = []chan<- string{}
	r.Mutex.Unlock()
}

func (r *RequestPool) processUrl(url string) {
	req, err := http.Get(url)
	if err != nil {
		r.broadcast(url, fmt.Sprintf("[%s] %s", RequestError, err.Error()))
		return
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		r.broadcast(url, fmt.Sprintf("[%s] %s", RequestError, err.Error()))
		return
	}

	req.Body.Close()
	r.broadcast(url, string(b))
}
