package geometry

import (
	"math"
)

type CollisionResult struct {
	HasCollision bool
	MTV          Vec3 // Вектор, на который нужно сдвинуть объект, у которого вызван метод, чтобы выйти из коллизии[cite: 1]
}

type Collider interface {
	Collide(other Collider) CollisionResult //[cite: 1]
	ChangeCenter(center Vec3)               //[cite: 1]
}

type BoxCollider struct {
	Center      Vec3       //[cite: 1]
	HalfExtents Vec3       //[cite: 1]
	Rotation    Quaternion // Заменили Axes [3]Vec3 на Quaternion[cite: 1, 2]
}

func NewBoxCollider(center, halfExtents Vec3, rotation Quaternion) *BoxCollider {
	return &BoxCollider{
		Center:      center,
		HalfExtents: halfExtents,
		Rotation:    rotation,
	}
}

func (b *BoxCollider) ChangeCenter(center Vec3) {
	b.Center = center //[cite: 1]
}

func (b *BoxCollider) Collide(other Collider) CollisionResult {
	switch o := other.(type) {
	case *BoxCollider:
		return boxBoxCollide(b, o) //[cite: 1]
	case *CapsuleCollider:
		return boxCapsuleCollide(b, o) //[cite: 1]
	case *RayCollider:
		// Инвертируем MTV, чтобы выталкивался Box, а не Ray[cite: 1]
		res := rayBoxCollide(o, b) //[cite: 1]
		if res.HasCollision {
			res.MTV = res.MTV.Scale(-1) //[cite: 1]
		}
		return res //[cite: 1]
	default:
		return CollisionResult{HasCollision: false} //[cite: 1]
	}
}

type CapsuleCollider struct {
	Center     Vec3    //[cite: 1]
	Direction  Vec3    //[cite: 1]
	HalfHeight float64 //[cite: 1]
	Radius     float64 //[cite: 1]
}

func NewCapsuleCollider(center, direction Vec3, halfHeight, radius float64) *CapsuleCollider {
	dir := direction.Normalize() //[cite: 1]
	return &CapsuleCollider{
		Center:     center,
		Direction:  dir,
		HalfHeight: halfHeight,
		Radius:     radius,
	}
}

func (b *CapsuleCollider) ChangeCenter(center Vec3) {
	b.Center = center //[cite: 1]
}

func (c *CapsuleCollider) Collide(other Collider) CollisionResult {
	switch o := other.(type) {
	case *BoxCollider:
		res := boxCapsuleCollide(o, c) //[cite: 1]
		if res.HasCollision {
			res.MTV = res.MTV.Scale(-1) //[cite: 1]
		}
		return res //[cite: 1]
	case *CapsuleCollider:
		return capsuleCapsuleCollide(c, o) //[cite: 1]
	case *RayCollider:
		// Инвертируем MTV[cite: 1]
		res := rayCapsuleCollide(o, c) //[cite: 1]
		if res.HasCollision {
			res.MTV = res.MTV.Scale(-1) //[cite: 1]
		}
		return res //[cite: 1]
	default:
		return CollisionResult{HasCollision: false} //[cite: 1]
	}
}

