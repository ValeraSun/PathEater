package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type GameOverSystem struct {
	getter      componentsGetter
	broadcaster Broadcaster
	closer      closer
}

func NewGameOverSystem(getter componentsGetter, subscriber subscriber, broadcaster Broadcaster, closer closer) *GameOverSystem {
	s := &GameOverSystem{
		getter:      getter,
		broadcaster: broadcaster,
		closer:      closer,
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
	s.broadcaster.SendGameOverState(ecs.GameOverInfo{
		Type: "gameOver",
		Data: GameOverData{
			Win: ev.Win,
			//Status: s.status,
		},
	})
	//s.closer.Close()
	return nil
}
