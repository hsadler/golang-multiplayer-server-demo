package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	// create hub and run the channel listeners
	h := &Hub{
		Clients:   make(map[*Client]bool),
		Add:       make(chan *Client),
		Remove:    make(chan *Client),
		Broadcast: make(chan []byte),
	}
	go h.RunListeners()
	// create game-state and run listeners for write channels
	gs := &GameState{
		MapHeight:         MAP_HEIGHT,
		MapWidth:          MAP_WIDTH,
		Players:           make(map[string]Player),
		Foods:             make(map[string]Food),
		Mines:             make(map[string]Mine),
		AddPlayer:         make(chan Player),
		RemovePlayer:      make(chan Player),
		UpdatePlayerState: make(chan Player),
		UpdateFoodState:   make(chan Food),
		UpdateMineState:   make(chan Mine),
		Mu:                &sync.RWMutex{},
	}
	go gs.RunWriteListeners()
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
