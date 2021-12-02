package main

// TODO: do we need this, or put it on the hub?
type GameStateManager struct {
	GameState *GameState
}

type GameState struct {
	Players                  map[string]*Player `json:"players"`
	Foods                    map[string]*Food   `json:"foods"`
	Mines                    map[string]*Mine   `json:"mines"`
	RoundHistory             map[string]*Round  `json:"roundHistory"`
	RoundCurrent             *Round             `json:"roundCurrent"`
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
		SecondsToCurrentRoundEnd: 0,
		SecondsToNextRoundStart:  0,
	}
	// TODO: do initialization of game state here
	return gs
}

func (gs *GameState) Run() {
	// stub
}
