package main

import "github.com/gorilla/websocket"

//////////////// MAIN MODELS ////////////////

// round
type PlayerScore struct {
	PlayerId string `json:"playerId"`
	Score    int    `json:"score"`
}
type Round struct {
	PlayerScores []PlayerScore `json:"playerScores"`
}

// player
type Player struct {
	Id               string   `json:"id"`
	Active           bool     `json:"active"`
	Name             string   `json:"name"`
	Position         Position `json:"position"`
	Size             int      `json:"size"`
	TimeUntilRespawn int      `json:"timeUntilRespawn"`
}

func NewPlayerFromMap(pData map[string]interface{}, ws *websocket.Conn) Player {
	posMap := pData["position"].(map[string]interface{})
	pos := Position{
		X: posMap["x"].(float64),
		Y: posMap["y"].(float64),
	}
	player := Player{
		Id:               pData["id"].(string),
		Active:           pData["active"].(bool),
		Name:             pData["name"].(string),
		Position:         pos,
		Size:             int(pData["size"].(float64)),
		TimeUntilRespawn: 0,
	}
	return player
}

// food
type Food struct {
	Id       string   `json:"id"`
	Active   bool     `json:"active"`
	Position Position `json:"position"`
	Size     int      `json:"size"`
}

// mine
type Mine struct {
	Id       string   `json:"id"`
	Active   bool     `json:"active"`
	Position Position `json:"position"`
	Size     int      `json:"size"`
}

//////////////// SUB MODELS ////////////////

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
