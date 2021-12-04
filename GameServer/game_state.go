package main

import (
	"fmt"
	"time"
)

type GameStateManager struct {
	Hub                      *Hub
	GameState                *GameState
	RoundIsInProgress        bool
	SecondsToNextRoundStart  int
	SecondsToCurrentRoundEnd int
}

func (gsm *GameStateManager) RunRoundTicker() {
	for range time.Tick(time.Second) {
		if gsm.RoundIsInProgress {
			fmt.Println("Seconds left in round:", gsm.SecondsToCurrentRoundEnd)
			if gsm.SecondsToCurrentRoundEnd == 0 {
				// end the current round
				gsm.RoundIsInProgress = false
				gsm.SecondsToNextRoundStart = SECONDS_BETWEEN_ROUNDS
				r := gsm.GameState.EndCurrentRound()
				message := NewRoundResultMessage(r)
				SerializeAndScheduleServerMessage(message, gsm.Hub.Broadcast)
			} else {
				// count down to round end
				gsm.SecondsToCurrentRoundEnd -= 1
			}
			message := NewSecondsToCurrentRoundEndMessage(gsm.SecondsToCurrentRoundEnd)
			SerializeAndScheduleServerMessage(message, gsm.Hub.Broadcast)
		} else {
			fmt.Println("Seconds until next round:", gsm.SecondsToNextRoundStart)
			if gsm.SecondsToNextRoundStart == 0 {
				// initialize game state for the new round and broadcast
				gsm.RoundIsInProgress = true
				gsm.SecondsToCurrentRoundEnd = SECONDS_PER_ROUND
				gsm.GameState.InitNewRoundGameState()
				message := NewGameStateMessage(gsm.GameState)
				SerializeAndScheduleServerMessage(message, gsm.Hub.Broadcast)
			} else {
				// count down to next round
				gsm.SecondsToNextRoundStart -= 1
			}
			message := NewSecondsToNextRoundStartMessage(gsm.SecondsToNextRoundStart)
			SerializeAndScheduleServerMessage(message, gsm.Hub.Broadcast)
		}
	}
}

type GameState struct {
	MapHeight int                `json:"mapHeight"`
	MapWidth  int                `json:"mapWidth"`
	Players   map[string]*Player `json:"players"`
	Foods     map[string]*Food   `json:"foods"`
	Mines     map[string]*Mine   `json:"mines"`
}

func (gs *GameState) EndCurrentRound() *Round {
	// aggregate game state to round object
	playerIdToScore := make(map[string]int)
	for _, player := range gs.Players {
		playerIdToScore[player.Id] = player.Size
	}
	r := &Round{
		PlayerIdToScore: playerIdToScore,
	}
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
		// TODO: add randomization to player positions and ensure they have space
		newPosition := &Position{
			X: 0,
			Y: 0,
		}
		p.Position = newPosition
	}
}
