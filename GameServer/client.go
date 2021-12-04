package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub       *Hub
	GameState *GameState
	Ws        *websocket.Conn
	Player    *Player
	Send      chan []byte
}

func (cl *Client) RecieveMessages() {
	// initialize client's game state by sending server's entire state
	cl.SendGameState()
	// do player removal from game state and websocket close on disconnect
	defer func() {
		cl.HandlePlayerExit(nil)
		cl.Ws.Close()
	}()
	for {
		// read message
		_, message, err := cl.Ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		// log message received
		LogJson("client message received:", message)
		// route message to handler
		messageTypeToHandler := map[string]func(map[string]interface{}){
			CLIENT_MESSAGE_TYPE_PLAYER_ENTER:       cl.HandlePlayerEnter,
			CLIENT_MESSAGE_TYPE_PLAYER_EXIT:        cl.HandlePlayerExit,
			CLIENT_MESSAGE_TYPE_PLAYER_POSITION:    cl.HandlePlayerPosition,
			CLIENT_MESSAGE_TYPE_PLAYER_EAT_FOOD:    cl.HandlePlayerEatFood,
			CLIENT_MESSAGE_TYPE_PLAYER_EAT_PLAYER:  cl.HandlePlayerEatPlayer,
			CLIENT_MESSAGE_TYPE_MINE_DAMAGE_PLAYER: cl.HandleMineDamagePlayer,
		}
		var mData map[string]interface{}
		if err := json.Unmarshal(message, &mData); err != nil {
			panic(err)
		}
		// process message with handler
		messageTypeToHandler[mData["messageType"].(string)](mData)
	}
}

func (cl *Client) SendGameState() {
	message := GameStateMessage{
		MessageType: SERVER_MESSAGE_TYPE_GAME_STATE,
		GameState:   cl.GameState,
	}
	SerializeAndScheduleServerMessage(message, cl.Send)
}

func (cl *Client) HandlePlayerEnter(mData map[string]interface{}) {
	player := NewPlayerFromMap(mData["player"].(map[string]interface{}), cl.Ws)
	cl.Player = player
	message := PlayerEnterMessage{
		MessageType: SERVER_MESSAGE_TYPE_PLAYER_ENTER,
		Player:      player,
	}
	SerializeAndScheduleServerMessage(message, cl.Hub.Broadcast)
}

func (cl *Client) HandlePlayerExit(mData map[string]interface{}) {
	message := PlayerExitMessage{
		MessageType: SERVER_MESSAGE_TYPE_PLAYER_EXIT,
		PlayerId:    cl.Player.Id,
	}
	SerializeAndScheduleServerMessage(message, cl.Hub.Broadcast)
	cl.Hub.Remove <- cl
}

func (cl *Client) HandlePlayerPosition(mData map[string]interface{}) {
	player := NewPlayerFromMap(mData["player"].(map[string]interface{}), cl.Ws)
	cl.Player = player
	message := PlayerStateUpdateMessage{
		MessageType: SERVER_MESSAGE_TYPE_PLAYER_STATE_UPDATE,
		Player:      player,
	}
	SerializeAndScheduleServerMessage(message, cl.Hub.Broadcast)
}

func (cl *Client) HandlePlayerEatFood(mData map[string]interface{}) {
	// TODO: STUB
}

func (cl *Client) HandlePlayerEatPlayer(mData map[string]interface{}) {
	// TODO: STUB
}

func (cl *Client) HandleMineDamagePlayer(mData map[string]interface{}) {
	// TODO: STUB
}

func (cl *Client) SendMessages() {
	for message := range cl.Send {
		SendJsonMessage(cl.Ws, message)
	}
}

func (cl *Client) Cleanup() {
	close(cl.Send)
}
