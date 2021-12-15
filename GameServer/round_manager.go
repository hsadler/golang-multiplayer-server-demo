package main

import (
	"time"
)

type RoundManager struct {
	Hub                      *Hub
	GameState                *GameState
	RoundIsInProgress        bool
	SecondsToNextRoundStart  int
	SecondsToCurrentRoundEnd int
}

func (rm *RoundManager) RunRoundTicker() {
	for range time.Tick(time.Second) {
		if rm.RoundIsInProgress {
			// LogForce("Seconds left in round:", rm.SecondsToCurrentRoundEnd)
			if rm.SecondsToCurrentRoundEnd == 0 {
				// end the current round
				rm.RoundIsInProgress = false
				rm.SecondsToNextRoundStart = SECONDS_BETWEEN_ROUNDS
				message := NewRoundResultMessage(rm.GameState.GetRoundResult())
				SerializeAndScheduleServerMessage(message, rm.Hub.Broadcast)
			} else {
				// count down to round end
				rm.SecondsToCurrentRoundEnd -= 1
			}
			message := NewSecondsToCurrentRoundEndMessage(rm.SecondsToCurrentRoundEnd)
			SerializeAndScheduleServerMessage(message, rm.Hub.Broadcast)
		} else {
			// LogForce("Seconds until next round:", rm.SecondsToNextRoundStart)
			if rm.SecondsToNextRoundStart == 0 {
				// initialize game state for the new round and broadcast
				rm.RoundIsInProgress = true
				rm.SecondsToCurrentRoundEnd = SECONDS_PER_ROUND
				rm.GameState.InitNewRoundGameState()
				message := NewGameStateMessage(rm.GameState.GetSerializable())
				SerializeAndScheduleServerMessage(message, rm.Hub.Broadcast)
			} else {
				// count down to next round
				rm.SecondsToNextRoundStart -= 1
			}
			message := NewSecondsToNextRoundStartMessage(rm.SecondsToNextRoundStart)
			SerializeAndScheduleServerMessage(message, rm.Hub.Broadcast)
		}
	}
}
