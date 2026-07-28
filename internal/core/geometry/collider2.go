package geometry

type TriangleCollider struct {
	Vertex1 Vec3
	Vertex2 Vec3
	Vertex3 Vec3
	Center  Vec3
}

func NewTriangleCollider(v1, v2, v3 Vec3) *TriangleCollider {
	c := &TriangleCollider{Vertex1: v1, Vertex2: v2, Vertex3: v3}
	c.Center = c.GetCenter()
	return c
}

func (b *TriangleCollider) Collide(other Collider) CollisionResult {
	switch o := other.(type) {
	case *CircleCollider:
		return triangleCircleCollide(b, o)
	default:
		return CollisionResult{HasCollision: false}
	}
}

func (c *TriangleCollider) GetCenter() Vec3 {
	return Vec3{
		X: average(c.Vertex1.X, c.Vertex2.X, c.Vertex3.X),
		Y: average(c.Vertex1.Y, c.Vertex2.Y, c.Vertex3.Y),
		Z: 0,
	}
}

func (c *TriangleCollider) ChangeCenter(center Vec3) {
	change := center.Sub(c.Center)
	c.Center = center
	c.Vertex1.Add(change)
	c.Vertex2.Add(change)
	c.Vertex3.Add(change)
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

func (c *CircleCollider) ChangeCenter(center Vec3) {
	c.Center = center
}

func (c *CircleCollider) GetCenter() Vec3 {
	return c.Center
}

func triangleCircleCollide(triangle *TriangleCollider, circle *CircleCollider) CollisionResult {
	dist, closestPoint, normal := pointTriangleDistance(
		circle.Center,
		triangle.Vertex1,
		triangle.Vertex2,
		triangle.Vertex3,
	)

	if dist > circle.Radius {
		return CollisionResult{HasCollision: false}
	}

	dir := closestPoint.Sub(circle.Center)
	dirLen := dir.Length()

	var mtv Vec3
	if dirLen < 0.0001 {
		mtv = normal.Scale(circle.Radius)
	} else {
		dirNorm := dir.Scale(1.0 / dirLen)
		overlap := circle.Radius - dist
		mtv = dirNorm.Scale(overlap)
	}

	return CollisionResult{
		HasCollision: true,
		MTV:          mtv,
	}
}

func circleCircleCollide(a, b *CircleCollider) CollisionResult {
	d := a.Center.Sub(b.Center)
	dist := d.Length()

	if dist >= a.Radius+b.Radius {
		return CollisionResult{HasCollision: false}
	}

	if dist < 0.0001 {
		return CollisionResult{
			HasCollision: true,
			MTV:          Vec3{X: 0, Y: a.Radius, Z: 0},
		}
	}

	overlap := (a.Radius + b.Radius) - dist
	direction := d.Normalize()
	mtv := direction.Scale(overlap)

	return CollisionResult{
		HasCollision: true,
		MTV:          mtv,
	}
}
