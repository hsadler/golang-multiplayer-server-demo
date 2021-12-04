package main

import "github.com/gorilla/websocket"

//////////////// MAIN MODELS ////////////////

// round
type Round struct {
	PlayerIdToScore map[string]int `json:"playerIdToScore"`
}

// player
type Player struct {
	Id       string    `json:"id"`
	Active   bool      `json:"active"`
	Name     string    `json:"name"`
	Position *Position `json:"position"`
	Size     int       `json:"size"`
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

// food
type Food struct {
	Id       string    `json:"id"`
	Active   bool      `json:"active"`
	Position *Position `json:"position"`
}

// mine
type Mine struct {
	Id       string    `json:"id"`
	Active   bool      `json:"active"`
	Position *Position `json:"position"`
}

//////////////// SUB MODELS ////////////////

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
