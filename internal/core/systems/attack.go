package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type AttackSystem struct {
	getter     componentsGetter
	eventQueue chan *events.AttackEvent
}

func NewAttackSystem(getter componentsGetter) *AttackSystem {
	return &AttackSystem{
		getter: getter,
	}
}

func (s *AttackSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("control")

	return nil
}

func (s *AttackSystem) drainEvents(comps map[types.Entity]types.Component) {
	for {
		select {
		case e := <-s.eventQueue:
			comp, ok := comps[e.ID]
			if !ok {
				continue
			}
			c, ok := comp.(*components.ControlComponent)
			if !ok {
				continue
			}
			c.Superimpose(&e.PlayerState)
		default:
			return
		}
	}
}
