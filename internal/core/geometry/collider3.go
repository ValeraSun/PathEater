package geometry

import (
	"math"
)

type CollisionResult struct {
	HasCollision bool
	MTV          Vec3
}

type Collider interface {
	Collide(other Collider) CollisionResult
	ChangeCenter(center Vec3)
	GetCenter() Vec3
}

type BoxCollider struct {
	Center      Vec3
	HalfExtents Vec3
	Axes        [3]Vec3
}

func NewBoxCollider(center, halfExtents Vec3, axes [3]Vec3) *BoxCollider {
	return &BoxCollider{
		Center:      center,
		HalfExtents: halfExtents,
		Axes: [3]Vec3{
			axes[0].Normalize(),
			axes[1].Normalize(),
			axes[2].Normalize(),
		},
	}
}

func (b *BoxCollider) ChangeCenter(center Vec3) {
	b.Center = center
}

func (b *BoxCollider) GetCenter() Vec3 {
	return b.Center
}

func (b *BoxCollider) GetNormalByAxis(index int) Vec3 {
	if index < 0 || index > 2 {
		return Vec3{0, 0, 0}
	}
	return b.Axes[index]
}

// Возвращает нормаль, направленную в center
func (b *BoxCollider) GetInwardNormal(center Vec3) Vec3 {
	rawNormal := b.GetThinnestNormal()

	toCenter := b.Center.Sub(center)

	if rawNormal.Dot(toCenter) < 0 {
		return rawNormal
	}
	return rawNormal.Scale(-1)
}

func (b *BoxCollider) GetThinnestNormal() Vec3 {
	minAxis := 0
	minVal := b.HalfExtents.X
	if b.HalfExtents.Y < minVal {
		minVal = b.HalfExtents.Y
		minAxis = 1
	}
	if b.HalfExtents.Z < minVal {
		minAxis = 2
	}
	return b.Axes[minAxis]
}

func (b *BoxCollider) Collide(other Collider) CollisionResult {
	switch o := other.(type) {
	case *BoxCollider:
		return boxBoxCollide(b, o)
	case *CapsuleCollider:
		return boxCapsuleCollide(b, o)
	case *RayCollider:
		res := rayBoxCollide(o, b)
		if res.HasCollision {
			res.MTV = res.MTV.Scale(-1)
		}
		return res
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

func (c *CapsuleCollider) ChangeCenter(center Vec3) {
	c.Center = center
}

func (c *CapsuleCollider) GetCenter() Vec3 {
	return c.Center
}

func (c *CapsuleCollider) Collide(other Collider) CollisionResult {
	switch o := other.(type) {
	case *BoxCollider:
		res := boxCapsuleCollide(o, c)
		if res.HasCollision {
			res.MTV = res.MTV.Scale(-1)
		}
		return res
	case *CapsuleCollider:
		return capsuleCapsuleCollide(c, o)
	case *RayCollider:

		res := rayCapsuleCollide(o, c)
		if res.HasCollision {
			res.MTV = res.MTV.Scale(-1)
		}
		return res
	default:
		return CollisionResult{HasCollision: false}
	}
}

func boxBoxCollide(a, b *BoxCollider) CollisionResult {

	axes := make([]Vec3, 0, 15)
	axes = append(axes, a.Axes[0], a.Axes[1], a.Axes[2])
	axes = append(axes, b.Axes[0], b.Axes[1], b.Axes[2])

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			cross := a.Axes[i].Cross(b.Axes[j])

			if cross.LengthSqr() > 1e-9 {
				axes = append(axes, cross.Normalize())
			}
		}
	}

	minOverlap := math.MaxFloat64
	var mtvAxis Vec3
	d := a.Center.Sub(b.Center)

	for _, axis := range axes {
		if axis.LengthSqr() < 1e-9 {
			continue
		}
		axis = axis.Normalize()

		rA := a.HalfExtents.X*math.Abs(a.Axes[0].Dot(axis)) +
			a.HalfExtents.Y*math.Abs(a.Axes[1].Dot(axis)) +
			a.HalfExtents.Z*math.Abs(a.Axes[2].Dot(axis))

		rB := b.HalfExtents.X*math.Abs(b.Axes[0].Dot(axis)) +
			b.HalfExtents.Y*math.Abs(b.Axes[1].Dot(axis)) +
			b.HalfExtents.Z*math.Abs(b.Axes[2].Dot(axis))

		dist := math.Abs(d.Dot(axis))

		overlap := rA + rB - dist
		if overlap <= 0 {
			return CollisionResult{HasCollision: false}
		}

		if overlap < minOverlap {
			minOverlap = overlap
			mtvAxis = axis
		}
	}

	if d.Dot(mtvAxis) < 0 {
		mtvAxis = mtvAxis.Scale(-1)
	}

	return CollisionResult{HasCollision: true, MTV: mtvAxis.Scale(minOverlap)}
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

		mtv = dir.Normalize().Scale(overlap)
	} else {

		mtv = Vec3{X: 0, Y: overlap, Z: 0}
	}

	return CollisionResult{HasCollision: true, MTV: mtv}
}

