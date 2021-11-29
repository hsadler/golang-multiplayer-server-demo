package main

import (
	"github.com/gorilla/websocket"
)

func SendJsonMessage(ws *websocket.Conn, messageJson []byte) {
	ws.WriteMessage(1, messageJson)
	// log that message was sent
	// fmt.Println("server message sent:")
	// ConsoleLogJsonByteArray(messageJson)
}

type PlayerMessage struct {
	MessageType string  `json:"messageType"`
	Player      *Player `json:"player"`
}

type GameStateJsonSerializable struct {
	Players []*Player `json:"players"`
}

type GameStateMessage struct {
	MessageType string                     `json:"messageType"`
	GameState   *GameStateJsonSerializable `json:"gameState"`
}
