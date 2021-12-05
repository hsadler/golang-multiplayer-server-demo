package main

import (
	"encoding/json"
	"sync"
)

type GameState struct {
	MapHeight         int
	MapWidth          int
	Players           map[string]Player
	Foods             map[string]Food
	Mines             map[string]Mine
	AddPlayer         chan Player
	RemovePlayer      chan Player
	UpdatePlayerState chan Player
	UpdateFoodState   chan Food
	UpdateMineState   chan Mine
	Mu                *sync.RWMutex
}

func (gs *GameState) RunWriteListeners() {
	for {
		select {
		case player := <-gs.AddPlayer:
			LogData("Adding player to game state:", player.Id)
			gs.Mu.Lock()
			gs.Players[player.Id] = player
			gs.Mu.Unlock()
		case player := <-gs.RemovePlayer:
			LogData("Removing player from game state:", player.Id)
			gs.Mu.Lock()
			delete(gs.Players, player.Id)
			gs.Mu.Unlock()
		case player := <-gs.UpdatePlayerState:
			LogData("Updating player state in game state:", player.Id)
			gs.Mu.Lock()
			gs.Players[player.Id] = player
			gs.Mu.Unlock()
		case food := <-gs.UpdateFoodState:
			LogData("Updating food state in game state:", food.Id)
			gs.Mu.Lock()
			gs.Foods[food.Id] = food
			gs.Mu.Unlock()
		case mine := <-gs.UpdateMineState:
			LogData("Updating mine state in game state:", mine.Id)
			gs.Mu.Lock()
			gs.Mines[mine.Id] = mine
			gs.Mu.Unlock()
		}
	}
}

func (gs *GameState) FetchPlayer(playerId string) Player {
	gs.Mu.RLock()
	defer gs.Mu.RUnlock()
	return gs.Players[playerId]
}

func (gs *GameState) GetRoundResult() Round {
	gs.Mu.RLock()
	defer gs.Mu.RUnlock()
	// aggregate game state to round object
	playerIdToScore := make(map[string]int)
	for _, player := range gs.Players {
		playerIdToScore[player.Id] = player.Size
	}
	r := Round{
		PlayerIdToScore: playerIdToScore,
	}
	res, _ := json.Marshal(r)
	LogJson("round result:", res)
	return r
}

func (gs *GameState) InitNewRoundGameState() {
	gs.Mu.Lock()
	defer gs.Mu.Unlock()
	// food placement
	foods := make(map[string]Food)
	for i := 0; i < FOOD_COUNT; i++ {
		f := Food{
			Id:       GenUUID(),
			Active:   true,
			Position: gs.GetNewSpawnFoodPosition(),
		}
		foods[f.Id] = f
	}
	gs.Foods = foods
	// mine placement
	mines := make(map[string]Mine)
	for i := 0; i < MINE_COUNT; i++ {
		m := Mine{
			Id:       GenUUID(),
			Active:   true,
			Position: gs.GetNewSpawnMinePosition(),
		}
		mines[m.Id] = m
	}
	gs.Mines = mines
	// player placement
	for _, p := range gs.Players {
		p.Position = gs.GetNewSpawnPlayerPosition()
	}
}

func (gs *GameState) GetNewSpawnMinePosition() *Position {
	// TODO: add randomization to mine position and ensure it has space
	// away from mines
	return &Position{
		X: 0,
		Y: 0,
	}
}

func (gs *GameState) GetNewSpawnFoodPosition() *Position {
	// TODO: add randomization to food position and ensure it has space
	// away from mines
	return &Position{
		X: 0,
		Y: 0,
	}
}

func (gs *GameState) GetNewSpawnPlayerPosition() *Position {
	// TODO: add randomization to player position and ensure it has space
	// away from mines and other players
	return &Position{
		X: 0,
		Y: 0,
	}
}

type GameStateSerializable struct {
	MapHeight int      `json:"mapHeight"`
	MapWidth  int      `json:"mapWidth"`
	Players   []Player `json:"players"`
	Foods     []Food   `json:"foods"`
	Mines     []Mine   `json:"mines"`
}

func (gs *GameState) GetSerializable() GameStateSerializable {
	gs.Mu.RLock()
	defer gs.Mu.RUnlock()
	players := make([]Player, 0)
	for _, p := range gs.Players {
		players = append(players, p)
	}
	foods := make([]Food, 0)
	for _, f := range gs.Foods {
		foods = append(foods, f)
	}
	mines := make([]Mine, 0)
	for _, m := range gs.Mines {
		mines = append(mines, m)
	}
	return GameStateSerializable{
		MapHeight: gs.MapHeight,
		MapWidth:  gs.MapWidth,
		Players:   players,
		Foods:     foods,
		Mines:     mines,
	}
}
