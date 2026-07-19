package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type WeaponSystem struct {
	getter     componentsGetter
	eventQueue chan *events.SetWeaponStateEvent
}

func NewWeaponSystem(getter componentsGetter, subscriber subscriber) *WeaponSystem {
	s := &WeaponSystem{
		getter:     getter,
		eventQueue: make(chan *events.SetWeaponStateEvent, 100),
	}
	subscriber.Subscribe("setWeaponState", s.OnEvent)
	return s
}

func (s *WeaponSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("weapon")

	s.handleWeapon(comps)

	return nil
}

func (s *WeaponSystem) handleWeapon(comps map[types.Entity]types.Component) {
	for {
		select {
		case e := <-s.eventQueue:
			for _, comp := range comps {
				weap := comp.(*components.WeaponComponent)
				var angle float64 = 0
				if e.WeaponState.TurnClockwise {
					angle = angle + weap.Speed
				}
				if e.WeaponState.TurnCounterclockwise {
					angle = angle - weap.Speed
				}
				weap.Direction.Rotate(angle)
				if e.WeaponState.Shoot {
					if weap.Ammo > 0 {
						weap.ShootSuccess = true
						events.NewShootEvent(weap.Direction, 10)
						weap.Ammo--
					} else {
						weap.ShootSuccess = false
					}
				}
			}
		default:
			return
		}
	}
}

func (s *WeaponSystem) OnEvent(event events.Event) error {
	ws, ok := event.(*events.SetWeaponStateEvent)

	if !ok {
		return nil
	}

	select {
	case s.eventQueue <- ws:

	default:
		fmt.Printf("Переполена очередь %v\n", s)
	}

	return nil
}
