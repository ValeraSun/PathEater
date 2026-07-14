package systems

import (
	"fmt"
	"log"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type ControlSystem struct {
	getter     componentsGetter
	eventQueue chan *events.SetPlayerStateEvent
}

func NewPlayerControlSystem(getter componentsGetter) *ControlSystem {
	return &ControlSystem{
		getter:     getter,
		eventQueue: make(chan *events.SetPlayerStateEvent, 100),
	}
}

func (s *ControlSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("control")

	s.resetComponents(comps)

	s.drainEvents(comps)

	return nil
}

func (s *ControlSystem) resetComponents(comps map[types.Entity]types.Component) {
	for _, comp := range comps {
		c, ok := comp.(*components.ControlComponent)
		if !ok {
			log.Printf("по комоненту control вернулся не control а %+v\n", comp)
		}
		*c = components.ControlComponent{}
	}
}

func (s *ControlSystem) drainEvents(comps map[types.Entity]types.Component) {

loop:
	for {
		select {
		case e := <-s.eventQueue:
			c, ok := comps[e.ID].(*components.ControlComponent)
			if !ok {
				continue
			}
			c.Superimpose(&e.PlayerState)
		default:
			break loop
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
		fmt.Printf("Преполена очередь %v\n", s)
	}

	return nil
}
