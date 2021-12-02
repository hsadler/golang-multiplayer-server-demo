package main

import "github.com/gorilla/websocket"

//////////////// MAIN MODELS ////////////////

// game state
type GameState struct {
	Players                  map[string]*Player `json:"players"`
	Foods                    map[string]*Food   `json:"foods"`
	Mines                    map[string]*Mine   `json:"mines"`
	RoundHistory             map[string]*Round  `json:"roundHistory"`
	RoundCurrent             *Round             `json:"roundCurrent"`
	SecondsToCurrentRoundEnd int                `json:"secondsToCurrentRoundEnd"`
	SecondsToNextRoundStart  int                `json:"secondsToNextRoundStart"`
}

// round
type Round struct {
	Id              string         `json:"id"`
	PlayerIdToScore map[string]int `json:"playerIdToScore"`
	TimeStart       int            `json:"timeStart"`
	TimeEnd         int            `json:"timeEnd"`
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
