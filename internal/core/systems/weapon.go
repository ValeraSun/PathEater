package systems

import (
	"errors"
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
	subscriber.Subscribe("setWeaponState", s.OnEvent)
	return s
}

func (s *WeaponSystem) Update(dt float32) error {

	for {
		select {
		case e := <-s.eventQueue:
			if !s.HasComponents(e.ID, "weapon") {
				return errors.New("weapon system error")
			}

			c, _ := s.GetComponent(e.ID, "weapon")
			weapon := c.(*components.WeaponComponent)

			weapon.ApplyState(e, dt)

		default:
			return nil
		}
	}

}

func (s *WeaponSystem) handleEvents() {

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
