package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type CollisionSystem2 struct {
	getter     componentsGetter
	remover    entitiesRemover
	entitiesID []types.Entity
}

func NewCollisionSystem2(getter, remover componentsGetter) *CollisionSystem2 {
	return &CollisionSystem2{
		getter:     getter,
		remover:    remover,
		entitiesID: make([]types.Entity, 0, 128),
	}
}
func (s *CollisionSystem2) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("collider2")

	compsShip := s.getter.GetEntitiesByComponent("navigation_ship")

	var navShip *components.NavigationShipComponent
	var idShip types.Entity
	for id, comp := range compsShip {
		idShip = id
		navShip = comp.(*components.NavigationShipComponent)
	}

	for id := range comps {
		s.entitiesID = append(s.entitiesID, id)
	}

	for i := 0; i < len(s.entitiesID); i++ {
		for j := i + 1; j < len(s.entitiesID); j++ {

			id1 := s.entitiesID[i]
			id2 := s.entitiesID[j]
			collision1, _ := comps[id1].(*components.ColliderComponent2)
			collision2, _ := comps[id2].(*components.ColliderComponent2)

			_, isColliding := collision1.Collide(collision2)
			if isColliding {
				if id1 == idShip {
					//s.getter.GetComponent(id1, "health").(*HealthComponent).Damage(10)
					s.remover.RemoveEntity(id2)
				}
				if id2 == idShip {
					//navShip.Damage()
					s.remover.RemoveEntity(id1)
				}
			}
		}
	}

	return nil
}
