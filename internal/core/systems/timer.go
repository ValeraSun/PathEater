package systems

import (
	"time"

	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type TimerSystem struct {
	getter      componentsGetter
	broadcaster Broadcaster
	timer       timer
}

func NewTimerSystem(getter componentsGetter, broadcaster Broadcaster, timer timer) *TimerSystem {
	return &TimerSystem{
		getter:      getter,
		broadcaster: broadcaster,
		timer:       timer,
	}
}

type timeData struct {
	time time.Duration
}

func (s *TimerSystem) Update(dt float32) error {
	time := s.timer.GetRemainingTime()
	s.broadcaster.SendTime(ecs.TimeInfo{
		Type: "time",
		Data: timeData{
			Time: time,
		}
	})
	return nil
}