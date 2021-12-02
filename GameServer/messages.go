package main

import (
	"github.com/gorilla/websocket"
)

const CLIENT_MESSAGE_TYPE_PLAYER_ENTER string = "CLIENT_MESSAGE_TYPE_PLAYER_ENTER"
const CLIENT_MESSAGE_TYPE_PLAYER_EXIT string = "CLIENT_MESSAGE_TYPE_PLAYER_EXIT"
const CLIENT_MESSAGE_TYPE_PLAYER_POSITION string = "CLIENT_MESSAGE_TYPE_PLAYER_POSITION"
const CLIENT_MESSAGE_TYPE_PLAYER_EAT_FOOD string = "CLIENT_MESSAGE_TYPE_PLAYER_EAT_FOOD"
const CLIENT_MESSAGE_TYPE_PLAYER_EAT_PLAYER string = "CLIENT_MESSAGE_TYPE_PLAYER_EAT_PLAYER"
const CLIENT_MESSAGE_TYPE_MINE_DAMAGE_PLAYER string = "CLIENT_MESSAGE_TYPE_MINE_DAMAGE_PLAYER"

const SERVER_MESSAGE_TYPE_GAME_STATE string = "SERVER_MESSAGE_TYPE_GAME_STATE"
const SERVER_MESSAGE_TYPE_PLAYER_ENTER string = "SERVER_MESSAGE_TYPE_PLAYER_ENTER"
const SERVER_MESSAGE_TYPE_PLAYER_EXIT string = "SERVER_MESSAGE_TYPE_PLAYER_EXIT"
const SERVER_MESSAGE_TYPE_PLAYER_STATE_UPDATE string = "SERVER_MESSAGE_TYPE_PLAYER_STATE_UPDATE"
const SERVER_MESSAGE_TYPE_FOOD_STATE_UPDATE string = "SERVER_MESSAGE_TYPE_FOOD_STATE_UPDATE"
const SERVER_MESSAGE_TYPE_MINE_STATE_UPDATE string = "SERVER_MESSAGE_TYPE_MINE_STATE_UPDATE"
const SERVER_MESSAGE_TYPE_ROUND_TIME_TO_START string = "SERVER_MESSAGE_TYPE_ROUND_TIME_TO_START"
const SERVER_MESSAGE_TYPE_ROUND_START string = "SERVER_MESSAGE_TYPE_ROUND_START"
const SERVER_MESSAGE_TYPE_ROUND_END string = "SERVER_MESSAGE_TYPE_ROUND_END"

func SendJsonMessage(ws *websocket.Conn, messageJson []byte) {
	ws.WriteMessage(1, messageJson)
	// log that message was sent
	// fmt.Println("server message sent:")
	// ConsoleLogJsonByteArray(messageJson)
}

// TODO: finalize server message schemas

// TODO: deprecate GameStateJsonSerializable in favor of game model struct
type GameStateJsonSerializable struct {
	Players []*Player `json:"players"`
}
type GameStateMessage struct {
	MessageType string                     `json:"messageType"`
	GameState   *GameStateJsonSerializable `json:"gameState"`
}

type PlayerEnterMessage struct {
	MessageType string  `json:"messageType"`
	Player      *Player `json:"player"`
}

type PlayerExitMessage struct {
	MessageType string `json:"messageType"`
	PlayerId    string `json:"playerId"`
}

type PlayerStateUpdateMessage struct {
	MessageType string  `json:"messageType"`
	Player      *Player `json:"player"`
}

type FoodStateUpdateMessage struct {
	MessageType string `json:"messageType"`
	Food        *Food  `json:"food"`
}

type MineStateUpdateMessage struct {
	MessageType string `json:"messageType"`
	Mine        *Mine  `json:"mine"`
}

type RoundTimeToStartMessage struct {
	MessageType string `json:"messageType"`
	Seconds     int    `json:"second"`
}

type RoundStartMessage struct {
	MessageType string `json:"messageType"`
}

type RoundEndMessage struct {
	MessageType string `json:"messageType"`
	Round       *Round `json:"round"`
}
