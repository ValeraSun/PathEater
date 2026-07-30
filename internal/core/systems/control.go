package systems

import (
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type ControlSystem struct {
	getter     componentsGetter
	eventQueue chan *events.SetPlayerStateEvent
}

func NewControlSystem(getter componentsGetter, subscriber subscriber) *ControlSystem {
	s := &ControlSystem{
		getter:     getter,
		eventQueue: make(chan *events.SetPlayerStateEvent, eventQueueSize*6),
	}
	subscriber.Subscribe("setPlayerState", s.OnEvent)
	return s
}

func (s *ControlSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("control")
	s.resetAll(comps)

	s.drainEvents(comps)

	return nil
}

func (s *ControlSystem) resetAll(comps map[types.Entity]types.Component) {
	for _, c := range comps {
		control := c.(*components.ControlComponent)

		control.Reset()
	}
}

func (s *ControlSystem) drainEvents(comps map[types.Entity]types.Component) {
	for {
		select {
		case e := <-s.eventQueue:
			comp, ok := comps[types.Entity(e.ID)]
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

func (s *ControlSystem) OnEvent(event events.Event) error {
	ps, ok := event.(*events.SetPlayerStateEvent)

	if !ok {
		return nil
	}

	select {
	case s.eventQueue <- ps:

	default:
		fmt.Printf("1 Переполена очередь %+v\n", *s)
	}

	return nil
}
