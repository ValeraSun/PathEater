package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type RenderSystem struct {
	getter      componentsGetter
	broadcaster Broadcaster
}

func NewRenderSystem(getter componentsGetter, broadcaster Broadcaster) *RenderSystem {
	return &RenderSystem{
		getter:      getter,
		broadcaster: broadcaster,
	}
}

type playerData struct {
	Position geometry.Vec3 `json:"position"`
	Rotation geometry.Vec3 `json:"rotation"`
	Health   int
}

type shipData struct {
	Collider        geometry.Collider `json:"collider"`
	WeaponDirection geometry.Vec3     `json:"weapon_direction"`
	Health          int               `json:"health"`
}

type asteroidData struct {
	Position  geometry.Vec3 `json:"position"`
	Radius    float64       `json:"radius"`
	Destroyed bool          `json:"destroyed"`
}

func (s *RenderSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("update")

	for id := range comps {
		if s.getter.HasComponents(id, "transform") && s.getter.HasComponents(id, "control") && s.getter.HasComponents(id, "health") {
			c, _ := s.getter.GetComponent(id, "transform")
			transform := c.(*components.TransformComponent)

			c, _ = s.getter.GetComponent(id, "health")
			hp := c.(*components.HealthComponent)

			s.broadcaster.SendEntityUpdate(ecs.EntityInfo{
				ID:   id,
				Type: "player",
				Data: playerData{
					Position: transform.Position,
					Rotation: transform.Direction,
					Health:   hp.Health,
				}})
		}

		if s.getter.HasComponents(id, "ship") {
			c, _ := s.getter.GetComponent(id, "health")
			hp := c.(*components.HealthComponent)

			c, _ = s.getter.GetComponent(id, "collider")
			col := c.(*components.ColliderComponent)

			c, _ = s.getter.GetComponent(id, "weapon")
			weap := c.(*components.WeaponComponent)

			s.broadcaster.SendEntityUpdate(ecs.EntityInfo{
				ID:   id,
				Type: "ship",
				Data: shipData{
					Collider:        col.Collider,
					WeaponDirection: weap.Direction,
					Health:          hp.Health,
				}})
		}

		if s.getter.HasComponents(id, "asteroid") {

			c, _ := s.getter.GetComponent(id, "asteroid")
			aster := c.(*components.AsteroidComponent)

			if aster.Visible {
				c, _ := s.getter.GetComponent(id, "transform")
				transform := c.(*components.TransformComponent)

				c, _ = s.getter.GetComponent(id, "collider")
				collider := c.(*components.ColliderComponent)

				s.broadcaster.SendEntityUpdate(ecs.EntityInfo{
					ID:   id,
					Type: "asteroid",
					Data: asteroidData{
						Position:  transform.Position,
						Radius:    collider.Collider.(*geometry.CircleCollider).GetRadius(),
						Destroyed: aster.Destroyed,
					}})
			}
		}
	}

	return nil
}
