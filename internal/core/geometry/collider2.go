package geometry

type CollisionResult2 struct {
	HasCollision bool
	MTV          Vec2
}

type Collider2 interface {
	Collide(other Collider2) CollisionResult2
}

type TriangleCollider struct {
	Vertex1 Vec2
	Vertex2 Vec2
	Vertex3 Vec2
}

func NewTriangleCollider(v1, v2, v3 Vec2) *TriangleCollider {
	return &TriangleCollider{Vertex1: v1, Vertex2: v2, Vertex3: v3}
}

func (b *TriangleCollider) Collide(other Collider2) CollisionResult2 {
	switch o := other.(type) {
	case *CircleCollider:
		return triangleCircleCollide(b, o)
	default:
		return CollisionResult2{HasCollision: false}
	}
}

type CircleCollider struct {
	Center Vec2
	Radius float64
}

func NewCircleCollider(center Vec2, radius float64) *CircleCollider {
	return &CircleCollider{
		Center: center,
		Radius: radius,
	}
}

func (c *CircleCollider) Collide(other Collider2) CollisionResult2 {
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
		return CollisionResult2{HasCollision: false}
	}
}

func triangleCircleCollide(a *TriangleCollider, b *CircleCollider) CollisionResult2 {
	v1, v2 := twoClosestPointsToTarget(b.Center, a.Vertex1, a.Vertex2, a.Vertex3)

	h := heightFromA(b.Center, v1, v2)

	overlapX := b.Radius - abs(h.X)
	if overlapX <= 0 {
		return CollisionResult2{HasCollision: false}
	}

	overlapY := b.Radius - abs(h.Y)
	if overlapY <= 0 {
		return CollisionResult2{HasCollision: false}
	}

	minOverlap := overlapX
	mtv := Vec2{X: overlapX, Y: 0}
	if h.X < 0 {
		mtv.X = -mtv.X
	}

	if overlapY < minOverlap {
		minOverlap = overlapY
		mtv = Vec2{X: 0, Y: overlapY}
		if h.Y < 0 {
			mtv.Y = -mtv.Y
		}
	}

	return CollisionResult2{HasCollision: true, MTV: mtv}
}

func circleCircleCollide(a, b *CircleCollider) CollisionResult2 {
	d := a.Center.Sub(b.Center)

	overlapX := (a.Radius + b.Radius) - abs(d.X)
	if overlapX <= 0 {
		return CollisionResult2{HasCollision: false}
	}

	overlapY := (a.Radius + b.Radius) - abs(d.Y)
	if overlapY <= 0 {
		return CollisionResult2{HasCollision: false}
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

	return CollisionResult2{HasCollision: true, MTV: mtv}
}
