package geometry

import (
	"math"
)

type CollisionResult struct {
	HasCollision bool
	MTV          Vec3 // Вектор, на который нужно сдвинуть объект, у которого вызван метод, чтобы выйти из коллизии
}

type Collider interface {
	Collide(other Collider) CollisionResult
}

type BoxCollider struct {
	Center      Vec3
	HalfExtents Vec3
}

func NewBoxCollider(center, halfExtents Vec3) *BoxCollider {
	return &BoxCollider{Center: center, HalfExtents: halfExtents}
}

func (b *BoxCollider) Collide(other Collider) CollisionResult {
	switch o := other.(type) {
	case *BoxCollider:
		return boxBoxCollide(b, o)
	case *CapsuleCollider:
		return boxCapsuleCollide(b, o)
	default:
		return CollisionResult{HasCollision: false}
	}
}

type CapsuleCollider struct {
	Center     Vec3
	Direction  Vec3
	HalfHeight float64
	Radius     float64
}

func NewCapsuleCollider(center, direction Vec3, halfHeight, radius float64) *CapsuleCollider {
	dir := direction.Normalize()
	return &CapsuleCollider{
		Center:     center,
		Direction:  dir,
		HalfHeight: halfHeight,
		Radius:     radius,
	}
}

func (c *CapsuleCollider) Collide(other Collider) CollisionResult {
	switch o := other.(type) {
	case *BoxCollider:
		// Вызываем проверку Box-Capsule, но инвертируем MTV,
		// так как вызываем со стороны капсулы, и выталкивать нужно именно её
		res := boxCapsuleCollide(o, c)
		if res.HasCollision {
			res.MTV = res.MTV.Scale(-1)
		}
		return res
	case *CapsuleCollider:
		return capsuleCapsuleCollide(c, o)
	default:
		return CollisionResult{HasCollision: false}
	}
}

func boxBoxCollide(a, b *BoxCollider) CollisionResult {
	d := a.Center.Sub(b.Center)

	overlapX := (a.HalfExtents.X + b.HalfExtents.X) - abs(d.X)
	if overlapX <= 0 {
		return CollisionResult{HasCollision: false}
	}

	overlapY := (a.HalfExtents.Y + b.HalfExtents.Y) - abs(d.Y)
	if overlapY <= 0 {
		return CollisionResult{HasCollision: false}
	}

	overlapZ := (a.HalfExtents.Z + b.HalfExtents.Z) - abs(d.Z)
	if overlapZ <= 0 {
		return CollisionResult{HasCollision: false}
	}

	// Ищем минимальное перекрытие для выбора оси выталкивания
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

func capsuleCapsuleCollide(a, b *CapsuleCollider) CollisionResult {
	a1 := a.Center.Sub(a.Direction.Scale(a.HalfHeight))
	a2 := a.Center.Add(a.Direction.Scale(a.HalfHeight))
	b1 := b.Center.Sub(b.Direction.Scale(b.HalfHeight))
	b2 := b.Center.Add(b.Direction.Scale(b.HalfHeight))

	c1, c2 := closestPtSegmentSegment(a1, a2, b1, b2)
	dir := c1.Sub(c2)
	distSqr := dir.LengthSqr()
	radSum := a.Radius + b.Radius

	if distSqr > radSum*radSum {
		return CollisionResult{HasCollision: false}
	}

	dist := math.Sqrt(distSqr)
	overlap := radSum - dist

	var mtv Vec3
	if dist > 1e-9 {
		// Стандартный случай: расталкиваем по вектору между ближайшими точками
		mtv = dir.Normalize().Scale(overlap)
	} else {
		// Центры точно совпали, расталкиваем по произвольной оси (например, Y)
		mtv = Vec3{X: 0, Y: overlap, Z: 0}
	}

	return CollisionResult{HasCollision: true, MTV: mtv}
}

func boxCapsuleCollide(box *BoxCollider, cap *CapsuleCollider) CollisionResult {
	a := cap.Center.Sub(cap.Direction.Scale(cap.HalfHeight))
	b := cap.Center.Add(cap.Direction.Scale(cap.HalfHeight))

	axes := [7]Vec3{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
		cap.Direction,
		cap.Direction.Cross(Vec3{1, 0, 0}),
		cap.Direction.Cross(Vec3{0, 1, 0}),
		cap.Direction.Cross(Vec3{0, 0, 1}),
	}

	minOverlap := math.MaxFloat64
	var mtvAxis Vec3

	for _, axis := range axes {
		if axis.LengthSqr() < 1e-9 {
			continue
		}
		// Нормализация обязательна для правильного расчета глубины проникновения
		axis = axis.Normalize()

		boxProj := abs(box.HalfExtents.X*axis.X) + abs(box.HalfExtents.Y*axis.Y) + abs(box.HalfExtents.Z*axis.Z)
		boxCenterProj := box.Center.Dot(axis)
		boxMin := boxCenterProj - boxProj
		boxMax := boxCenterProj + boxProj

		projA := a.Dot(axis)
		projB := b.Dot(axis)
		segMin := min(projA, projB) - cap.Radius
		segMax := max(projA, projB) + cap.Radius

		// Если найдена разделяющая ось — коллизии нет
		if segMax < boxMin || boxMax < segMin {
			return CollisionResult{HasCollision: false}
		}

		// Вычисляем глубину перекрытия
		overlap := min(boxMax-segMin, segMax-boxMin)

		if overlap < minOverlap {
			minOverlap = overlap
			mtvAxis = axis
		}
	}

	// Определяем направление вектора MTV (должен выталкивать Box из Capsule)
	dir := box.Center.Sub(cap.Center)
	if dir.Dot(mtvAxis) < 0 {
		mtvAxis = mtvAxis.Scale(-1)
	}

	return CollisionResult{HasCollision: true, MTV: mtvAxis.Scale(minOverlap)}
}

// ... Оставшиеся вспомогательные функции (closestPtSegmentSegment) остаются без изменений ...
func closestPtSegmentSegment(p1, q1, p2, q2 Vec3) (Vec3, Vec3) {
	d1 := q1.Sub(p1)
	d2 := q2.Sub(p2)
	r := p1.Sub(p2)
	a := d1.Dot(d1)
	e := d2.Dot(d2)
	f := d2.Dot(r)

	s, t := 0.0, 0.0

	if a <= 1e-9 && e <= 1e-9 {
		return p1, p2
	}
	if a <= 1e-9 {
		t = clamp(f/e, 0.0, 1.0)
	} else {
		c := d1.Dot(r)
		if e <= 1e-9 {
			s = clamp(-c/a, 0.0, 1.0)
			t = 0.0
		} else {
			b := d1.Dot(d2)
			denom := a*e - b*b

			if denom != 0.0 {
				s = clamp((b*f-c*e)/denom, 0.0, 1.0)
			} else {
				s = 0.0
			}

			t = (b*s + f) / e
			if t < 0.0 {
				t = 0.0
				s = clamp(-c/a, 0.0, 1.0)
			} else if t > 1.0 {
				t = 1.0
				s = clamp((b-c)/a, 0.0, 1.0)
			}
		}
	}
	c1 := p1.Add(d1.Scale(s))
	c2 := p2.Add(d2.Scale(t))
	return c1, c2
}
