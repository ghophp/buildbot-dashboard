package handler

import (
	"encoding/base64"
	"encoding/json"
	"sync"
	"time"

	"github.com/ghophp/buildbot-dashboard/container"

	"github.com/beatrichartz/martini-sockets"
	"github.com/go-martini/martini"
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

func monitorBuilders(c *container.ContainerBag) {
	responses := make(chan string)

	for {
		if len(clientList.clients) > 0 {
			builders, err := GetBuilders(c)
			if err != nil || len(builders) <= 0 {
				return
			}

			filtered := make(map[string]Builder)
			for id, builder := range builders {
				if len(builder.CachedBuilds) > 0 {
					filtered[id] = builder
				}
			}

			var wg sync.WaitGroup
			wg.Add(len(filtered))

			for id, builder := range filtered {
				go func(id string, builder Builder) {
					defer wg.Done()

					b, err := GetBuilder(c, id, builder)
					if err == nil {
						if r, err := json.Marshal(b); err == nil {
							responses <- base64.StdEncoding.EncodeToString(r)
						}
					}

				}(id, builder)
			}

			go func() {
				for response := range responses {
					clientList.broadcast(&Message{response})
				}
			}()

			wg.Wait()
		}

		time.Sleep(time.Second * time.Duration(c.RefreshSec))
	}
}

func AddWs(m *martini.ClassicMartini, c *container.ContainerBag) {
	clientList = newClientList()

	go monitorBuilders(c)

	m.Get("/ws", sockets.JSON(Message{}),
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