func boxCapsuleCollide(box *BoxCollider, cap *CapsuleCollider) CollisionResult {
	a := cap.Center.Sub(cap.Direction.Scale(cap.HalfHeight))
	b := cap.Center.Add(cap.Direction.Scale(cap.HalfHeight))

	axes := []Vec3{
		box.Axes[0],
		box.Axes[1],
		box.Axes[2],
		cap.Direction,
		cap.Direction.Cross(box.Axes[0]),
		cap.Direction.Cross(box.Axes[1]),
		cap.Direction.Cross(box.Axes[2]),
	}

	minOverlap := math.MaxFloat64
	var mtvAxis Vec3

	for _, axis := range axes {
		if axis.LengthSqr() < 1e-9 {
			continue
		}
		axis = axis.Normalize()

		boxProj := box.HalfExtents.X*math.Abs(box.Axes[0].Dot(axis)) +
			box.HalfExtents.Y*math.Abs(box.Axes[1].Dot(axis)) +
			box.HalfExtents.Z*math.Abs(box.Axes[2].Dot(axis))

		boxCenterProj := box.Center.Dot(axis)
		boxMin := boxCenterProj - boxProj
		boxMax := boxCenterProj + boxProj

		projA := a.Dot(axis)
		projB := b.Dot(axis)
		segMin := min(projA, projB) - cap.Radius
		segMax := max(projA, projB) + cap.Radius

		if segMax < boxMin || boxMax < segMin {
			return CollisionResult{HasCollision: false}
		}

		overlap := min(boxMax-segMin, segMax-boxMin)

		if overlap < minOverlap {
			minOverlap = overlap
			mtvAxis = axis
		}
	}

	dir := box.Center.Sub(cap.Center)
	if dir.Dot(mtvAxis) < 0 {
		mtvAxis = mtvAxis.Scale(-1)
	}

	return CollisionResult{HasCollision: true, MTV: mtvAxis.Scale(minOverlap)}
}

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

type RayCollider struct {
	Origin    Vec3
	Direction Vec3
	Length    float64
}

func NewRayCollider(origin, direction Vec3, length float64) *RayCollider {
	return &RayCollider{
		Origin:    origin,
		Direction: direction.Normalize(),
		Length:    length,
	}
}

func (c *RayCollider) Change(start, end Vec3) {
	c.Origin = start
	dir := end.Sub(start)
	c.Direction = dir.Normalize()
	c.Length = dir.Length()

}

func (r *RayCollider) ChangeCenter(center Vec3) {
	r.Origin = center
}

func (r *RayCollider) GetCenter() Vec3 {
	return r.Origin
}

func (r *RayCollider) Collide(other Collider) CollisionResult {
	switch o := other.(type) {
	case *BoxCollider:
		return rayBoxCollide(r, o)
	case *CapsuleCollider:
		return rayCapsuleCollide(r, o)
	case *RayCollider:

		return CollisionResult{HasCollision: false}
	default:
		return CollisionResult{HasCollision: false}
	}
}

func rayBoxCollide(ray *RayCollider, box *BoxCollider) CollisionResult {
	tmin := math.Inf(-1)
	tmax := math.Inf(1)

	delta := box.Center.Sub(ray.Origin)

	he := []float64{box.HalfExtents.X, box.HalfExtents.Y, box.HalfExtents.Z}

	for i := 0; i < 3; i++ {
		axis := box.Axes[i]
		e := axis.Dot(delta)
		f := ray.Direction.Dot(axis)

		if math.Abs(f) > 1e-9 {
			invF := 1.0 / f
			t1 := (e - he[i]) * invF
			t2 := (e + he[i]) * invF

			if t1 > t2 {
				t1, t2 = t2, t1
			}
			if t1 > tmin {
				tmin = t1
			}
			if t2 < tmax {
				tmax = t2
			}
			if tmin > tmax {
				return CollisionResult{HasCollision: false}
			}
			if tmax < 0 {
				return CollisionResult{HasCollision: false}
			}
		} else {

			if -e-he[i] > 0 || -e+he[i] < 0 {
				return CollisionResult{HasCollision: false}
			}
		}
	}

	if tmin > ray.Length {
		return CollisionResult{HasCollision: false}
	}

	hitT := tmin
	if hitT < 0 {
		hitT = 0
	}
	penetration := tmax - hitT
	mtv := ray.Direction.Scale(-penetration)

	return CollisionResult{HasCollision: true, MTV: mtv}
}

func rayCapsuleCollide(ray *RayCollider, cap *CapsuleCollider) CollisionResult {

	rayEnd := ray.Origin.Add(ray.Direction.Scale(ray.Length))

	capA := cap.Center.Sub(cap.Direction.Scale(cap.HalfHeight))
	capB := cap.Center.Add(cap.Direction.Scale(cap.HalfHeight))

	c1, c2 := closestPtSegmentSegment(ray.Origin, rayEnd, capA, capB)

	dir := c1.Sub(c2)
	distSqr := dir.LengthSqr()

	if distSqr > cap.Radius*cap.Radius {
		return CollisionResult{HasCollision: false}
	}

	dist := math.Sqrt(distSqr)
	overlap := cap.Radius - dist

	var mtv Vec3
	if dist > 1e-9 {

		mtv = dir.Normalize().Scale(overlap)
	} else {

		mtv = Vec3{X: 0, Y: overlap, Z: 0}
	}

	return CollisionResult{HasCollision: true, MTV: mtv}
}
