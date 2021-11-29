package main

import "github.com/gorilla/websocket"

type Player struct {
	Id       string    `json:"id"`
	Position *Position `json:"position"`
}

func NewPlayerFromMap(pData map[string]interface{}, ws *websocket.Conn) *Player {
	posMap := pData["position"].(map[string]interface{})
	pos := Position{
		X: posMap["x"].(float64),
		Y: posMap["y"].(float64),
	}
	player := Player{
		Id:       pData["id"].(string),
		Position: &pos,
	}
	return &player
}
