package geometry

import "math"

type Vec2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func GetZeroVector2() Vec2 {
	return Vec2{
		X: 0,
		Y: 0,
	}
}

func (v Vec2) Add(u Vec2) Vec2      { return Vec2{v.X + u.X, v.Y + u.Y} }
func (v Vec2) Sub(u Vec2) Vec2      { return Vec2{v.X - u.X, v.Y - u.Y} }
func (v Vec2) Scale(s float64) Vec2 { return Vec2{v.X * s, v.Y * s} }
func (v Vec2) Dot(u Vec2) float64   { return v.X*u.X + v.Y*u.Y }
func (v Vec2) LengthSqr() float64   { return v.Dot(v) }
func (v Vec2) Length() float64      { return math.Sqrt(v.LengthSqr()) }

func (v Vec2) Normalize() Vec2 {
	l := v.Length()
	if l < 1e-9 {
		return Vec2{0, 1}
	}
	return Vec2{v.X / l, v.Y / l}
}

//Векторное умножение векторов
func (v Vec2) Cross(u Vec2) float64 {
	return v.X*u.Y - v.Y*u.X
}

func (v Vec2) Distance(other Vec2) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// DistanceSq - квадрат расстояния
func (v Vec2) DistanceSq(other Vec2) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return dx*dx + dy*dy
}

//комбинирование векторов
func CombineVectors2(direction, controlInput Vec2) Vec2 {
	if controlInput.X == 0 && controlInput.Y == 0 {
		return Vec2{}
	}

	right := Vec2{
		X: -direction.Y,
		Y: direction.X,
	}

	return Vec2{
		X: right.X*controlInput.X + direction.X*controlInput.Y,
		Y: right.Y*controlInput.X + direction.Y*controlInput.Y,
	}
}

func minVec2(v1, v2 Vec2) Vec2 {
	len1 := v1.Length()
	len2 := v2.Length()

	if len1 < len2 {
		return v1
	}
	return v2
}

func heightFromA(A, B, C Vec2) Vec2 {
	BC := C.Sub(B)

	BA := A.Sub(B)

	bcLenSq := BC.Dot(BC)

	if bcLenSq == 0 {
		panic("Сторона BC имеет нулевую длину")
	}

	t := BA.Dot(BC) / bcLenSq

	H := B.Add(BC.Scale(t))

	AH := H.Sub(A)

	return AH
}

func twoClosestPointsToTarget(target, p1, p2, p3 Vec2) (Vec2, Vec2) {
	d1 := target.DistanceSq(p1)
	d2 := target.DistanceSq(p2)
	d3 := target.DistanceSq(p3)

	points := []Vec2{p1, p2, p3}
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
