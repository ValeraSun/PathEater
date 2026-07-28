package components

import (
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type WeaponComponent struct {
	StartAngle      float64
	EndAngle        float64
	Direction       float64
	Speed           float64
	TurningSpeed    float64
	Ammo            int
	currentCooldown float32
	Cooldown        float32
	Length          float64
}

func (*WeaponComponent) Type() string {
	return "weapon"
}

func (w *WeaponComponent) ApplyState(state *events.WeaponEvent, dt float32) {
	if w.currentCooldown-dt > 0 {
		w.currentCooldown -= dt
	} else {
		w.currentCooldown = 0
	}

	if state.Left {
		if w.Direction-w.TurningSpeed*float64(dt) > w.StartAngle {
			w.Direction -= w.TurningSpeed * float64(dt)
		} else {
			w.Direction = w.StartAngle
		}
	}

	if state.Right {
		if w.Direction+w.TurningSpeed*float64(dt) < w.EndAngle {
			w.Direction += w.TurningSpeed * float64(dt)
		} else {
			w.Direction = w.EndAngle
		}
	}

	if state.Shooting && w.currentCooldown == 0 {
		w.currentCooldown = w.Cooldown
		events.NewCreateBulletEvent(
			geometry.AngleToVec(w.Direction).Scale(w.Length),
			*geometry.AngleToVec(w.Direction),
		)
	}

}
