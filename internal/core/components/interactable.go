package components

import (
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type publisher interface {
	Publish(events.Event) error
}
type Interactionor func(types.Entity, publisher)

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

func InteractTerminal(id types.Entity, publisher publisher) {
	e := events.NewInteractTerminalEvent(id)
	publisher.Publish(e)
}

func InteractDoor(id types.Entity, publisher publisher) {
	e := events.NewDoorEvent(id)
	publisher.Publish(e)
}
