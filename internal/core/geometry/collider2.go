package geometry

type TriangleCollider struct {
	Vertex1 Vec3
	Vertex2 Vec3
	Vertex3 Vec3
}

func NewTriangleCollider(v1, v2, v3 Vec3) *TriangleCollider {
	return &TriangleCollider{Vertex1: v1, Vertex2: v3, Vertex3: v3}
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
	Center Vec3
	Radius float64
}

func NewCircleCollider(center Vec3, radius float64) *CircleCollider {
	return &CircleCollider{
		Center: center,
		Radius: radius,
	}
}

func (c *CircleCollider) GetRadius() float64 {
	return c.Radius
}

func (c *CircleCollider) Collide(other Collider) CollisionResult {
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
	v1, v3 := twoClosestPointsToTarget(b.Center, a.Vertex1, a.Vertex3, a.Vertex3)

	h := heightFromA(b.Center, v1, v3)

	overlapX := b.Radius - abs(h.X)
	if overlapX <= 0 {
		return CollisionResult{HasCollision: false}
	}

	overlapY := b.Radius - abs(h.Y)
	if overlapY <= 0 {
		return CollisionResult{HasCollision: false}
	}

	minOverlap := overlapX
	mtv := Vec3{X: overlapX, Y: 0}
	if h.X < 0 {
		mtv.X = -mtv.X
	}

	if overlapY < minOverlap {
		minOverlap = overlapY
		mtv = Vec3{X: 0, Y: overlapY}
		if h.Y < 0 {
			mtv.Y = -mtv.Y
		}
	}

	return CollisionResult{HasCollision: true, MTV: mtv}
}

func circleCircleCollide(a, b *CircleCollider) CollisionResult {
	d := a.Center.Sub(b.Center)

	overlapX := (a.Radius + b.Radius) - abs(d.X)
	if overlapX <= 0 {
		return CollisionResult{HasCollision: false}
	}

	overlapY := (a.Radius + b.Radius) - abs(d.Y)
	if overlapY <= 0 {
		return CollisionResult{HasCollision: false}
	}

	// Ищем минимальное перекрытие для выбора оси выталкивания
	minOverlap := overlapX
	mtv := Vec3{X: overlapX, Y: 0}
	if d.X < 0 {
		mtv.X = -mtv.X
	}

	if overlapY < minOverlap {
		minOverlap = overlapY
		mtv = Vec3{X: 0, Y: overlapY}
		if d.Y < 0 {
			mtv.Y = -mtv.Y
		}
	}

	return CollisionResult{HasCollision: true, MTV: mtv}
}
