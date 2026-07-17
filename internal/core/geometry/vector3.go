package geometry

import "math"

type Vec3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func GetZeroVector() Vec3 {
	return Vec3{
		X: 0,
		Y: 0,
		Z: 0,
	}
}

func (v Vec3) Add(u Vec3) Vec3      { return Vec3{v.X + u.X, v.Y + u.Y, v.Z + u.Z} }
func (v Vec3) Sub(u Vec3) Vec3      { return Vec3{v.X - u.X, v.Y - u.Y, v.Z - u.Z} }
func (v Vec3) Scale(s float64) Vec3 { return Vec3{v.X * s, v.Y * s, v.Z * s} }
func (v Vec3) Dot(u Vec3) float64   { return v.X*u.X + v.Y*u.Y + v.Z*u.Z }
func (v Vec3) LengthSqr() float64   { return v.Dot(v) }
func (v Vec3) Length() float64      { return math.Sqrt(v.LengthSqr()) }

func (v Vec3) Normalize() Vec3 {
	l := v.Length()
	if l < 1e-9 {
		return Vec3{0, 0, 0}
	}
	return Vec3{v.X / l, v.Y / l, v.Z / l}
}

func (v Vec3) Cross(u Vec3) Vec3 {
	return Vec3{v.Y*u.Z - v.Z*u.Y, v.Z*u.X - v.X*u.Z, v.X*u.Y - v.Y*u.X}
}

<<<<<<< HEAD:internal/core/geometry/vector3.go
=======
// func abs(f float64) float64 {
// 	if f < 0 {
// 		return -f
// 	}
// 	return f
// }

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

>>>>>>> connection:internal/core/geometry/vector.go
func CombineVectors(direction, controlInput Vec3) Vec3 {
	if controlInput.X == 0 && controlInput.Y == 0 && controlInput.Z == 0 {
		return Vec3{}
	}

	forward := direction

	// right — перпендикуляр к forward в плоскости forward-worldUp
	worldUp := Vec3{0, 1, 0}
	right := worldUp.Cross(forward).Normalize()

	// Если forward параллелен worldUp (корабль смотрит строго вверх/вниз)
	if right.X == 0 && right.Y == 0 && right.Z == 0 {
		worldUp = Vec3{0, 0, 1} // берём другую ось как up
		right = worldUp.Cross(forward).Normalize()
	}

	// up — перпендикуляр к forward и right
	up := forward.Cross(right)

	// Глобальный вектор = right*X + up*Y + forward*Z
	return Vec3{
		right.X*controlInput.X + up.X*controlInput.Y + forward.X*controlInput.Z,
		right.Y*controlInput.X + up.Y*controlInput.Y + forward.Y*controlInput.Z,
		right.Z*controlInput.X + up.Z*controlInput.Y + forward.Z*controlInput.Z,
	}
}
