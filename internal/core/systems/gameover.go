package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/sendler"
)

type GameOverSystem struct {
	componentsGetter
	broadcaster
	closer
}

func NewGameOverSystem(getter componentsGetter, subscriber subscriber, broadcaster broadcaster, closer closer) *GameOverSystem {
	s := &GameOverSystem{
		getter,
		broadcaster,
		closer,
	}
	subscriber.Subscribe("gameOver", s.OnEvent)
	return s
}

func (s *GameOverSystem) Update(dt float32) error {
	return nil
}

type GameOverData struct {
	Win    bool `json:"win"`
	Status int  `json:"status"`
}

func (s *GameOverSystem) OnEvent(event events.Event) error {
	ev := event.(*events.GameOverEvent)
	s.broadcaster.SendGameOverState(sendler.GameOverInfo{
		Type: "gameOver",
		Data: GameOverData{
			Win: ev.Win,
			//Status: s.status,
		},
	})
	//s.closer.Close()
	return nil
}