func boxBoxCollide(a, b *BoxCollider) CollisionResult {
	// Извлекаем матрицы поворота (локальные оси) из кватернионов[cite: 2]
	axesA := a.Rotation.ToRotationMatrix()
	axesB := b.Rotation.ToRotationMatrix()

	// Подготавливаем 15 потенциальных разделяющих осей[cite: 1]
	axes := make([]Vec3, 0, 15) //[cite: 1]
	axes = append(axes, axesA[0], axesA[1], axesA[2])
	axes = append(axes, axesB[0], axesB[1], axesB[2])

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			cross := axesA[i].Cross(axesB[j]) //[cite: 1]
			// Избегаем нулевых векторов при параллельных осях[cite: 1]
			if cross.LengthSqr() > 1e-9 { //[cite: 1]
				axes = append(axes, cross.Normalize()) //[cite: 1]
			}
		}
	}

	minOverlap := math.MaxFloat64 //[cite: 1]
	var mtvAxis Vec3              //[cite: 1]
	d := a.Center.Sub(b.Center)   // Вектор от центра B к центру A[cite: 1]

	for _, axis := range axes {
		if axis.LengthSqr() < 1e-9 { //[cite: 1]
			continue //[cite: 1]
		}
		axis = axis.Normalize() //[cite: 1]

		// Проецируем полуразмеры (HalfExtents) коробки A на ось[cite: 1]
		rA := a.HalfExtents.X*math.Abs(axesA[0].Dot(axis)) +
			a.HalfExtents.Y*math.Abs(axesA[1].Dot(axis)) +
			a.HalfExtents.Z*math.Abs(axesA[2].Dot(axis)) //[cite: 1]

		// Проецируем полуразмеры (HalfExtents) коробки B на ось[cite: 1]
		rB := b.HalfExtents.X*math.Abs(axesB[0].Dot(axis)) +
			b.HalfExtents.Y*math.Abs(axesB[1].Dot(axis)) +
			b.HalfExtents.Z*math.Abs(axesB[2].Dot(axis)) //[cite: 1]

		// Расстояние между центрами, спроецированное на ось[cite: 1]
		dist := math.Abs(d.Dot(axis)) //[cite: 1]

		// Если расстояние между центрами больше суммы спроецированных радиусов, то коллизии нет[cite: 1]
		overlap := rA + rB - dist //[cite: 1]
		if overlap <= 0 {
			return CollisionResult{HasCollision: false} //[cite: 1]
		}

		// Ищем минимальное перекрытие для выбора оси выталкивания[cite: 1]
		if overlap < minOverlap {
			minOverlap = overlap //[cite: 1]
			mtvAxis = axis       //[cite: 1]
		}
	}

	// Направление вектора MTV должно выталкивать A из B[cite: 1]
	if d.Dot(mtvAxis) < 0 {
		mtvAxis = mtvAxis.Scale(-1) //[cite: 1]
	}

	return CollisionResult{HasCollision: true, MTV: mtvAxis.Scale(minOverlap)} //[cite: 1]
}

func capsuleCapsuleCollide(a, b *CapsuleCollider) CollisionResult {
	a1 := a.Center.Sub(a.Direction.Scale(a.HalfHeight)) //[cite: 1]
	a2 := a.Center.Add(a.Direction.Scale(a.HalfHeight)) //[cite: 1]
	b1 := b.Center.Sub(b.Direction.Scale(b.HalfHeight)) //[cite: 1]
	b2 := b.Center.Add(b.Direction.Scale(b.HalfHeight)) //[cite: 1]

	c1, c2 := closestPtSegmentSegment(a1, a2, b1, b2) //[cite: 1]
	dir := c1.Sub(c2)                                 //[cite: 1]
	distSqr := dir.LengthSqr()                        //[cite: 1]
	radSum := a.Radius + b.Radius                     //[cite: 1]

	if distSqr > radSum*radSum { //[cite: 1]
		return CollisionResult{HasCollision: false} //[cite: 1]
	}

	dist := math.Sqrt(distSqr) //[cite: 1]
	overlap := radSum - dist   //[cite: 1]

	var mtv Vec3
	if dist > 1e-9 {
		// Стандартный случай: расталкиваем по вектору между ближайшими точками[cite: 1]
		mtv = dir.Normalize().Scale(overlap) //[cite: 1]
	} else {
		// Центры точно совпали, расталкиваем по произвольной оси (например, Y)[cite: 1]
		mtv = Vec3{X: 0, Y: overlap, Z: 0} //[cite: 1]
	}

	return CollisionResult{HasCollision: true, MTV: mtv} //[cite: 1]
}

