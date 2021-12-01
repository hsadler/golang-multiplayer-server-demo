package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub    *Hub
	Ws     *websocket.Conn
	Player *Player
	Send   chan []byte
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
		// fmt.Println("client message received:")
		// ConsoleLogJsonByteArray(message)
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
	allPlayers := make([]*Player, 0)
	for client := range cl.Hub.Clients {
		if client.Player != nil {
			allPlayers = append(allPlayers, client.Player)
		}
	}
	messageData := GameStateMessage{
		MessageType: SERVER_MESSAGE_TYPE_GAME_STATE,
		GameState: &GameStateJsonSerializable{
			Players: allPlayers,
		},
	}
	serialized, _ := json.Marshal(messageData)
	cl.Send <- serialized
}

func (cl *Client) HandlePlayerEnter(mData map[string]interface{}) {
	player := NewPlayerFromMap(mData["player"].(map[string]interface{}), cl.Ws)
	cl.Player = player
	message := PlayerEnterMessage{
		MessageType: SERVER_MESSAGE_TYPE_PLAYER_ENTER,
		Player:      player,
	}
	serialized, _ := json.Marshal(message)
	cl.Hub.Broadcast <- serialized
}

func (cl *Client) HandlePlayerExit(mData map[string]interface{}) {
	message := PlayerExitMessage{
		MessageType: SERVER_MESSAGE_TYPE_PLAYER_EXIT,
		PlayerId:    cl.Player.Id,
	}
	serialized, _ := json.Marshal(message)
	cl.Hub.Broadcast <- serialized
	cl.Hub.Remove <- cl
}

func (cl *Client) HandlePlayerPosition(mData map[string]interface{}) {
	player := NewPlayerFromMap(mData["player"].(map[string]interface{}), cl.Ws)
	cl.Player = player
	message := PlayerStateUpdateMessage{
		MessageType: SERVER_MESSAGE_TYPE_PLAYER_STATE_UPDATE,
		Player:      player,
	}
	serialized, _ := json.Marshal(message)
	cl.Hub.Broadcast <- serialized
}

func (cl *Client) HandlePlayerEatFood(mData map[string]interface{}) {
	// STUB
}

func (cl *Client) HandlePlayerEatPlayer(mData map[string]interface{}) {
	// STUB
}

func (cl *Client) HandleMineDamagePlayer(mData map[string]interface{}) {
	// STUB
}

func (cl *Client) SendMessages() {
	for {
		select {
		case message, ok := <-cl.Send:
			if !ok {
				return
			}
			SendJsonMessage(cl.Ws, message)
		}
	}
}

func (cl *Client) Cleanup() {
	close(cl.Send)
}
