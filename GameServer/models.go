package main

import "github.com/gorilla/websocket"

// TODO: finalize game entity schemas

//////////////// MAIN MODELS ////////////////

// game state
type GameState struct {
	Players      map[string]*Player `json:"players"`
	Foods        map[string]*Food   `json:"foods"`
	Mines        map[string]*Mine   `json:"mines"`
	RoundHistory map[string]*Round  `json:"roundHistory"`
	RoundCurrent *Round             `json:"roundCurrent"`
}

// round
type Round struct {
}

// player
type Player struct {
	Id       string    `json:"id"`
	Active   bool      `json:"active"`
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

// food
type Food struct {
	Id       string    `json:"id"`
	Position *Position `json:"position"`
}

// mine
type Mine struct {
	Id       string    `json:"id"`
	Position *Position `json:"position"`
}

//////////////// SUB MODELS ////////////////

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
