package svc

import (
	"encoding/json"
	"sync"
	"time"
)

type wsHub struct {
	// the mutex to protect connections
	connectionsMx sync.RWMutex

	// Registered connections.
	connections map[*connection]struct{}

	// Inbound messages from the connections.
	broadcast chan []byte

	logMx sync.RWMutex
	log   [][]byte
}

func newHub() *wsHub {
	h := &wsHub{
		connectionsMx: sync.RWMutex{},
		broadcast:     make(chan []byte),
		connections:   make(map[*connection]struct{}),
	}

	go func() {
		for {
			msg := <-h.broadcast
			h.connectionsMx.RLock()
			for c := range h.connections {
				select {
				case c.send <- msg:
				// stop trying to send to this connection after trying for 1 second.
				// if we have to stop, it means that a reader died so remove the connection also.
				case <-time.After(1 * time.Second):
					h.removeConnection(c)
				}
			}
			h.connectionsMx.RUnlock()
		}
	}()
	return h
}

func (h *wsHub) BroadcastReload() {
	reload := map[string]string{
		"action": "reload",
	}
	reloadBytes, _ := json.Marshal(reload)
	h.broadcast <- reloadBytes
}

func (h *wsHub) BroadcastChat(name, message string) {
	chatMessage := map[string]string{
		"action":  "chat",
		"name":    name,
		"message": message,
	}
	chatBytes, _ := json.Marshal(chatMessage)
	h.broadcast <- chatBytes
}

func (h *wsHub) addConnection(conn *connection) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()
	h.connections[conn] = struct{}{}
}

func (h *wsHub) removeConnection(conn *connection) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()
	if _, ok := h.connections[conn]; ok {
		delete(h.connections, conn)
		close(conn.send)
	}
}
