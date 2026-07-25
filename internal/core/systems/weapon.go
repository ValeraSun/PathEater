package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/transfer"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type WeaponSystem struct {
	getter     componentsGetter
	publisher  transfer.EventPublisher
	eventQueue chan *events.SetWeaponStateEvent
}

func NewWeaponSystem(getter componentsGetter, publisher transfer.EventPublisher, subscriber subscriber) *WeaponSystem {
	s := &WeaponSystem{
		getter:     getter,
		publisher:  publisher,
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
						event := events.NewCreateBulletEvent(BulletPos(weap.Direction), weap.Direction)
						s.publisher.Publish(event)
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

func BulletPos(dir geometry.Vec3) geometry.Vec3 {
	return dir.Scale(1.5)
}

func (s *WeaponSystem) OnEvent(event events.Event) error {
	ws, ok := event.(*events.SetWeaponStateEvent)

	if !ok {
		return nil
	}

	select {
	case s.eventQueue <- ws:

	default:
		fmt.Printf("Переполена очередь %+v\n", *s)
	}

	return nil
}
