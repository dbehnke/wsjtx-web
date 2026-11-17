package hub

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// BroadcastMessage is a framed WebSocket payload to send to all clients.
type BroadcastMessage struct {
	Type    int
	Payload []byte
}

// Hub tracks WebSocket clients and fans out broadcast messages.
type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan BroadcastMessage
	mu        sync.Mutex
}

// New creates a Hub ready for Run and Register.
func New() *Hub {
	return &Hub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan BroadcastMessage),
	}
}

// Run processes the broadcast queue until the process exits.
func (h *Hub) Run() {
	for {
		message := <-h.broadcast
		h.mu.Lock()
		for client := range h.clients {
			err := client.WriteMessage(message.Type, message.Payload)
			if err != nil {
				log.Printf("Websocket error: %v", err)
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mu.Unlock()
	}
}

// Broadcast enqueues a message for all connected clients.
func (h *Hub) Broadcast(msg BroadcastMessage) {
	h.broadcast <- msg
}

// Register adds a client connection.
func (h *Hub) Register(c *websocket.Conn) {
	h.mu.Lock()
	h.clients[c] = true
	h.mu.Unlock()
}

// Unregister removes a client connection.
func (h *Hub) Unregister(c *websocket.Conn) {
	h.mu.Lock()
	delete(h.clients, c)
	h.mu.Unlock()
}
