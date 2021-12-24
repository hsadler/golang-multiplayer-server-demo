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
				rm.RoundEndProcs()
				rm.CountNextRoundStart()
			} else {
				rm.CountRoundInProgress()
				rm.CountPlayerRespawns()
			}
		} else {
			// LogForce("Seconds until next round:", rm.SecondsToNextRoundStart)
			if rm.SecondsToNextRoundStart == 0 {
				rm.RoundStartProcs()
				rm.CountRoundInProgress()
			} else {
				rm.CountNextRoundStart()
			}
		}
	}
}

func (rm *RoundManager) CountRoundInProgress() {
	// count down to round end
	SerializeAndScheduleServerMessage(
		NewSecondsToCurrentRoundEndMessage(rm.SecondsToCurrentRoundEnd),
		rm.Hub.Broadcast,
	)
	rm.SecondsToCurrentRoundEnd -= 1
}

func (rm *RoundManager) CountPlayerRespawns() {
	// count down to respawn for players who are waiting
	for _, pData := range rm.GameState.Players.Values() {
		p := pData.(Player)
		if !p.Active && p.TimeUntilRespawn > 0 {
			p.TimeUntilRespawn -= 1
			if p.TimeUntilRespawn == 0 {
				p.Active = true
			}
			rm.GameState.Players.Set(p.Id, p)
			SerializeAndScheduleServerMessage(
				NewPlayerStateUpdateMessage(p),
				rm.Hub.Broadcast,
			)
		}
	}
}

func (rm *RoundManager) CountNextRoundStart() {
	// count down to next round
	SerializeAndScheduleServerMessage(
		NewSecondsToNextRoundStartMessage(rm.SecondsToNextRoundStart),
		rm.Hub.Broadcast,
	)
	rm.SecondsToNextRoundStart -= 1
}

func (rm *RoundManager) RoundStartProcs() {
	// initialize game state for the new round and broadcast
	rm.RoundIsInProgress = true
	rm.SecondsToCurrentRoundEnd = SECONDS_PER_ROUND
	rm.GameState.InitNewRoundGameState()
	SerializeAndScheduleServerMessage(
		NewGameStateMessage(rm.GameState.GetSerializable()),
		rm.Hub.Broadcast,
	)
	// activate all players
	for _, pData := range rm.GameState.Players.Values() {
		player := pData.(Player)
		player.Active = true
		rm.GameState.Players.Set(player.Id, player)
		SerializeAndScheduleServerMessage(
			NewPlayerStateUpdateMessage(player),
			rm.Hub.Broadcast,
		)
	}
}

func (rm *RoundManager) RoundEndProcs() {
	// end the current round
	rm.RoundIsInProgress = false
	rm.SecondsToNextRoundStart = SECONDS_BETWEEN_ROUNDS
	SerializeAndScheduleServerMessage(
		NewRoundResultMessage(rm.GameState.GetRoundResult()),
		rm.Hub.Broadcast,
	)
	// deactivate and initialize all players, reset respawn times
	for _, pData := range rm.GameState.Players.Values() {
		player := pData.(Player)
		player.Active = false
		player.Size = 1
		player.TimeUntilRespawn = 0
		rm.GameState.Players.Set(player.Id, player)
		SerializeAndScheduleServerMessage(
			NewPlayerStateUpdateMessage(player),
			rm.Hub.Broadcast,
		)
	}
}
