package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type InteractionSystem struct {
	componentsGetter
	publisher
}

func NewInteractionSystem(getter componentsGetter, publisher publisher) *InteractionSystem {
	s := &InteractionSystem{
		getter,
		publisher,
	}

	return s
}

func (s *InteractionSystem) Update(dt float32) error {
	comps := s.GetEntitiesByComponent("interactable")

	type interactable struct {
		*components.InteractableComponent
		types.Entity
	}
	interactables := make([]*interactable, 0, len(comps))

	for id, c := range comps {
		inter := c.(*components.InteractableComponent)
		interactables = append(interactables, &interactable{
			inter,
			id,
		},
		)
	}

	comps = s.GetEntitiesByComponent("interactionDetector")

	for source, c := range comps {

		if !s.HasComponents(source, "transform", "control") {
			continue
		}

		interactor := c.(*components.InteractionDetectorComponent)

		c, _ = s.GetComponent(source, "transform")
		transform := c.(*components.TransformComponent)

		ray := interactor.GetRay(transform.Position, transform.Direction)

		for _, inter := range interactables {
			r := ray.Collide(&inter.Collider)

			if r.HasCollision {

				c, _ = s.GetComponent(source, "control")
				control := c.(*components.ControlComponent)

				if inter.OnKey == "down" && control.InteractDown {
					s.Publish(inter.Interaction(source, inter.Entity))
				}

				if inter.OnKey == "always" && control.Interact {
					s.Publish(inter.Interaction(source, inter.Entity))
				}

			}
		}

	}
	return nil
}
