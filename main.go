package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
	"wsjtx-web/pkg/wsjtx"

	"github.com/gorilla/websocket"
)

//go:embed dist
var distFS embed.FS

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

type Hub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mu        sync.Mutex
}

func newHub() *Hub {
	return &Hub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

func (h *Hub) run() {
	for {
		message := <-h.broadcast
		h.mu.Lock()
		for client := range h.clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Websocket error: %v", err)
				client.Close()
				delete(h.clients, client)
			}
		}
		h.mu.Unlock()
	}
}

func main() {
	// UDP Server
	addr, err := net.ResolveUDPAddr("udp", ":2237")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("Listening for WSJT-X on %s...", addr)

	// WebSocket Hub
	hub := newHub()
	go hub.run()

	// HTTP Server
	// Strip the "dist" prefix from the embedded filesystem
	fsys, err := fs.Sub(distFS, "dist")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.FS(fsys)))
	// UDP Loop
	var lastAddr *net.UDPAddr
	var addrMu sync.Mutex

	// Handle WebSocket messages
	go func() {
		for {
			// We need a way to receive messages from clients and send them to UDP.
			// The current Hub implementation only broadcasts FROM UDP TO Clients.
			// We need to change how we handle client messages.
			// Let's iterate over clients and read? No, that's blocking.
			// The client reading loop is already there:
			/*
				for {
					_, _, err := ws.ReadMessage()
					if err != nil { ... break }
				}
			*/
			// We should change that loop to actually parse the message.
			time.Sleep(1 * time.Second)
		}
	}()

	// We need to modify the WebSocket handler to process incoming messages
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		defer ws.Close()

		hub.mu.Lock()
		hub.clients[ws] = true
		hub.mu.Unlock()

		// Keep connection alive and read messages
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				hub.mu.Lock()
				delete(hub.clients, ws)
				hub.mu.Unlock()
				break
			}

			// Parse command from frontend
			var cmd struct {
				Type string          `json:"type"`
				Data json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(message, &cmd); err != nil {
				log.Printf("Error parsing command: %v", err)
				continue
			}

			addrMu.Lock()
			targetAddr := lastAddr
			addrMu.Unlock()

			if targetAddr == nil {
				log.Println("No WSJT-X instance seen yet, cannot send command")
				continue
			}

			var msg wsjtx.Message
			switch cmd.Type {
			case "reply":
				var d struct {
					Time      uint32  `json:"Time"`
					SNR       int32   `json:"SNR"`
					DeltaTime float64 `json:"DeltaTime"`
					DeltaFreq uint32  `json:"DeltaFrequency"`
					Mode      string  `json:"Mode"`
					Message   string  `json:"Message"`
					LowConf   bool    `json:"LowConfidence"`
					Modifiers uint8   `json:"Modifiers"`
				}
				if err := json.Unmarshal(cmd.Data, &d); err != nil {
					log.Printf("Error parsing reply data: %v", err)
					continue
				}
				msg = &wsjtx.ReplyMessage{
					BaseMessage: wsjtx.BaseMessage{Id: "WSJTX-WEB"},
					Time:        d.Time,
					SNR:         d.SNR,
					DeltaTime:   d.DeltaTime,
					DeltaFreq:   d.DeltaFreq,
					Mode:        d.Mode,
					Message:     d.Message,
					LowConf:     d.LowConf,
					Modifiers:   d.Modifiers,
				}
			case "halt":
				var d struct {
					AutoTxOnly bool `json:"AutoTxOnly"`
				}
				if err := json.Unmarshal(cmd.Data, &d); err != nil {
					log.Printf("Error parsing halt data: %v", err)
					continue
				}
				msg = &wsjtx.HaltTxMessage{
					BaseMessage: wsjtx.BaseMessage{Id: "WSJTX-WEB"},
					AutoTxOnly:  d.AutoTxOnly,
				}
			default:
				log.Printf("Unknown command type: %s", cmd.Type)
				continue
			}

			// Encode and send
			var buf bytes.Buffer
			encoder := wsjtx.NewEncoder(&buf)
			if err := encoder.Encode(msg); err != nil {
				log.Printf("Error encoding message: %v", err)
				continue
			}

			if _, err := conn.WriteToUDP(buf.Bytes(), targetAddr); err != nil {
				log.Printf("Error sending UDP: %v", err)
			} else {
				log.Printf("Sent %s to %s", cmd.Type, targetAddr)
			}
		}
	})

	go func() {
		log.Println("Starting HTTP server on :8080...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// UDP Loop
	buf := make([]byte, 4096)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error reading UDP:", err)
			continue
		}

		addrMu.Lock()
		lastAddr = remoteAddr
		addrMu.Unlock()

		msg, err := wsjtx.ParsePacket(buf[:n])
		if err != nil {
			continue
		}

		// Broadcast to WebSockets
		wrapper := struct {
			Type string      `json:"type"`
			Data interface{} `json:"data"`
		}{
			Type: fmt.Sprintf("%d", msg.Type()),
			Data: msg,
		}

		wrappedJson, err := json.Marshal(wrapper)
		if err == nil {
			hub.broadcast <- wrappedJson
		}

		switch m := msg.(type) {
		case *wsjtx.DecodeMessage:
			fmt.Printf("Decode: %s %s SNR=%d %d\n", m.Mode, m.Message, m.SNR, m.Time)
		}
	}
}
