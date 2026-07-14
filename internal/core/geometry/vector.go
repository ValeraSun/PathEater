package geometry

import "math"

type Vector3 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

func GetZeroVector() Vector3 {
	return Vector3{
		X: 0,
		Y: 0,
		Z: 0,
	}
}

func (v Vector3) Length() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

func (v Vector3) Normalize() Vector3 {
	length := v.Length()
	if length == 0 {
		return Vector3{0, 0, 0}
	}
	return Vector3{
		v.X / length,
		v.Y / length,
		v.Z / length,
	}
}

func CombineVectors(direction, controlInput Vector3) Vector3 {
	if controlInput.X == 0 && controlInput.Y == 0 && controlInput.Z == 0 {
		return Vector3{}
	}

	forward := direction

	// right — перпендикуляр к forward в плоскости forward-worldUp
	worldUp := Vector3{0, 1, 0}
	right := cross(worldUp, forward).Normalize()

	// Если forward параллелен worldUp (корабль смотрит строго вверх/вниз)
	if right.X == 0 && right.Y == 0 && right.Z == 0 {
		worldUp = Vector3{0, 0, 1} // берём другую ось как up
		right = cross(worldUp, forward).Normalize()
	}

	// up — перпендикуляр к forward и right
	up := cross(forward, right)

	// Глобальный вектор = right*X + up*Y + forward*Z
	return Vector3{
		right.X*controlInput.X + up.X*controlInput.Y + forward.X*controlInput.Z,
		right.Y*controlInput.X + up.Y*controlInput.Y + forward.Y*controlInput.Z,
		right.Z*controlInput.X + up.Z*controlInput.Y + forward.Z*controlInput.Z,
	}
}

func cross(a, b Vector3) Vector3 {
	return Vector3{
		a.Y*b.Z - a.Z*b.Y,
		a.Z*b.X - a.X*b.Z,
		a.X*b.Y - a.Y*b.X,
	}
}
