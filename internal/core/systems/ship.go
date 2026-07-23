package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type ShipSystem struct {
	getter     componentsGetter
	eventQueue chan *events.SetPlayerStateEvent
}

func NewShipSystem(getter componentsGetter, subscriber subscriber) *ShipSystem {
	s := &ShipSystem{
		getter: getter,
	}
	subscriber.Subscribe("setShipState", s.OnEvent)
	return s
}

func (s *ShipSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("ship")

	s.moveShip(comps)

	return nil
}

func (s *ShipSystem) moveShip(comps map[types.Entity]types.Component) {
	for id := range comps {
		select {
		case e := <-s.eventQueue:
			if e.ID == string(id) {
				c, _ := s.getter.GetComponent(id, "transform")
				transform := c.(*components.MovementComponent)

				transform.Direction = GetMoveVector(e)
			}
		default:
			return
		}
	}
}

func GetMoveVector(e *events.SetPlayerStateEvent) geometry.Vec3 {
	vec := geometry.GetZeroVector()
	if e.PlayerState.MoveFront {
		vec.Add(geometry.Vec3{X: 1, Y: 1, Z: 0})
	}
	if e.PlayerState.MoveBack {
		vec.Add(geometry.Vec3{X: 1, Y: -1, Z: 0})
	}
	return vec.Normalize()
}

func (s *ShipSystem) OnEvent(event events.Event) error {
	ps, ok := event.(*events.SetPlayerStateEvent)

	if !ok {
		return nil
	}

	select {
	case s.eventQueue <- ps:
	default:
		fmt.Printf("Переполена очередь %v\n", s)
	}

	return nil
}
