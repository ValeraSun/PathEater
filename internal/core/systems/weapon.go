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

func NewWeaponSystem(getter componentsGetter) *WeaponSystem {
	s := &WeaponSystem{
		getter:     getter,
		eventQueue: make(chan *events.SetWeaponStateEvent, 100),
	}
	subscriber.Subscribe("setWeaponState", s.OnEvent)
	return s
}

func (s *WeaponSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("weapon")

	s.rotateWeapon(comps)

	return nil
}

func (s *WeaponSystem) rotateWeapon(comps map[types.Entity]types.Component) {
	for {
		select {
		case e := <-s.eventQueue:
			for _, comp := range comps {
				weap := comp.(*components.WeaponComponent)
				var angle float64
				if e.WeaponState.TurnClockwise {
					angle = weap.Speed
				}
				if e.WeaponState.TurnCounterclockwise {
					angle = -weap.Speed
				}
				weap.Direction.Rotate(angle)
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
