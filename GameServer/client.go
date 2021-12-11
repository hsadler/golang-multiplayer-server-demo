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
	// load serializable game state
	gameState := cl.GameState.GetSerializable()
	// send message to client
	msg := NewGameStateMessage(gameState)
	SerializeAndScheduleServerMessage(msg, cl.Send)
	// logging
	j, _ := json.Marshal(gameState)
	LogJsonForce("Sending game state:", j)
}

func (cl *Client) HandlePlayerEnter(mData map[string]interface{}) {
	// parse inputs
	pData := mData["player"].(map[string]interface{})
	player := NewPlayerFromMap(pData, cl.Ws)
	// set player id on client
	cl.PlayerId = player.Id
	// commit datastore insert
	cl.GameState.Players.Set(player.Id, player)
	// broadcast message
	message := NewPlayerEnterMessage(player)
	SerializeAndScheduleServerMessage(message, cl.Hub.Broadcast)
	LogDataForce("Handling player enter:", player.Id)
}

func (cl *Client) HandlePlayerExit(mData map[string]interface{}) {
	// parse inputs
	playerId := mData["playerId"].(string)
	// datastore load
	playerData := cl.GameState.Players.Get(playerId)
	if playerData != nil {
		player := playerData.(Player)
		// commit datastore delete
		cl.GameState.Players.Delete(player.Id)
		// broadcast message
		message := NewPlayerExitMessage(player.Id)
		SerializeAndScheduleServerMessage(message, cl.Hub.Broadcast)
		LogDataForce("Handling player exit:", player.Id)
	}
	// schedule client removal from hub
	cl.Hub.Remove <- cl
}

func (cl *Client) HandlePlayerPosition(mData map[string]interface{}) {
	// parse inputs
	playerId := mData["playerId"].(string)
	posMap := mData["position"].(map[string]interface{})
	newPosition := Position{
		X: posMap["x"].(float64),
		Y: posMap["y"].(float64),
	}
	// datastore load
	player := cl.GameState.Players.Get(playerId).(Player)
	// player position update
	player.Position = newPosition
	// datastore save
	cl.GameState.Players.Set(playerId, player)
	// broadcast entity update message
	msg := NewPlayerStateUpdateMessage(player)
	SerializeAndScheduleServerMessage(msg, cl.Hub.Broadcast)
}

func (cl *Client) HandlePlayerEatFood(mData map[string]interface{}) {
	// parse inputs
	playerId := mData["playerId"].(string)
	foodId := mData["foodId"].(string)
	// datastore loads
	player := cl.GameState.Players.Get(playerId).(Player)
	food := cl.GameState.Foods.Get(foodId).(Food)
	// player grows in size
	player.Size += 1
	// food position changes
	food.Position = cl.GameState.GetNewSpawnFoodPosition()
	// datastore saves
	cl.GameState.Players.Set(playerId, player)
	cl.GameState.Foods.Set(foodId, food)
	// broadcast entity update messages
	pmsg := NewPlayerStateUpdateMessage(player)
	SerializeAndScheduleServerMessage(pmsg, cl.Hub.Broadcast)
	fmsg := NewFoodStateUpdateMessage(food)
	SerializeAndScheduleServerMessage(fmsg, cl.Hub.Broadcast)
}

func (cl *Client) HandlePlayerEatPlayer(mData map[string]interface{}) {
	// parse inputs
	playerId := mData["playerId"].(string)
	otherPlayerId := mData["otherPlayerId"].(string)
	// datastore loads
	player := cl.GameState.Players.Get(playerId).(Player)
	otherPlayer := cl.GameState.Players.Get(otherPlayerId).(Player)
	// player who ate other grows in size
	player.Size += otherPlayer.Size
	// player who got eaten respawns with reset size to 1
	otherPlayer.Size = 1
	otherPlayer.Position = cl.GameState.GetNewSpawnPlayerPosition()
	// datastore saves
	cl.GameState.Players.Set(playerId, player)
	cl.GameState.Players.Set(otherPlayerId, otherPlayer)
	// broadcast entity update messages
	pmsg := NewPlayerStateUpdateMessage(player)
	SerializeAndScheduleServerMessage(pmsg, cl.Hub.Broadcast)
	opmsg := NewPlayerStateUpdateMessage(otherPlayer)
	SerializeAndScheduleServerMessage(opmsg, cl.Hub.Broadcast)
}

func (cl *Client) HandleMineDamagePlayer(mData map[string]interface{}) {
	// parse inputs
	playerId := mData["playerId"].(string)
	mineId := mData["mineId"].(string)
	// datastore loads
	player := cl.GameState.Players.Get(playerId).(Player)
	mine := cl.GameState.Mines.Get(mineId).(Mine)
	// player loses size points
	player.Size -= 3
	// mine position changes
	mine.Position = cl.GameState.GetNewSpawnMinePosition()
	// datastore saves
	cl.GameState.Players.Set(playerId, player)
	cl.GameState.Mines.Set(mineId, mine)
	// broadcast entity update messages
	pmsg := NewPlayerStateUpdateMessage(player)
	SerializeAndScheduleServerMessage(pmsg, cl.Hub.Broadcast)
	mmsg := NewMineStateUpdateMessage(mine)
	SerializeAndScheduleServerMessage(mmsg, cl.Hub.Broadcast)
}

func (cl *Client) SendMessages() {
	for message := range cl.Send {
		SendJsonMessage(cl.Ws, message)
	}
}

func (cl *Client) Cleanup() {
	close(cl.Send)
}
