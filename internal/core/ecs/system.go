package ecs

import (
	"errors"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type moveSystem struct {
	world *World
}

func NewMoveSystem(world *World) *moveSystem {
	return &moveSystem{
		world: world,
	}
}
func (*moveSystem) Update(dt float32) error {
	return nil
}

func (s *moveSystem) onEvent(event events.Event) error {
	moveEvent, ok := event.(*events.MoveEvent)

	if !ok {
		return nil // Не наше событие, игнорируем
	}

	entity := moveEvent.Id

	component, ok := s.world.GetComponent(entity, "transform")
	transformComponent, ok := component.(*components.TransformComponent)
	if !ok {
		return errors.New("не содержит компонет необходимый")
	}

	transformComponent.Direction = moveEvent.Direction
	transformComponent.Position = moveEvent.Position

	return nil
}
