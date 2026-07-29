package systems

import (
	"math"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	minDistant      = 0.6
	minCoord        = 0.2
	alienAttackView = 120
	unreachableMTV  = 0.4
)

type VisionSystem struct {
	componentsGetter
}

func NewVisionSystem(getter componentsGetter) *VisionSystem {
	return &VisionSystem{
		getter,
	}
}

func (s *VisionSystem) Update(dt float32) error {

	collidersRaw := s.GetEntitiesByComponent("collider")

	colliders := make([](*components.ColliderComponent), 0, len(collidersRaw))

	for id, collider := range collidersRaw {
		if !s.HasComponents(id, "movable") {
			c, _ := collider.(*components.ColliderComponent)

			if c.Enable {
				colliders = append(colliders, c)
			}

		}

	}

	visionRaw := s.GetEntitiesByComponent("vision")

	type enemy struct {
		transform *components.TransformComponent
		vision    *components.VisionComponent
		privMTV   geometry.Vec3
	}

	enemies := make([]enemy, 0, len(visionRaw))

	for id, v := range visionRaw {

		if s.HasComponents(id, "transform") {
			v := v.(*components.VisionComponent)
			v.CanSee = false

			c, _ := s.GetComponent(id, "transform")
			t := c.(*components.TransformComponent)

			c, _ = s.GetComponent(id, "collider")
			collider := c.(*components.ColliderComponent)

			enemies = append(enemies, enemy{
				transform: t,
				vision:    v,
				privMTV:   collider.PrivMTV,
			})
		}

	}

	targetsRaw := s.GetEntitiesByComponent("target")

	type target struct {
		hitbox *components.HitboxComponent
		id     types.Entity
	}

	targets := make([]*target, 0, len(targetsRaw))

	for id := range targetsRaw {
		if s.HasComponents(id, "hitbox") {
			c, _ := s.GetComponent(id, "hitbox")
			hitbox := c.(*components.HitboxComponent)

			targets = append(targets,
				&target{
					hitbox,
					id,
				},
			)
		}
	}

	ray := geometry.NewRayCollider(geometry.Vec3{}, geometry.Vec3{}, 0)
	var canSee bool

	for _, enemy := range enemies {

		enemy.vision.CanSee = false
		enemy.vision.Distant = geometry.Vec3{X: math.MaxFloat64}

		for _, target := range targets {
			targetPos := target.hitbox.Collider.GetCenter()
			ray.Change(enemy.transform.Position, targetPos)

			canSee = true
			var enemyTarget components.Target
			if s.HasComponents(target.id, "baggage") {
				enemyTarget = components.Baggage
			} else {
				enemyTarget = components.Player
			}

			if enemyTarget == components.Player {
				for _, collider := range colliders {

					result := ray.Collide(collider.Collider)

					if result.HasCollision {
						canSee = false
						break
					}

				}
			}

			targetDistant := targetPos.Sub(enemy.transform.Position)
			angle := geometry.AngleBetweenDegrees(enemy.vision.Distant.Normalize(), targetDistant.Normalize())
			canSee = canSee && angle <= alienAttackView
			lastSeenDistant := enemy.vision.LastSeen.Sub(enemy.transform.Position)
			lastSeenDistant.Y = 0

			switch {
			case canSee && enemyTarget == components.Player && enemy.vision.Target == components.Player:
				if targetDistant.Length() < enemy.vision.Distant.Length() {
					enemy.vision.CanSee = true
					enemy.vision.Distant = targetDistant
					enemy.vision.LastSeen = targetPos
					enemy.vision.Target = enemyTarget
				}

			case enemy.vision.Target == components.Player:

				if lastSeenDistant.Length() > minCoord && lastSeenDistant.Normalize().Add(enemy.privMTV.Normalize()).Length() > unreachableMTV {
					enemy.vision.CanSee = true
					enemy.vision.Distant = lastSeenDistant
				} else {
					enemy.vision.CanSee = false
					enemy.vision.Target = components.Nothing
				}
			default:
				enemy.vision.CanSee = true
				enemy.vision.Distant = targetDistant
				enemy.vision.LastSeen = targetPos
				enemy.vision.Target = enemyTarget
			}

			enemy.vision.Distant.Y = 0
			enemy.vision.Distant = enemy.vision.Distant

		}

	}

	return nil
}
