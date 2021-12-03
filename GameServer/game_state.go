package main

import (
	"fmt"
	"time"
)

type GameStateManager struct {
	GameState *GameState
}

func NewGameStateManager() *GameStateManager {
	gs := NewGameState()
	gsm := &GameStateManager{
		GameState: gs,
	}
	return gsm
}

func (gsm *GameStateManager) RunRoundTicker() {
	gs := gsm.GameState
	for range time.Tick(time.Second) {
		if gs.RoundIsInProgress {
			fmt.Println("Seconds left in round:", gs.SecondsToCurrentRoundEnd)
			if gs.SecondsToCurrentRoundEnd == 0 {
				gs.RoundIsInProgress = false
				gs.SecondsToNextRoundStart = SECONDS_BETWEEN_ROUNDS
				// TODO: end current round
			} else {
				// count down
				gs.SecondsToCurrentRoundEnd -= 1
			}
		} else {
			fmt.Println("Seconds until next round:", gs.SecondsToNextRoundStart)
			if gs.SecondsToNextRoundStart == 0 {
				gs.RoundIsInProgress = true
				gs.SecondsToCurrentRoundEnd = SECONDS_PER_ROUND
				// TODO: start next round
			} else {
				// count down
				gs.SecondsToNextRoundStart -= 1
			}
		}
	}
}

type GameState struct {
	Players                  map[string]*Player `json:"players"`
	Foods                    map[string]*Food   `json:"foods"`
	Mines                    map[string]*Mine   `json:"mines"`
	RoundHistory             map[string]*Round  `json:"roundHistory"`
	RoundCurrent             *Round             `json:"roundCurrent"`
	RoundIsInProgress        bool               `json:"roundIsInProgress"`
	SecondsToCurrentRoundEnd int                `json:"secondsToCurrentRoundEnd"`
	SecondsToNextRoundStart  int                `json:"secondsToNextRoundStart"`
}

func NewGameState() *GameState {
	gs := &GameState{
		Players:                  make(map[string]*Player),
		Foods:                    make(map[string]*Food),
		Mines:                    make(map[string]*Mine),
		RoundHistory:             make(map[string]*Round),
		RoundCurrent:             nil,
		RoundIsInProgress:        false,
		SecondsToCurrentRoundEnd: 0,
		SecondsToNextRoundStart:  0,
	}
	// TODO: do initialization of game state here
	return gs
}

func (gs *GameState) StartNewRound() {
	// stub
	// logic:
	// - finalize RoundCurrent and store to RoundHistory
	// - create new round and assign to RoundCurrent
}
