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
	PlayerId  string
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
			log.Print("Message read error:", err)
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
	gameState := cl.GameState.GetSerializable()
	message := NewGameStateMessage(gameState)
	SerializeAndScheduleServerMessage(message, cl.Send)
	j, _ := json.Marshal(gameState)
	LogJsonForce("Sending game state:", j)
}

func (cl *Client) HandlePlayerEnter(mData map[string]interface{}) {
	pData := mData["player"].(map[string]interface{})
	player := NewPlayerFromMap(pData, cl.Ws)
	cl.PlayerId = player.Id
	cl.GameState.Players.Set(player.Id, player)
	message := NewPlayerEnterMessage(player)
	SerializeAndScheduleServerMessage(message, cl.Hub.Broadcast)
	LogDataForce("Handling player enter:", player.Id)
}

func (cl *Client) HandlePlayerExit(mData map[string]interface{}) {
	playerData := cl.GameState.Players.Get(cl.PlayerId)
	if playerData != nil {
		player := playerData.(Player)
		cl.GameState.Players.Delete(player.Id)
		message := NewPlayerExitMessage(player.Id)
		SerializeAndScheduleServerMessage(message, cl.Hub.Broadcast)
		LogDataForce("Handling player exit:", player.Id)
	}
	cl.Hub.Remove <- cl
}

func (cl *Client) HandlePlayerPosition(mData map[string]interface{}) {
	posMap := mData["position"].(map[string]interface{})
	newPosition := Position{
		X: posMap["x"].(float64),
		Y: posMap["y"].(float64),
	}
	playerId := mData["playerId"].(string)
	player := cl.GameState.Players.Get(playerId).(Player)
	player.Position = newPosition
	cl.GameState.Players.Set(playerId, player)
	message := NewPlayerStateUpdateMessage(player)
	SerializeAndScheduleServerMessage(message, cl.Hub.Broadcast)
}

func (cl *Client) HandlePlayerEatFood(mData map[string]interface{}) {
	// player grows in size
	playerId := mData["playerId"].(string)
	player := cl.GameState.Players.Get(playerId).(Player)
	player.Size += 1
	cl.GameState.Players.Set(playerId, player)
	pmsg := NewPlayerStateUpdateMessage(player)
	SerializeAndScheduleServerMessage(pmsg, cl.Hub.Broadcast)
	// food position changes
	foodId := mData["foodId"].(string)
	food := cl.GameState.Foods.Get(foodId).(Food)
	food.Position = cl.GameState.GetNewSpawnFoodPosition()
	cl.GameState.Foods.Set(foodId, food)
	fmsg := NewFoodStateUpdateMessage(food)
	SerializeAndScheduleServerMessage(fmsg, cl.Hub.Broadcast)
}

func (cl *Client) HandlePlayerEatPlayer(mData map[string]interface{}) {
	// load players
	playerId := mData["playerId"].(string)
	player := cl.GameState.Players.Get(playerId).(Player)
	otherPlayerId := mData["otherPlayerId"].(string)
	otherPlayer := cl.GameState.Players.Get(otherPlayerId).(Player)
	// player who ate other grows in size
	player.Size += otherPlayer.Size
	// player who got eaten respawns with reset size to 1
	otherPlayer.Size = 1
	otherPlayer.Position = cl.GameState.GetNewSpawnPlayerPosition()
	// save players
	cl.GameState.Players.Set(playerId, player)
	cl.GameState.Players.Set(otherPlayerId, otherPlayer)
	// broadcast player update messages
	pmsg := NewPlayerStateUpdateMessage(player)
	opmsg := NewPlayerStateUpdateMessage(otherPlayer)
	SerializeAndScheduleServerMessage(pmsg, cl.Hub.Broadcast)
	SerializeAndScheduleServerMessage(opmsg, cl.Hub.Broadcast)
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
