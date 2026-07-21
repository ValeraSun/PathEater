package geometry

import "math"

//Векторное умножение векторов
func (v Vec3) Cross2(u Vec3) float64 {
	return v.X*u.Y - v.Y*u.X
}

//возвращает перпендикулярный по часовой стрелке и равный по модулю вектор
func (v Vec3) Perpend() Vec3 {
	return Vec3{
		X: v.Y,
		Y: v.X,
		Z: 0,
	}
}

func (v Vec3) Distance(other Vec3) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (v Vec3) DistanceSq(other Vec3) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return dx*dx + dy*dy
}

func (v1 Vec3) CosOfAngleBetweenVec2(v2 Vec3) float64 {
	return v1.Dot(v2) / (v1.Length() * v2.Length())
}

func (v Vec3) RotateBySinCos(sin, cos float64) {
	v.X = v.X*cos - v.Y*sin
	v.Y = v.Y*cos + v.X*sin
}

func (v Vec3) Rotate(angle float64) {
	v.RotateBySinCos(math.Sincos(angle))
}

func (v Vec3) RotateAroundPoint(pos *Vec3, center Vec3, cos float64) {
	sin := cosToSin(cos)
	v.RotateBySinCos(sin, cos)
	dxA := pos.X - center.X
	dyA := pos.Y - center.Y
	pos.X = center.X + dxA*cos - dyA*sin
	pos.Y = center.Y + dxA*sin + dyA*cos
}

func heightFromA(A, B, C Vec3) Vec3 {
	BC := C.Sub(B)

	BA := A.Sub(B)

	bcLenSq := BC.Dot(BC)

	t := BA.Dot(BC) / bcLenSq

	H := B.Add(BC.Scale(t))

	AH := H.Sub(A)

	return AH
}

func twoClosestPointsToTarget(target, p1, p2, p3 Vec3) (Vec3, Vec3) {
	d1 := target.DistanceSq(p1)
	d2 := target.DistanceSq(p2)
	d3 := target.DistanceSq(p3)

	points := []Vec3{p1, p2, p3}
	distances := []float64{d1, d2, d3}

	minIdx := 0
	for i := 1; i < 3; i++ {
		if distances[i] < distances[minIdx] {
			minIdx = i
		}
	}

	secondMinIdx := -1
	secondMinDist := math.Inf(1)

	for i := 0; i < 3; i++ {
		if i == minIdx {
			continue
		}
		if distances[i] < secondMinDist {
			secondMinIdx = i
			secondMinDist = distances[i]
		}
	}

	return points[minIdx], points[secondMinIdx]
}
