package events

type GameOverEvent struct {
	Result int
}

func (*GameOverEvent) Type() string { return "gameOver" }

func NewGameOverEvent(res int) *GameOverEvent {
	return &GameOverEvent{
		Result: res,
	}
}
