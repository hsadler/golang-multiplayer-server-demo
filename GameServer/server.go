package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	log.SetFlags(log.LstdFlags)
	// create hub and run the channel listeners
	h := &Hub{
		Clients:   make(map[*Client]bool),
		Add:       make(chan *Client),
		Remove:    make(chan *Client),
		Broadcast: make(chan []byte),
	}
	go h.RunListeners()
	// create game-state
	gs := &GameState{
		MapHeight: MAP_HEIGHT,
		MapWidth:  MAP_WIDTH,
		Players:   NewCMap(),
		Foods:     NewCMap(),
		Mines:     NewCMap(),
	}
	// create round-manager and run the round management process
	rm := &RoundManager{
		Hub:                      h,
		GameState:                gs,
		RoundIsInProgress:        false,
		SecondsToCurrentRoundEnd: 0,
		SecondsToNextRoundStart:  SECONDS_BETWEEN_ROUNDS,
	}
	go rm.RunRoundTicker()
	// handle client connections
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// upgrade request to websocket and use default options
		upgrader := websocket.Upgrader{}
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("Request upgrade error:", err)
			return
		}
		// create client, run processes, and add to hub
		cl := &Client{
			Hub:       h,
			GameState: gs,
			Ws:        ws,
			PlayerId:  "",
			Send:      make(chan []byte, 256),
		}
		go cl.RecieveMessages()
		go cl.SendMessages()
		h.Add <- cl
	})
	// run the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	addr := flag.String("addr", "0.0.0.0:"+port, "http service address")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("Server start error:", err)
	}
}
