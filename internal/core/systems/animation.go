package systems

import (
	"errors"
	"fmt"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type AnimationSystem struct {
	componentsGetter
	publisher
	eventQueue chan *events.AnimationEvent
}

func NewAnimationSystem(getter componentsGetter, subscriber subscriber, publisher publisher) *AnimationSystem {
	s := &AnimationSystem{
		getter,
		publisher,
		make(chan *events.AnimationEvent, 20),
	}
	subscriber.Subscribe("animation", s.OnEvent)
	return s

}

func (s *AnimationSystem) Update(dt float32) error {
	s.resetCooldown(dt)
	err := s.drainEvents()
	return err
}

func (s *AnimationSystem) resetCooldown(dt float32) {
	comps := s.GetEntitiesByComponent("animation")

	for _, c := range comps {
		anime := c.(*components.AnimationComponent)

		anime.ResetCooldown(dt)
	}
}
func (s *AnimationSystem) drainEvents() error {
	for {
		select {
		case e := <-s.eventQueue:
			err := s.handleEvent(e)

			if err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

func (s *AnimationSystem) handleEvent(event *events.AnimationEvent) error {
	c, exist := s.GetComponent(event.ID, "animation")

	if !exist {
		return errors.New("Нет компонента animation")
	}

	anime := c.(*components.AnimationComponent)

	anime.ChangeState(event)
	return nil
}

func (s *AnimationSystem) OnEvent(event events.Event) error {
	ps, ok := event.(*events.AnimationEvent)

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
