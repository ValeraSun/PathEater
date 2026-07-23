package events

type GameOverEvent struct {
	Win bool
}

func (*GameOverEvent) Type() string { return "gameOver" }

func NewGameOverEvent(win bool) *GameOverEvent {
	return &GameOverEvent{
		Win: win,
	}
}
