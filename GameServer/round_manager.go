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
				SerializeAndScheduleServerMessage(
					NewRoundResultMessage(rm.GameState.GetRoundResult()),
					rm.Hub.Broadcast,
				)
			} else {
				// count down to round end
				rm.SecondsToCurrentRoundEnd -= 1
				// count down to respawn for players who are waiting
				for _, pData := range rm.GameState.Players.Values() {
					p := pData.(Player)
					if !p.Active && p.TimeUntilRespawn > 0 {
						p.TimeUntilRespawn -= 1
						if p.TimeUntilRespawn == 0 {
							p.Active = true
							p.Position = rm.GameState.GetNewSpawnPlayerPosition()
						}
						rm.GameState.Players.Set(p.Id, p)
						SerializeAndScheduleServerMessage(
							NewPlayerStateUpdateMessage(p),
							rm.Hub.Broadcast,
						)
					}
				}
			}
			SerializeAndScheduleServerMessage(
				NewSecondsToCurrentRoundEndMessage(rm.SecondsToCurrentRoundEnd),
				rm.Hub.Broadcast,
			)
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
			SerializeAndScheduleServerMessage(
				NewSecondsToNextRoundStartMessage(rm.SecondsToNextRoundStart),
				rm.Hub.Broadcast,
			)
		}
	}
}
