package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	// create hub singleton and run
	h := NewHub()
	go h.Run()
	// create and run game-state-manager singleton and run
	gsm := NewGameStateManager()
	go gsm.RunRoundTicker()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{} // use default options
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		// create client, run processes, and add to hub
		cl := &Client{
			Hub:       h,
			GameState: gsm.GameState,
			Ws:        ws,
			Player:    nil,
			Send:      make(chan []byte, 256),
		}
		go cl.RecieveMessages()
		go cl.SendMessages()
		h.Add <- cl
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	addr := flag.String("addr", "0.0.0.0:"+port, "http service address")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
