package handler

import (
	"encoding/base64"
	"encoding/json"
	"runtime"
	"sync"
	"time"

	"github.com/ghophp/buildbot-dashboard/container"

	"github.com/beatrichartz/martini-sockets"
	"github.com/go-martini/martini"
)

const (
	CloseNormalClosure = 1000
)

type (
	ClientList struct {
		sync.Mutex
		clients []*Client
	}

	Client struct {
		in         <-chan *Message
		out        chan<- *Message
		done       <-chan bool
		err        <-chan error
		disconnect chan<- int
	}

	HttpResponse struct {
		res string
		err error
	}

	Message struct {
		Text string `json:"text"`
	}
)

// clientList hold the connection with the websockets
var clientList *ClientList

func newClientList() *ClientList {
	return &ClientList{sync.Mutex{}, make([]*Client, 0)}
}

func (r *ClientList) appendClient(client *Client) {
	r.Lock()
	r.clients = append(r.clients, client)
	r.Unlock()
}

func (r *ClientList) removeClient(client *Client) {
	r.Lock()
	defer r.Unlock()

	for index, c := range r.clients {
		if c == client {
			r.clients = append(r.clients[:index], r.clients[(index+1):]...)
		}
	}
}

func (r *ClientList) broadcast(msg *Message) {
	r.Lock()
	for _, c := range r.clients {
		c.out <- msg
	}
	defer r.Unlock()
}

func asyncBuilderFetch(c *container.ContainerBag, builders map[string]Builder) []*HttpResponse {
	ch := make(chan *HttpResponse)
	responses := []*HttpResponse{}

	for id, builder := range builders {
		go func(id string, builder Builder) {
			var response string

			b, err := GetBuilder(c, id, builder)
			if err == nil {
				if r, err := json.Marshal(b); err == nil {
					response = base64.StdEncoding.EncodeToString(r)
				}
			}

			ch <- &HttpResponse{response, err}
		}(id, builder)
	}

	for {
		select {
		case r := <-ch:
			responses = append(responses, r)
			if len(responses) == len(builders) {
				return responses
			}
		}
	}

	return responses
}

func monitorBuilders(c *container.ContainerBag) {
	for {
		if len(clientList.clients) > 0 {
			builders, err := GetBuilders(c, false)
			if err != nil || len(builders) <= 0 {
				continue
			}

			results := asyncBuilderFetch(c, builders)
			for _, result := range results {
				if len(result.res) > 0 {
					clientList.broadcast(&Message{result.res})
				}
			}
		}

		c.Logger.Printf("fetching builders at %d goroutines", runtime.NumGoroutine())
		time.Sleep(time.Second * time.Duration(c.RefreshSec))
	}

	c.Logger.Printf("monitorBuilders left the routine loop")
}

func AddWs(m *martini.ClassicMartini, c *container.ContainerBag) {
	clientList = newClientList()

	go monitorBuilders(c)

	m.Get("/ws", sockets.JSON(Message{}, &sockets.Options{WriteWait: 0, PongWait: 0}),
		func(params martini.Params,
			receiver <-chan *Message,
			sender chan<- *Message,
			done <-chan bool,
			disconnect chan<- int,
			err <-chan error) (int, string) {

			client := &Client{receiver, sender, done, err, disconnect}
			clientList.appendClient(client)

			for {
				select {
				case <-client.done:
					clientList.removeClient(client)
					return 200, "OK"
				}
			}
		})
}