func boxCapsuleCollide(box *BoxCollider, cap *CapsuleCollider) CollisionResult {
	// Конвертируем кватернион в матрицу 3x3 для получения локальных осей[cite: 2]
	boxAxes := box.Rotation.ToRotationMatrix()

	a := cap.Center.Sub(cap.Direction.Scale(cap.HalfHeight)) //[cite: 1]
	b := cap.Center.Add(cap.Direction.Scale(cap.HalfHeight)) //[cite: 1]

	// Тестируем локальные оси OBB, ось капсулы и их попарные векторные произведения[cite: 1]
	axes := []Vec3{
		boxAxes[0],
		boxAxes[1],
		boxAxes[2],
		cap.Direction, //[cite: 1]
		cap.Direction.Cross(boxAxes[0]),
		cap.Direction.Cross(boxAxes[1]),
		cap.Direction.Cross(boxAxes[2]),
	}

	minOverlap := math.MaxFloat64 //[cite: 1]
	var mtvAxis Vec3              //[cite: 1]

	for _, axis := range axes {
		if axis.LengthSqr() < 1e-9 { //[cite: 1]
			continue //[cite: 1]
		}
		axis = axis.Normalize() //[cite: 1]

		// Проецируем OBB на ось[cite: 1]
		boxProj := box.HalfExtents.X*math.Abs(boxAxes[0].Dot(axis)) +
			box.HalfExtents.Y*math.Abs(boxAxes[1].Dot(axis)) +
			box.HalfExtents.Z*math.Abs(boxAxes[2].Dot(axis)) //[cite: 1]

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

func closestPtSegmentSegment(p1, q1, p2, q2 Vec3) (Vec3, Vec3) {
	d1 := q1.Sub(p1) //[cite: 1]
	d2 := q2.Sub(p2) //[cite: 1]
	r := p1.Sub(p2)  //[cite: 1]
	a := d1.Dot(d1)  //[cite: 1]
	e := d2.Dot(d2)  //[cite: 1]
	f := d2.Dot(r)   //[cite: 1]

	s, t := 0.0, 0.0 //[cite: 1]

	if a <= 1e-9 && e <= 1e-9 { //[cite: 1]
		return p1, p2 //[cite: 1]
	}
	if a <= 1e-9 { //[cite: 1]
		t = clamp(f/e, 0.0, 1.0) //[cite: 1]
	} else {
		c := d1.Dot(r) //[cite: 1]
		if e <= 1e-9 { //[cite: 1]
			s = clamp(-c/a, 0.0, 1.0) //[cite: 1]
			t = 0.0                   //[cite: 1]
		} else {
			b := d1.Dot(d2)    //[cite: 1]
			denom := a*e - b*b //[cite: 1]

			if denom != 0.0 { //[cite: 1]
				s = clamp((b*f-c*e)/denom, 0.0, 1.0) //[cite: 1]
			} else {
				s = 0.0 //[cite: 1]
			}

			t = (b*s + f) / e //[cite: 1]
			if t < 0.0 {      //[cite: 1]
				t = 0.0                   //[cite: 1]
				s = clamp(-c/a, 0.0, 1.0) //[cite: 1]
			} else if t > 1.0 { //[cite: 1]
				t = 1.0                      //[cite: 1]
				s = clamp((b-c)/a, 0.0, 1.0) //[cite: 1]
			}
		}
	}
	c1 := p1.Add(d1.Scale(s)) //[cite: 1]
	c2 := p2.Add(d2.Scale(t)) //[cite: 1]
	return c1, c2             //[cite: 1]
}

type RayCollider struct {
	Origin    Vec3    //[cite: 1]
	Direction Vec3    // Вектор должен быть нормализован[cite: 1]
	Length    float64 // Длина луча (для бесконечного луча можно использовать math.MaxFloat64)[cite: 1]
}

func NewRayCollider(origin, direction Vec3, length float64) *RayCollider {
	return &RayCollider{
		Origin:    origin,
		Direction: direction.Normalize(), //[cite: 1]
		Length:    length,
	}
}

func (r *RayCollider) ChangeCenter(center Vec3) {
	r.Origin = center //[cite: 1]
}

func (r *RayCollider) Collide(other Collider) CollisionResult {
	switch o := other.(type) {
	case *BoxCollider:
		return rayBoxCollide(r, o) //[cite: 1]
	case *CapsuleCollider:
		return rayCapsuleCollide(r, o) //[cite: 1]
	case *RayCollider:
		// Пересечение двух бесконечно тонких лучей маловероятно и редко имеет смысл в физике[cite: 1]
		return CollisionResult{HasCollision: false} //[cite: 1]
	default:
		return CollisionResult{HasCollision: false} //[cite: 1]
	}
}

func rayBoxCollide(ray *RayCollider, box *BoxCollider) CollisionResult {
	tmin := math.Inf(-1) //[cite: 1]
	tmax := math.Inf(1)  //[cite: 1]

	// Получаем базис из кватерниона[cite: 2]
	boxAxes := box.Rotation.ToRotationMatrix()

	// Вектор от начала луча до центра коробки[cite: 1]
	delta := box.Center.Sub(ray.Origin) //[cite: 1]

	// Полуразмеры для доступа по индексу[cite: 1]
	he := []float64{box.HalfExtents.X, box.HalfExtents.Y, box.HalfExtents.Z} //[cite: 1]

	// Проверяем каждую из 3-х локальных осей OBB[cite: 1]
	for i := 0; i < 3; i++ {
		axis := boxAxes[i]           // Используем извлеченные оси
		e := axis.Dot(delta)         // Проекция вектора к центру на ось OBB[cite: 1]
		f := ray.Direction.Dot(axis) // Проекция направления луча на ось OBB[cite: 1]

		if math.Abs(f) > 1e-9 { //[cite: 1]
			invF := 1.0 / f          //[cite: 1]
			t1 := (e - he[i]) * invF //[cite: 1]
			t2 := (e + he[i]) * invF //[cite: 1]

			if t1 > t2 { //[cite: 1]
				t1, t2 = t2, t1 //[cite: 1]
			}
			if t1 > tmin { //[cite: 1]
				tmin = t1 //[cite: 1]
			}
			if t2 < tmax { //[cite: 1]
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
			if -e-he[i] > 0 || -e+he[i] < 0 { //[cite: 1]
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
	// Конечная точка луча[cite: 1]
	rayEnd := ray.Origin.Add(ray.Direction.Scale(ray.Length)) //[cite: 1]

	// Крайние точки внутреннего отрезка капсулы[cite: 1]
	capA := cap.Center.Sub(cap.Direction.Scale(cap.HalfHeight)) //[cite: 1]
	capB := cap.Center.Add(cap.Direction.Scale(cap.HalfHeight)) //[cite: 1]

	// Находим ближайшие точки между отрезком луча и стержнем капсулы[cite: 1]
	c1, c2 := closestPtSegmentSegment(ray.Origin, rayEnd, capA, capB) //[cite: 1]

	// Вектор от оси капсулы к лучу[cite: 1]
	dir := c1.Sub(c2)          //[cite: 1]
	distSqr := dir.LengthSqr() //[cite: 1]

	// Если расстояние больше радиуса капсулы, коллизии нет[cite: 1]
	if distSqr > cap.Radius*cap.Radius { //[cite: 1]
		return CollisionResult{HasCollision: false} //[cite: 1]
	}

	dist := math.Sqrt(distSqr)   //[cite: 1]
	overlap := cap.Radius - dist //[cite: 1]

	var mtv Vec3
	if dist > 1e-9 {
		// Стандартный случай: расталкиваем по нормали[cite: 1]
		mtv = dir.Normalize().Scale(overlap) //[cite: 1]
	} else {
		// Отрезки идеально пересеклись, выталкиваем по произвольной перпендикулярной оси[cite: 1]
		mtv = Vec3{X: 0, Y: overlap, Z: 0} //[cite: 1]
	}

	return CollisionResult{HasCollision: true, MTV: mtv} //[cite: 1]
}
