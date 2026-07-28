package components

import (
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type Interactionor func(types.Entity, types.Entity) events.Event

type InteractableComponent struct {
	Collider    geometry.BoxCollider
	Interaction Interactionor
}

func NewInteractableComponent(center, halfExtents geometry.Vec3, interaction Interactionor) *InteractableComponent {
	return &InteractableComponent{
		Collider: *geometry.NewBoxCollider(
			center,
			halfExtents,
			[3]geometry.Vec3{
				{X: 1},
				{Y: 1},
				{Z: 1},
			},
		),
		Interaction: interaction,
	}
}

func (*InteractableComponent) Type() string { return "interactable" }

func InteractTerminal(source, target types.Entity) events.Event {
	return events.NewInteractTerminalEvent(source)

}

func InteractDoor(source, target types.Entity) events.Event {
	return events.NewDoorEvent(target)

}
