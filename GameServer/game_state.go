package main

import (
	"fmt"
	"time"
)

type GameStateManager struct {
	GameState                *GameState
	RoundIsInProgress        bool
	SecondsToCurrentRoundEnd int
	SecondsToNextRoundStart  int
}

func NewGameStateManager() *GameStateManager {
	gs := NewGameState()
	gsm := &GameStateManager{
		GameState:                gs,
		RoundIsInProgress:        false,
		SecondsToCurrentRoundEnd: 0,
		SecondsToNextRoundStart:  SECONDS_BETWEEN_ROUNDS,
	}
	return gsm
}

func (gsm *GameStateManager) RunRoundTicker() {
	for range time.Tick(time.Second) {
		if gsm.RoundIsInProgress {
			fmt.Println("Seconds left in round:", gsm.SecondsToCurrentRoundEnd)
			if gsm.SecondsToCurrentRoundEnd == 0 {
				gsm.RoundIsInProgress = false
				gsm.SecondsToNextRoundStart = SECONDS_BETWEEN_ROUNDS
				gsm.GameState.EndCurrentRound()
			} else {
				// count down
				gsm.SecondsToCurrentRoundEnd -= 1
			}
		} else {
			fmt.Println("Seconds until next round:", gsm.SecondsToNextRoundStart)
			if gsm.SecondsToNextRoundStart == 0 {
				gsm.RoundIsInProgress = true
				gsm.SecondsToCurrentRoundEnd = SECONDS_PER_ROUND
				gsm.GameState.StartNewRound()
			} else {
				// count down
				gsm.SecondsToNextRoundStart -= 1
			}
		}
	}
}

type GameState struct {
	Players      map[string]*Player `json:"players"`
	Foods        map[string]*Food   `json:"foods"`
	Mines        map[string]*Mine   `json:"mines"`
	RoundHistory map[string]*Round  `json:"roundHistory"`
	RoundCurrent *Round             `json:"roundCurrent"`
}

func NewGameState() *GameState {
	gs := &GameState{
		Players:      make(map[string]*Player),
		Foods:        make(map[string]*Food),
		Mines:        make(map[string]*Mine),
		RoundHistory: make(map[string]*Round),
		RoundCurrent: nil,
	}
	// TODO: do initialization of game state here
	return gs
}

func (gs *GameState) EndCurrentRound() {
	// stub
	// - finalize RoundCurrent and store to RoundHistory
}

func (gs *GameState) StartNewRound() {
	// stub
	// - create new round and assign to RoundCurrent
}
