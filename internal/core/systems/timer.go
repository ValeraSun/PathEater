package systems

import (
	"time"

	"github.com/ValeraSun/PathEater/internal/core/sendler"
)

type TimerSystem struct {
	componentsGetter
	broadcaster
	timer
}

func NewTimerSystem(getter componentsGetter, broadcaster broadcaster, timer timer) *TimerSystem {
	return &TimerSystem{
		getter,
		broadcaster,
		timer,
	}
}

type timeData struct {
	Time time.Duration
}

func (s *TimerSystem) Update(dt float32) error {
	remaining := s.timer.GetRemainingTime()
	seconds := int(remaining.Seconds())

	if err := s.SendTime(sendler.TimeInfo{
		Type: "time",
		Data: seconds,
	}); err != nil {
		return err
	}

	return nil
}
