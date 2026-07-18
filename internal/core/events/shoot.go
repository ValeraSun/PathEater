package events

import (
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type ShootEvent struct {
	Direction geometry.Vec3
	Speed     float64
}

func (*ShootEvent) Type() string { return "shoot" }

func NewShootEvent(dir geometry.Vec3, speed float64) *ShootEvent {
	return &ShootEvent{
		Direction: dir,
		Speed:     speed,
	}
}
