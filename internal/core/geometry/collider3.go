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
	Axes        [3]Vec3 // Локальные оси OBB (X, Y, Z). Должны быть ортогональны и нормализованы.
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

func (b *BoxCollider) Collide(other Collider) CollisionResult {
	switch o := other.(type) {
	case *BoxCollider:
		return boxBoxCollide(b, o) //
	case *CapsuleCollider:
		return boxCapsuleCollide(b, o) //
	case *RayCollider:
		// Инвертируем MTV, чтобы выталкивался Box, а не Ray
		res := rayBoxCollide(o, b) //
		if res.HasCollision {
			res.MTV = res.MTV.Scale(-1) //
		}
		return res
	default:
		return CollisionResult{HasCollision: false} //[cite: 1]
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
		res := boxCapsuleCollide(o, c)
		if res.HasCollision {
			res.MTV = res.MTV.Scale(-1)
		}
		return res
	case *CapsuleCollider:
		return capsuleCapsuleCollide(c, o)
	case *RayCollider:
		// Инвертируем MTV
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
	// Подготавливаем 15 потенциальных разделяющих осей
	axes := make([]Vec3, 0, 15)
	axes = append(axes, a.Axes[0], a.Axes[1], a.Axes[2])
	axes = append(axes, b.Axes[0], b.Axes[1], b.Axes[2])

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			cross := a.Axes[i].Cross(b.Axes[j])
			// Избегаем нулевых векторов при параллельных осях
			if cross.LengthSqr() > 1e-9 {
				axes = append(axes, cross.Normalize())
			}
		}
	}

	minOverlap := math.MaxFloat64
	var mtvAxis Vec3
	d := a.Center.Sub(b.Center) // Вектор от центра B к центру A[cite: 1]

	for _, axis := range axes {
		if axis.LengthSqr() < 1e-9 {
			continue //[cite: 1]
		}
		axis = axis.Normalize()

		// Проецируем полуразмеры (HalfExtents) коробки A на ось
		rA := a.HalfExtents.X*math.Abs(a.Axes[0].Dot(axis)) +
			a.HalfExtents.Y*math.Abs(a.Axes[1].Dot(axis)) +
			a.HalfExtents.Z*math.Abs(a.Axes[2].Dot(axis))

		// Проецируем полуразмеры (HalfExtents) коробки B на ось
		rB := b.HalfExtents.X*math.Abs(b.Axes[0].Dot(axis)) +
			b.HalfExtents.Y*math.Abs(b.Axes[1].Dot(axis)) +
			b.HalfExtents.Z*math.Abs(b.Axes[2].Dot(axis))

		// Расстояние между центрами, спроецированное на ось
		dist := math.Abs(d.Dot(axis))

		// Если расстояние между центрами больше суммы спроецированных радиусов, то коллизии нет
		overlap := rA + rB - dist
		if overlap <= 0 {
			return CollisionResult{HasCollision: false}
		}

		// Ищем минимальное перекрытие для выбора оси выталкивания[cite: 1]
		if overlap < minOverlap {
			minOverlap = overlap
			mtvAxis = axis
		}
	}

	// Направление вектора MTV должно выталкивать A из B[cite: 1]
	if d.Dot(mtvAxis) < 0 {
		mtvAxis = mtvAxis.Scale(-1) //[cite: 1]
	}

	return CollisionResult{HasCollision: true, MTV: mtvAxis.Scale(minOverlap)} //[cite: 1]
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
	a := cap.Center.Sub(cap.Direction.Scale(cap.HalfHeight)) //[cite: 1]
	b := cap.Center.Add(cap.Direction.Scale(cap.HalfHeight)) //[cite: 1]

	// Тестируем локальные оси OBB, ось капсулы и их попарные векторные произведения[cite: 1]
	axes := []Vec3{
		box.Axes[0],
		box.Axes[1],
		box.Axes[2],
		cap.Direction,                    //[cite: 1]
		cap.Direction.Cross(box.Axes[0]), //[cite: 1]
		cap.Direction.Cross(box.Axes[1]), //[cite: 1]
		cap.Direction.Cross(box.Axes[2]), //[cite: 1]
	}

	minOverlap := math.MaxFloat64 //[cite: 1]
	var mtvAxis Vec3              //[cite: 1]

	for _, axis := range axes {
		if axis.LengthSqr() < 1e-9 { //[cite: 1]
			continue //[cite: 1]
		}
		axis = axis.Normalize() //[cite: 1]

		// Проецируем OBB на ось
		boxProj := box.HalfExtents.X*math.Abs(box.Axes[0].Dot(axis)) +
			box.HalfExtents.Y*math.Abs(box.Axes[1].Dot(axis)) +
			box.HalfExtents.Z*math.Abs(box.Axes[2].Dot(axis)) //[cite: 1]

		boxCenterProj := box.Center.Dot(axis) //[cite: 1]
		boxMin := boxCenterProj - boxProj     //[cite: 1]
		boxMax := boxCenterProj + boxProj     //[cite: 1]

		projA := a.Dot(axis)                     //[cite: 1]
		projB := b.Dot(axis)                     //[cite: 1]
		segMin := min(projA, projB) - cap.Radius //[cite: 1]
		segMax := max(projA, projB) + cap.Radius //[cite: 1]

		// Если найдена разделяющая ось — коллизии нет[cite: 1]
		if segMax < boxMin || boxMax < segMin { //[cite: 1]
			return CollisionResult{HasCollision: false} //[cite: 1]
		}

		// Вычисляем глубину перекрытия[cite: 1]
		overlap := min(boxMax-segMin, segMax-boxMin) //[cite: 1]

		if overlap < minOverlap { //[cite: 1]
			minOverlap = overlap //[cite: 1]
			mtvAxis = axis       //[cite: 1]
		}
	}

	// Определяем направление вектора MTV (должен выталкивать Box из Capsule)[cite: 1]
	dir := box.Center.Sub(cap.Center) //[cite: 1]
	if dir.Dot(mtvAxis) < 0 {         //[cite: 1]
		mtvAxis = mtvAxis.Scale(-1) //[cite: 1]
	}

	return CollisionResult{HasCollision: true, MTV: mtvAxis.Scale(minOverlap)} //[cite: 1]
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

