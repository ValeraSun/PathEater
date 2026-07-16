package geometry

import (
	"math"
)

type CollisionResult struct {
	HasCollision bool
	MTV          Vec2 // Вектор, на который нужно сдвинуть объект, у которого вызван метод, чтобы выйти из коллизии
}

type Collider interface {
	Collide(other Collider) CollisionResult
}

type TriangleCollider struct {
	Vertex1 Vec2
	Vertex2 Vec2
	Vertex3 Vec2
}

func NewTriangleCollider(v1, v2, v3 Vec2) *TriangleCollider {
	return &TriangleCollider{Vertex1: v1, Vertex2: v2, Vertex3: v3}
}

func (b *TriangleCollider) Collide(other Collider) CollisionResult {
	switch o := other.(type) {
	case *CircleCollider:
		return triangleCircleCollide(b, o)
	default:
		return CollisionResult{HasCollision: false}
	}
}

type CircleCollider struct {
	Center     Vec2
	Radius     float64
}

func NewCircleCollider(center Vec2, radius float64) *CircleCollider {
	return &CircleCollider{
		Center: center,
		Radius: radius,
	}
}

func (c *CapsuleCollider) Collide(other Collider) CollisionResult {
	switch o := other.(type) {
	case *TriangleCollider:
		res := triangleCircleCollide(o, c)
		if res.HasCollision {
			res.MTV = res.MTV.Scale(-1)
		}
		return res
	case *CircleCollider:
		return circleCircleCollide(c, o)
	default:
		return CollisionResult{HasCollision: false}
	}
}

func triangleCircleCollide(a *TriangleCollider, b *CircleCollider) CollisionResult {
	v1, v2 := twoClosestPointsToTarget(b.Center, a.V1, a.V2, a.V3)

	h := heightFromA(b.center, v1, v2)

	overlapX := b.radius - abs(h.X)
	if overlapX <= 0 {
		return CollisionResult{HasCollision: false}
	}

	overlapY := b.radius - abs(h.Y)
	if overlapY <= 0 {
		return CollisionResult{HasCollision: false}
	}

	minOverlap := overlapX
	mtv := Vec3{X: overlapX, Y: 0, Z: 0}
	if d.X < 0 {
		mtv.X = -mtv.X
	}

	if overlapY < minOverlap {
		minOverlap = overlapY
		mtv = Vec3{X: 0, Y: overlapY, Z: 0}
		if d.Y < 0 {
			mtv.Y = -mtv.Y
		}
	}

	if overlapZ < minOverlap {
		mtv = Vec3{X: 0, Y: 0, Z: overlapZ}
		if d.Z < 0 {
			mtv.Z = -mtv.Z
		}
	}

	return CollisionResult{HasCollision: true, MTV: mtv}
}

func circleCircleCollide(a, b *CircleCollider) CollisionResult {
	d := a.Center.Sub(b.Center)

	overlapX := (a.radius + b.radius) - abs(d.X)
	if overlapX <= 0 {
		return CollisionResult{HasCollision: false}
	}

	overlapY := (a.radius + b.radius) - abs(d.Y)
	if overlapY <= 0 {
		return CollisionResult{HasCollision: false}
	}

	// Ищем минимальное перекрытие для выбора оси выталкивания
	minOverlap := overlapX
	mtv := Vec2{X: overlapX, Y: 0}
	if d.X < 0 {
		mtv.X = -mtv.X
	}

	if overlapY < minOverlap {
		minOverlap = overlapY
		mtv = Vec2{X: 0, Y: overlapY}
		if d.Y < 0 {
			mtv.Y = -mtv.Y
		}
	}

	return CollisionResult{HasCollision: true, MTV: mtv}
}