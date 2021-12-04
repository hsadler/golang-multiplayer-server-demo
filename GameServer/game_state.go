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

func (gs *GameState) EndCurrentRound() {
	// TODO: stub
	// - finalize RoundCurrent and store to RoundHistory
}

func (gs *GameState) StartNewRound() {
	// TODO: stub
	// - create new round and assign to RoundCurrent
}
