package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type GameOverSystem struct {
	getter componentsGetter
}

func NewGameOverSystem(getter componentsGetter, subscriber subscriber, world ecs.) *GameOverSystem {
	s := &GameOverSystem{
		getter: getter,
	}
	subscriber.Subscribe("gameOver", s.OnEvent)
	return s
}

func (s *GameOverSystem) Update(dt float32) error {
	return nil
}

func (s *GameOverSystem) OnEvent(event events.Event) error {

	return nil
}
