package main

import (
	"encoding/json"
	"fmt"
)

type GameState struct {
	MapHeight         int                `json:"mapHeight"`
	MapWidth          int                `json:"mapWidth"`
	Players           map[string]*Player `json:"players"`
	Foods             map[string]*Food   `json:"foods"`
	Mines             map[string]*Mine   `json:"mines"`
	AddPlayer         chan *Player       `json:"-"`
	RemovePlayer      chan *Player       `json:"-"`
	UpdatePlayerState chan *Player       `json:"-"`
	UpdateFoodState   chan *Food         `json:"-"`
	UpdateMineState   chan *Mine         `json:"-"`
}

func (gs *GameState) RunListeners() {
	for {
		select {
		case player := <-gs.AddPlayer:
			fmt.Println("Adding player to game state:", player)
			gs.Players[player.Id] = player
		case player := <-gs.RemovePlayer:
			fmt.Println("Removing player from game state:", player)
			delete(gs.Players, player.Id)
		case player := <-gs.UpdatePlayerState:
			fmt.Println("Updating player state in game state:", player)
			gs.Players[player.Id] = player
		case food := <-gs.UpdateFoodState:
			fmt.Println("Updating food state in game state:", food)
			gs.Foods[food.Id] = food
		case mine := <-gs.UpdateMineState:
			fmt.Println("Updating mine state in game state:", mine)
			gs.Mines[mine.Id] = mine
		}
	}
}

func (gs *GameState) GetRoundResult() *Round {
	// aggregate game state to round object
	playerIdToScore := make(map[string]int)
	for _, player := range gs.Players {
		playerIdToScore[player.Id] = player.Size
	}
	r := &Round{
		PlayerIdToScore: playerIdToScore,
	}
	res, _ := json.Marshal(r)
	LogJson("round result:", res)
	return r
}

func (gs *GameState) InitNewRoundGameState() {
	// randomize food placement
	foods := make(map[string]*Food)
	for i := 0; i < FOOD_COUNT; i++ {
		f := &Food{
			Id:     GenUUID(),
			Active: true,
			// TODO: add randomization to food positions
			Position: &Position{
				X: 0,
				Y: 0,
			},
		}
		foods[f.Id] = f
	}
	gs.Foods = foods
	// randomize mine placement
	mines := make(map[string]*Mine)
	for i := 0; i < MINE_COUNT; i++ {
		m := &Mine{
			Id:     GenUUID(),
			Active: true,
			// TODO: add randomization to mine positions
			Position: &Position{
				X: 0,
				Y: 0,
			},
		}
		mines[m.Id] = m
	}
	gs.Mines = mines
	// randomize player placement and make sure they're not too close to another
	// player or mine
	for _, p := range gs.Players {
		p.Position = gs.GetNewSpawnPlayerPosition()
	}
	res, _ := json.Marshal(gs)
	LogJson("initialized game-state:", res)
}

func (gs *GameState) GetNewSpawnPlayerPosition() *Position {
	// TODO: add randomization to player positions and ensure they have space
	return &Position{
		X: 0,
		Y: 0,
	}
}
