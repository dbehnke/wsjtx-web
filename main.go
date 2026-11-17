package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"

	"wsjtx-web/internal/audio"
	"wsjtx-web/internal/hub"
	"wsjtx-web/internal/server"
	"wsjtx-web/pkg/wsjtx"

	"github.com/gorilla/websocket"
)

//go:embed dist
var distFS embed.FS

func main() {
	udpAddr := flag.String("udp", ":2237", "UDP address to receive WSJT-X traffic")
	httpAddr := flag.String("http", ":8080", "HTTP/WebSocket listen address")
	enableAudio := flag.Bool("audio", false, "Enable audio capture for waterfall (malgo/CGO)")
	flag.Parse()

	addr, err := net.ResolveUDPAddr("udp", *udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("listen udp %s: %v (is another app like GridTracker already bound to this port? use -udp to pick another)", addr, err)
	}
	defer conn.Close()

	log.Printf("Listening for WSJT-X on %s...", addr)

	h := hub.New()
	go h.Run()

	fsys, err := fs.Sub(distFS, "dist")
	if err != nil {
		log.Fatal(err)
	}

	addrs := &server.AddrStore{}

	var audioMgr *audio.Manager
	if *enableAudio {
		audioMgr, err = audio.New(h)
		if err != nil {
			log.Println("Failed to init audio manager:", err)
		} else if err := audioMgr.Start(""); err != nil {
			log.Println("Failed to start audio:", err)
		}
	} else {
		log.Println("Audio capture disabled (pass -audio to enable waterfall capture)")
	}

	go func() {
		if err := server.Start(server.Config{
			Addr:     *httpAddr,
			StaticFS: fsys,
			Hub:      h,
			UDPConn:  conn,
			Addrs:    addrs,
			Audio:    audioMgr,
		}); err != nil {
			log.Fatal(err)
		}
	}()

	buf := make([]byte, 4096)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error reading UDP:", err)
			continue
		}

		addrs.Set(remoteAddr)

		msg, err := wsjtx.ParsePacket(buf[:n])
		if err != nil {
			continue
		}

		wrapper := struct {
			Type string      `json:"type"`
			Data interface{} `json:"data"`
		}{
			Type: fmt.Sprintf("%d", msg.Type()),
			Data: msg,
		}

		wrappedJSON, err := json.Marshal(wrapper)
		if err == nil {
			h.Broadcast(hub.BroadcastMessage{
				Type:    websocket.TextMessage,
				Payload: wrappedJSON,
			})
		}
	}
}
