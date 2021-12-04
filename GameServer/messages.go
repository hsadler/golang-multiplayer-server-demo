package main

import (
	"encoding/json"

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
const SERVER_MESSAGE_TYPE_SECONDS_TO_NEXT_ROUND_START string = "SERVER_MESSAGE_TYPE_SECONDS_TO_NEXT_ROUND_START"
const SERVER_MESSAGE_TYPE_SECONDS_TO_CURRENT_ROUND_END string = "SERVER_MESSAGE_TYPE_SECONDS_TO_CURRENT_ROUND_END"
const SERVER_MESSAGE_TYPE_ROUND_RESULT string = "SERVER_MESSAGE_TYPE_ROUND_RESULT"

// message related functions

func SerializeAndScheduleServerMessage(message interface{}, ch chan []byte) {
	serialized, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	ch <- serialized
}

func SendJsonMessage(ws *websocket.Conn, messageJson []byte) {
	ws.WriteMessage(1, messageJson)
	// log that message was sent
	LogJson("server message sent:", messageJson)
}

// messages

type GameStateMessage struct {
	MessageType string     `json:"messageType"`
	GameState   *GameState `json:"gameState"`
}

func (m *GameStateMessage) NewGameStateMessage(gs *GameState) *GameStateMessage {
	return &GameStateMessage{
		MessageType: SERVER_MESSAGE_TYPE_GAME_STATE,
		GameState:   gs,
	}
}

type PlayerEnterMessage struct {
	MessageType string  `json:"messageType"`
	Player      *Player `json:"player"`
}

func (m *PlayerEnterMessage) NewPlayerEnterMessage(p *Player) *PlayerEnterMessage {
	return &PlayerEnterMessage{
		MessageType: SERVER_MESSAGE_TYPE_PLAYER_ENTER,
		Player:      p,
	}
}

type PlayerExitMessage struct {
	MessageType string `json:"messageType"`
	PlayerId    string `json:"playerId"`
}

func (m *PlayerExitMessage) NewPlayerExitMessage(pId string) *PlayerExitMessage {
	return &PlayerExitMessage{
		MessageType: SERVER_MESSAGE_TYPE_PLAYER_EXIT,
		PlayerId:    pId,
	}
}

type PlayerStateUpdateMessage struct {
	MessageType string  `json:"messageType"`
	Player      *Player `json:"player"`
}

func (m *PlayerStateUpdateMessage) NewPlayerStateUpdateMessage(p *Player) *PlayerStateUpdateMessage {
	return &PlayerStateUpdateMessage{
		MessageType: SERVER_MESSAGE_TYPE_PLAYER_STATE_UPDATE,
		Player:      p,
	}
}

type FoodStateUpdateMessage struct {
	MessageType string `json:"messageType"`
	Food        *Food  `json:"food"`
}

func (m *FoodStateUpdateMessage) NewFoodStateUpdateMessage(f *Food) *FoodStateUpdateMessage {
	return &FoodStateUpdateMessage{
		MessageType: SERVER_MESSAGE_TYPE_FOOD_STATE_UPDATE,
		Food:        f,
	}
}

type MineStateUpdateMessage struct {
	MessageType string `json:"messageType"`
	Mine        *Mine  `json:"mine"`
}

func (m *MineStateUpdateMessage) NewMineStateUpdateMessage(mine *Mine) *MineStateUpdateMessage {
	return &MineStateUpdateMessage{
		MessageType: SERVER_MESSAGE_TYPE_MINE_STATE_UPDATE,
		Mine:        mine,
	}
}

type SecondsToNextRoundStartMessage struct {
	MessageType string `json:"messageType"`
	Seconds     int    `json:"seconds"`
}

func (m *SecondsToNextRoundStartMessage) NewSecondsToNextRoundStartMessage(s int) *SecondsToNextRoundStartMessage {
	return &SecondsToNextRoundStartMessage{
		MessageType: SERVER_MESSAGE_TYPE_SECONDS_TO_NEXT_ROUND_START,
		Seconds:     s,
	}
}

type SecondsToCurrentRoundEndMessage struct {
	MessageType string `json:"messageType"`
	Seconds     int    `json:"seconds"`
}

func (m *SecondsToCurrentRoundEndMessage) NewSecondsToCurrentRoundEndMessage(s int) *SecondsToCurrentRoundEndMessage {
	return &SecondsToCurrentRoundEndMessage{
		MessageType: SERVER_MESSAGE_TYPE_SECONDS_TO_CURRENT_ROUND_END,
		Seconds:     s,
	}
}

type RoundResultMessage struct {
	MessageType string `json:"messageType"`
	Round       *Round `json:"round"`
}

func (m *RoundResultMessage) NewRoundResultMessage(r *Round) *RoundResultMessage {
	return &RoundResultMessage{
		MessageType: SERVER_MESSAGE_TYPE_ROUND_RESULT,
		Round:       r,
	}
}
