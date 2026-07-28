package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/transfer"
)

type WeaponSystem struct {
	componentsGetter
	transfer.EventPublisher
	eventQueue chan *events.WeaponEvent
}

func NewWeaponSystem(getter componentsGetter, publisher transfer.EventPublisher, subscriber subscriber) *WeaponSystem {
	s := &WeaponSystem{
		getter,
		publisher,
		make(chan *events.WeaponEvent, 100),
	}
	subscriber.Subscribe("weapon", s.OnEvent)
	return s
}

func (s *WeaponSystem) Update(dt float32) error {
	comps := s.GetEntitiesByComponent("weapon")
	for _, c := range comps {
		weapon := c.(*components.WeaponComponent)
		s.handleEvents(weapon, dt)
	}
	return nil
}

func (s *WeaponSystem) handleEvents(weapon *components.WeaponComponent, dt float32) {
	for {
		select {
		case e := <-s.eventQueue:

			event := weapon.ApplyState(e, dt)
			if event != nil {
				s.Publish(event)
			}
		default:
			return
		}
	}
}

func BulletPos(dir geometry.Vec3) geometry.Vec3 {
	return dir.Scale(1.5)
}

func (s *WeaponSystem) OnEvent(event events.Event) error {
	ps, _ := event.(*events.WeaponEvent)

	select {
	case s.eventQueue <- ps:

	default:
		fmt.Printf(" 3 Переполена очередь %+v\n", *s)
	}

	return nil
}