type RayCollider struct {
	Origin    Vec3
	Direction Vec3    // Вектор должен быть нормализован
	Length    float64 // Длина луча (для бесконечного луча можно использовать math.MaxFloat64)
}

func NewRayCollider(origin, direction Vec3, length float64) *RayCollider {
	return &RayCollider{
		Origin:    origin,
		Direction: direction.Normalize(),
		Length:    length,
	}
}

func (r *RayCollider) Collide(other Collider) CollisionResult {
	switch o := other.(type) {
	case *BoxCollider:
		return rayBoxCollide(r, o)
	case *CapsuleCollider:
		return rayCapsuleCollide(r, o)
	case *RayCollider:
		// Пересечение двух бесконечно тонких лучей маловероятно и редко имеет смысл в физике
		return CollisionResult{HasCollision: false}
	default:
		return CollisionResult{HasCollision: false}
	}
}

func rayBoxCollide(ray *RayCollider, box *BoxCollider) CollisionResult {
	tmin := math.Inf(-1) //[cite: 1]
	tmax := math.Inf(1)  //[cite: 1]

	// Вектор от начала луча до центра коробки
	delta := box.Center.Sub(ray.Origin)

	// Полуразмеры для доступа по индексу
	he := []float64{box.HalfExtents.X, box.HalfExtents.Y, box.HalfExtents.Z}

	// Проверяем каждую из 3-х локальных осей OBB
	for i := 0; i < 3; i++ {
		axis := box.Axes[i]
		e := axis.Dot(delta)         // Проекция вектора к центру на ось OBB
		f := ray.Direction.Dot(axis) // Проекция направления луча на ось OBB

		if math.Abs(f) > 1e-9 {
			invF := 1.0 / f
			t1 := (e - he[i]) * invF
			t2 := (e + he[i]) * invF

			if t1 > t2 {
				t1, t2 = t2, t1 //[cite: 1]
			}
			if t1 > tmin {
				tmin = t1 //[cite: 1]
			}
			if t2 < tmax {
				tmax = t2 //[cite: 1]
			}
			if tmin > tmax { //[cite: 1]
				return CollisionResult{HasCollision: false} //[cite: 1]
			}
			if tmax < 0 { //[cite: 1]
				return CollisionResult{HasCollision: false} //[cite: 1]
			}
		} else {
			// Луч параллелен плоскости, проверяем, находится ли он вне коробки по этой оси[cite: 1]
			if -e-he[i] > 0 || -e+he[i] < 0 {
				return CollisionResult{HasCollision: false} //[cite: 1]
			}
		}
	}

	// Если луч имеет конечную длину, и пересечение дальше этой длины[cite: 1]
	if tmin > ray.Length { //[cite: 1]
		return CollisionResult{HasCollision: false} //[cite: 1]
	}

	// Вычисляем MTV (выталкивает луч назад по его направлению)[cite: 1]
	hitT := tmin  //[cite: 1]
	if hitT < 0 { //[cite: 1]
		hitT = 0 // Луч начинается внутри коробки[cite: 1]
	}
	penetration := tmax - hitT               //[cite: 1]
	mtv := ray.Direction.Scale(-penetration) //[cite: 1]

	return CollisionResult{HasCollision: true, MTV: mtv} //[cite: 1]
}

func rayCapsuleCollide(ray *RayCollider, cap *CapsuleCollider) CollisionResult {
	// Конечная точка луча
	rayEnd := ray.Origin.Add(ray.Direction.Scale(ray.Length))

	// Крайние точки внутреннего отрезка капсулы
	capA := cap.Center.Sub(cap.Direction.Scale(cap.HalfHeight))
	capB := cap.Center.Add(cap.Direction.Scale(cap.HalfHeight))

	// Находим ближайшие точки между отрезком луча и стержнем капсулы
	c1, c2 := closestPtSegmentSegment(ray.Origin, rayEnd, capA, capB)

	// Вектор от оси капсулы к лучу
	dir := c1.Sub(c2)
	distSqr := dir.LengthSqr()

	// Если расстояние больше радиуса капсулы, коллизии нет
	if distSqr > cap.Radius*cap.Radius {
		return CollisionResult{HasCollision: false}
	}

	dist := math.Sqrt(distSqr)
	overlap := cap.Radius - dist

	var mtv Vec3
	if dist > 1e-9 {
		// Стандартный случай: расталкиваем по нормали
		mtv = dir.Normalize().Scale(overlap)
	} else {
		// Отрезки идеально пересеклись, выталкиваем по произвольной перпендикулярной оси
		mtv = Vec3{X: 0, Y: overlap, Z: 0}
	}

	return CollisionResult{HasCollision: true, MTV: mtv}
}
