package systems

import (
	"fmt"

	//"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	//"github.com/ValeraSun/PathEater/internal/core/types"
)

type AsteroidSystem struct {
	getter     componentsGetter
	eventQueue chan *events.MeteoriteZoneEvent
}

func NewAsteroidSystem(getter componentsGetter) *AsteroidSystem {
	s := &AsteroidSystem{
		getter:     getter,
		eventQueue: make(chan *events.MeteoriteZoneEvent, 100),
	}
	subscriber.Subscribe("meteoriteZone", s.OnEvent)
	return s
}

func (*AsteroidSystem) Update(dt float32) error {
	return nil
}

func (s *AsteroidSystem) OnEvent(event events.Event) error {
	mz, ok := event.(*events.MeteoriteZoneEvent)

	if !ok {
		return nil
	}

	select {
	case s.eventQueue <- mz:

	default:
		fmt.Printf("Переполена очередь %v\n", s)
	}

	return nil
}
