package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type InteractionSystem struct {
	getter     componentsGetter
	publisher  publisher
	eventQueue chan events.Event
}

func NewInteractionSystem(getter componentsGetter, publisher publisher) *InteractionSystem {
	s := &InteractionSystem{
		getter:    getter,
		publisher: publisher,
	}

	return s
}

func (s *InteractionSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("interactable")
	interactables := make([]*components.InteractableComponent, 0, len(comps))

	for _, c := range comps {
		interactables = append(interactables, c.(*components.InteractableComponent))
	}

	comps = s.getter.GetEntitiesByComponent("interactionDetector")

	for id, c := range comps {

		if !s.getter.HasComponents(id, "transform", "control") {
			continue
		}

		interactor := c.(*components.InteractionDetectorComponent)

		c, _ = s.getter.GetComponent(id, "transform")
		transform := c.(*components.TransformComponent)

		ray := interactor.GetRay(transform.Position, transform.Direction)

		for _, col := range interactables {
			r := ray.Collide(&col.Collider)

			if r.HasCollision {

				c, _ = s.getter.GetComponent(id, "control")
				control := c.(*components.ControlComponent)

				if control.Interact {
					col.Interaction(id, s.publisher)
				}
			}
		}

	}
	return nil
}
