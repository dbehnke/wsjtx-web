package server

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"log"
	"net"
	"net/http"
	"sync"

	"wsjtx-web/internal/audio"
	"wsjtx-web/internal/hub"
	"wsjtx-web/pkg/wsjtx"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

// AddrStore remembers the last WSJT-X UDP peer for outbound commands.
type AddrStore struct {
	mu   sync.Mutex
	addr *net.UDPAddr
}

// Set records the latest remote UDP address.
func (s *AddrStore) Set(addr *net.UDPAddr) {
	s.mu.Lock()
	s.addr = addr
	s.mu.Unlock()
}

// Get returns the latest remote UDP address, or nil if none yet.
func (s *AddrStore) Get() *net.UDPAddr {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.addr
}

// Config holds dependencies for the HTTP/WebSocket server.
type Config struct {
	Addr     string
	StaticFS fs.FS
	Hub      *hub.Hub
	UDPConn  *net.UDPConn
	Addrs    *AddrStore
	Audio    *audio.Manager // nil when audio is disabled
}

// Start registers routes and begins ListenAndServe (blocking).
func Start(cfg Config) error {
	http.Handle("/", http.FileServer(http.FS(cfg.StaticFS)))
	http.HandleFunc("/ws", cfg.handleWS)
	log.Printf("Starting HTTP server on %s...", cfg.Addr)
	return http.ListenAndServe(cfg.Addr, nil)
}

func (cfg Config) handleWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	cfg.Hub.Register(ws)
	defer cfg.Hub.Unregister(ws)

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		var cmd struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(message, &cmd); err != nil {
			log.Printf("Error parsing command: %v", err)
			continue
		}

		switch cmd.Type {
		case "getAudioDevices":
			cfg.handleGetAudioDevices()
			continue
		case "setAudioDevice":
			cfg.handleSetAudioDevice(cmd.Data)
			continue
		}

		targetAddr := cfg.Addrs.Get()
		if targetAddr == nil {
			log.Println("No WSJT-X instance seen yet, cannot send command")
			continue
		}

		msg, ok := buildUDPMessage(cmd.Type, cmd.Data)
		if !ok {
			continue
		}

		var buf bytes.Buffer
		encoder := wsjtx.NewEncoder(&buf)
		if err := encoder.Encode(msg); err != nil {
			log.Printf("Error encoding message: %v", err)
			continue
		}

		if _, err := cfg.UDPConn.WriteToUDP(buf.Bytes(), targetAddr); err != nil {
			log.Printf("Error sending UDP: %v", err)
		} else {
			log.Printf("Sent %s to %s", cmd.Type, targetAddr)
		}
	}
}

func (cfg Config) handleGetAudioDevices() {
	if cfg.Audio == nil {
		log.Println("Audio disabled; ignoring getAudioDevices")
		resp, _ := json.Marshal(struct {
			Type string      `json:"type"`
			Data interface{} `json:"data"`
		}{
			Type: "audioDevices",
			Data: []map[string]string{},
		})
		cfg.Hub.Broadcast(hub.BroadcastMessage{
			Type:    websocket.TextMessage,
			Payload: resp,
		})
		return
	}

	devices, err := cfg.Audio.ListDevices()
	if err != nil {
		log.Printf("Error listing devices: %v", err)
		return
	}
	resp, _ := json.Marshal(struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		Type: "audioDevices",
		Data: devices,
	})
	cfg.Hub.Broadcast(hub.BroadcastMessage{
		Type:    websocket.TextMessage,
		Payload: resp,
	})
}

func (cfg Config) handleSetAudioDevice(data json.RawMessage) {
	if cfg.Audio == nil {
		log.Println("Audio disabled; ignoring setAudioDevice")
		return
	}
	var d struct {
		Id string `json:"id"`
	}
	if err := json.Unmarshal(data, &d); err != nil {
		log.Printf("Error parsing setAudioDevice data: %v", err)
		return
	}
	if err := cfg.Audio.Start(d.Id); err != nil {
		log.Printf("Error setting audio device: %v", err)
	}
}

func buildUDPMessage(cmdType string, data json.RawMessage) (wsjtx.Message, bool) {
	switch cmdType {
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
		if err := json.Unmarshal(data, &d); err != nil {
			log.Printf("Error parsing reply data: %v", err)
			return nil, false
		}
		return &wsjtx.ReplyMessage{
			BaseMessage: wsjtx.BaseMessage{Id: "WSJTX-WEB"},
			Time:        d.Time,
			SNR:         d.SNR,
			DeltaTime:   d.DeltaTime,
			DeltaFreq:   d.DeltaFreq,
			Mode:        d.Mode,
			Message:     d.Message,
			LowConf:     d.LowConf,
			Modifiers:   d.Modifiers,
		}, true
	case "halt":
		var d struct {
			AutoTxOnly bool `json:"AutoTxOnly"`
		}
		if err := json.Unmarshal(data, &d); err != nil {
			log.Printf("Error parsing halt data: %v", err)
			return nil, false
		}
		return &wsjtx.HaltTxMessage{
			BaseMessage: wsjtx.BaseMessage{Id: "WSJTX-WEB"},
			AutoTxOnly:  d.AutoTxOnly,
		}, true
	default:
		log.Printf("Unknown command type: %s", cmdType)
		return nil, false
	}
}
