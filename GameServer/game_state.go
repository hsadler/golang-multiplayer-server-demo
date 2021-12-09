package main

import (
	"encoding/json"
)

type GameState struct {
	MapHeight int
	MapWidth  int
	Players   *CMap
	Foods     *CMap
	Mines     *CMap
}

func (gs *GameState) GetRoundResult() Round {
	// aggregate game state to round object
	playerIdToScore := make(map[string]int)
	for _, playerData := range gs.Players.Values() {
		player := playerData.(Player)
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
	// food placement
	gs.Foods = NewCMap()
	for i := 0; i < FOOD_COUNT; i++ {
		f := Food{
			Id:       GenUUID(),
			Active:   true,
			Position: gs.GetNewSpawnFoodPosition(),
			Size:     FOOD_SIZE,
		}
		gs.Foods.Set(f.Id, f)
	}
	// mine placement
	gs.Mines = NewCMap()
	for i := 0; i < MINE_COUNT; i++ {
		m := Mine{
			Id:       GenUUID(),
			Active:   true,
			Position: gs.GetNewSpawnMinePosition(),
			Size:     MINE_SIZE,
		}
		gs.Mines.Set(m.Id, m)
	}
	// player placement
	for _, pData := range gs.Players.Values() {
		player := pData.(Player)
		player.Position = gs.GetNewSpawnPlayerPosition()
		gs.Players.Set(player.Id, player)
	}
}

func (gs *GameState) GetNewSpawnMinePosition() Position {
	// TODO: add randomization to mine position and ensure it has space
	// away from mines
	return GenRandPosition(gs)
}

func (gs *GameState) GetNewSpawnFoodPosition() Position {
	// TODO: add randomization to food position and ensure it has space
	// away from mines
	return GenRandPosition(gs)
}

func (gs *GameState) GetNewSpawnPlayerPosition() Position {
	// TODO: add randomization to player position and ensure it has space
	// away from mines and other players
	return GenRandPosition(gs)
}

type GameStateSerializable struct {
	MapHeight int      `json:"mapHeight"`
	MapWidth  int      `json:"mapWidth"`
	Players   []Player `json:"players"`
	Foods     []Food   `json:"foods"`
	Mines     []Mine   `json:"mines"`
}

func (gs *GameState) GetSerializable() GameStateSerializable {
	players := make([]Player, 0)
	for _, p := range gs.Players.Values() {
		players = append(players, p.(Player))
	}
	foods := make([]Food, 0)
	for _, f := range gs.Foods.Values() {
		foods = append(foods, f.(Food))
	}
	mines := make([]Mine, 0)
	for _, m := range gs.Mines.Values() {
		mines = append(mines, m.(Mine))
	}
	return GameStateSerializable{
		MapHeight: gs.MapHeight,
		MapWidth:  gs.MapWidth,
		Players:   players,
		Foods:     foods,
		Mines:     mines,
	}
}
