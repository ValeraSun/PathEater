package components

import (
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type InteractableComponent struct {
	Collider    geometry.BoxCollider
	OnKey       string
	Interaction func(types.Entity, types.Entity) events.Event
}

func NewInteractableComponent(center, halfExtents geometry.Vec3, onKey string, interaction func(types.Entity, types.Entity) events.Event) *InteractableComponent {
	return &InteractableComponent{
		Collider: *geometry.NewBoxCollider(
			center,
			halfExtents,
			geometry.GetStandartAxes(),
		),
		OnKey:       onKey,
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
