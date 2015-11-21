package handler

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/ghophp/buildbot-dashing/container"
	"github.com/gorilla/websocket"
)

type (
	WsHandler struct {
		c *container.ContainerBag
	}

	ClientConn struct {
		websocket *websocket.Conn
		clientIP  net.Addr
	}
)

var (
	ActiveClients = make(map[ClientConn]int)
	wsMutex       sync.RWMutex
)

func addClient(cc ClientConn) {
	wsMutex.Lock()
	ActiveClients[cc] = 0
	wsMutex.Unlock()
}

func deleteClient(cc ClientConn) {
	wsMutex.Lock()
	delete(ActiveClients, cc)
	wsMutex.Unlock()
}

func broadcastMessage(messageType int, message []byte) {
	wsMutex.RLock()
	for client, _ := range ActiveClients {
		client.websocket.WriteMessage(messageType, message)
	}
	wsMutex.RUnlock()
}

func MonitorBuilders(c *container.ContainerBag) {
	responses := make(chan string)

	for {
		if len(ActiveClients) > 0 {
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
							responses <- string(r)
						}
					}

				}(id, builder)
			}

			go func() {
				for response := range responses {
					broadcastMessage(websocket.TextMessage, []byte(response))
				}
			}()

			wg.Wait()
		}

		time.Sleep(time.Second * time.Duration(c.RefreshSec))
	}
}

func NewWsHandler(c *container.ContainerBag) *WsHandler {
	go MonitorBuilders(c)

	return &WsHandler{
		c: c,
	}
}

func (h WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	sockCli := ClientConn{ws, ws.RemoteAddr()}
	addClient(sockCli)
}
