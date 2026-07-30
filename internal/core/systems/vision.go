package systems

import (
	"math"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	walkDistant     = 2.5
	cameDistant     = 0.2
	attackDistant   = 2.8
	alienAttackView = 120
	unreachableMTV  = 0.3
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
		*components.TransformComponent
		*components.VisionComponent
		privMTV geometry.Vec3
	}

	enemies := make([]enemy, 0, len(visionRaw))

	for id, v := range visionRaw {

		if s.HasComponents(id, "transform") {
			v := v.(*components.VisionComponent)

			c, _ := s.GetComponent(id, "transform")
			t := c.(*components.TransformComponent)

			c, _ = s.GetComponent(id, "collider")
			collider := c.(*components.ColliderComponent)

			enemies = append(enemies, enemy{
				t,
				v,
				collider.PrivMTV,
			})
		}

	}

	targetsRaw := s.GetEntitiesByComponent("target")

	type target struct {
		*components.HitboxComponent
		id  types.Entity
		typ components.Target
	}

	targets := make([]*target, 0, len(targetsRaw))

	for id := range targetsRaw {
		if s.HasComponents(id, "hitbox") {
			c, _ := s.GetComponent(id, "hitbox")
			hitbox := c.(*components.HitboxComponent)

			var typ components.Target

			if s.HasComponents(id, "baggage") {
				typ = components.Baggage
			} else {
				typ = components.Player
			}

			targets = append(targets,
				&target{
					hitbox,
					id,
					typ,
				},
			)
		}
	}

	ray := geometry.NewRayCollider(geometry.Vec3{}, geometry.Vec3{}, 0)

	for _, enemy := range enemies {

		wantDistant := geometry.Vec3{X: math.MaxFloat64}
		var bestTarget *target

		for _, target := range targets {

			canSee := true
			targetPos := target.Collider.GetCenter()
			ray.Change(enemy.Position, targetPos)

			if target.typ == components.Player {
				for _, collider := range colliders {

					result := ray.Collide(collider.Collider)

					if result.HasCollision {
						canSee = false
						break
					}

				}
			}

			if !canSee {
				continue
			}

			targetDistant := targetPos.Sub(enemy.Position)
			angle := geometry.AngleBetweenDegrees(wantDistant.Normalize(), targetDistant.Normalize())
			canSee = canSee && angle <= alienAttackView

			switch {
			case canSee && bestTarget == nil:
				bestTarget = target
			case canSee && target.typ == components.Player:

				bestDistant := bestTarget.Collider.GetCenter().Sub(enemy.Position).Length()
				currDistant := targetDistant.Length()

				if currDistant < bestDistant {
					wantDistant = bestTarget.Collider.GetCenter().Sub(enemy.Position)
					wantDistant.Y = 0

					bestTarget = target
				}

			case canSee && target.typ == components.Player && (bestTarget.typ == components.Nothing || bestTarget.typ == components.Baggage):
				wantDistant = targetDistant
				wantDistant.Y = 0

				bestTarget = target

			case target.typ == components.Baggage && bestTarget.typ == components.Nothing:

				bestDistant := bestTarget.Collider.GetCenter().Sub(enemy.Position).Length()
				currDistant := targetDistant.Length()

				if currDistant < bestDistant {
					bestTarget = target

				}
			}
		}

		switch {

		case bestTarget.typ == components.Player:
			enemy.GoingToLastSee = false
			enemy.WantPossition = bestTarget.Collider.GetCenter()
			enemy.Target = components.Player

		case enemy.Target == components.Player:
			wantDistant = bestTarget.Collider.GetCenter().Sub(enemy.Position)
			wantDistant.Y = 0

			if enemy.WantPossition.Sub(enemy.Position).Length() > cameDistant && wantDistant.Normalize().Sub(enemy.privMTV.Normalize()).Length() > unreachableMTV {
				enemy.GoingToLastSee = true
			} else {
				enemy.GoingToLastSee = false
				enemy.Target = components.Nothing
			}
		default:
			enemy.GoingToLastSee = false
			enemy.WantPossition = bestTarget.Collider.GetCenter()
			enemy.Target = bestTarget.typ
		}
	}

	return nil
}
