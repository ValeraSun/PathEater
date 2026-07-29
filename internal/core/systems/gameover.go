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

func (s *GameOverSystem) OnEvent(event events.Event) error {
	ev := event.(*events.GameOverEvent)
	_, ship := getShip(s)
	s.broadcaster.SendGameOverState(sendler.GameOverInfo{
		Win:    ev.Win,
		Status: ship.BaggageStatus,
	},
	)
	s.closer.Close()
	return nil
}
