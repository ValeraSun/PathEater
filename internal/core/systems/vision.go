package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

type VisionSystem struct {
	getter componentsGetter
}

func NewVisionSystem(getter componentsGetter) *VisionSystem {
	return &VisionSystem{
		getter: getter,
	}
}

func (s *VisionSystem) Update(dt float32) error {

	collidersRaw := s.getter.GetEntitiesByComponent("collider")

	colliders := make([](*components.ColliderComponent), 0, len(collidersRaw))

	for id, collider := range collidersRaw {
		if !s.getter.HasComponents(id, "movable") {
			c, _ := collider.(*components.ColliderComponent)

			colliders = append(colliders, c)
		}

	}

	visionRaw := s.getter.GetEntitiesByComponent("vision")

	type enemy struct {
		transform *components.TransformComponent
		vision    *components.VisionComponent
	}

	enemies := make([]enemy, 0, len(visionRaw))

	for id, v := range visionRaw {

		if s.getter.HasComponents(id, "transform") {
			v := v.(*components.VisionComponent)
			v.CanSee = false

			c, _ := s.getter.GetComponent(id, "transform")
			t := c.(*components.TransformComponent)

			enemies = append(enemies, enemy{
				transform: t,
				vision:    v,
			})
		}

	}

	playersRaw := s.getter.GetEntitiesByComponent("player")

	players := make([]*components.TransformComponent, 0, len(playersRaw))

	for id := range playersRaw {
		if s.getter.HasComponents(id, "transform") {
			c, _ := s.getter.GetComponent(id, "transform")
			player := c.(*components.TransformComponent)

			players = append(players, player)
		}
	}

	ray := geometry.NewRayCollider(geometry.Vec3{}, geometry.Vec3{}, 0)
	var canSee bool

	for _, enemy := range enemies {
		for _, player := range players {
			ray.Change(enemy.transform.Position, player.Position)
			canSee = true

			for _, collider := range colliders {

				result := ray.Collide(collider.Collider)

				if result.HasCollision {
					canSee = false
					break
				}

			}

			playerDirection := player.Position.Sub(enemy.transform.Position)
			angle := geometry.AngleBetweenDegrees(enemy.vision.Direction, playerDirection)
			canSee = canSee && angle <= 120

			if canSee {
				enemy.vision.CanSee = true
				enemy.vision.Direction = playerDirection
				enemy.vision.LastSeen = player.Position
			} else {

				lastSeenDirection := enemy.vision.LastSeen.Sub(enemy.transform.Position)
				if lastSeenDirection.Length() > 0.2 {
					enemy.vision.CanSee = true
					enemy.vision.Direction = lastSeenDirection
				} else {
					enemy.vision.CanSee = false
				}

			}

			enemy.vision.Direction.Y = 0
			enemy.vision.Direction = enemy.vision.Direction.Normalize()

		}

	}

	return nil
}
