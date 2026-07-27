package systems

import (
	"time"

	"github.com/ValeraSun/PathEater/internal/core/ecs"
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
	Time time.Duration
}

func (s *TimerSystem) Update(dt float32) error {
	remaining := s.timer.GetRemainingTime()
	seconds := int(remaining.Seconds())

	if err := s.broadcaster.SendTime(ecs.TimeInfo{
		Type: "time",
		Data: seconds,
	}); err != nil {
		return err
	}

	return nil
}
