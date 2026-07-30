package geometry

import "math"

type Quaternion struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	W float64 `json:"w"`
}

func (q Quaternion) ToRotationMatrix() [3]Vec3 {
	x, y, z, w := q.X, q.Y, q.Z, q.W

	xx := x * x
	yy := y * y
	zz := z * z
	xy := x * y
	xz := x * z
	yz := y * z
	wx := w * x
	wy := w * y
	wz := w * z

	return [3]Vec3{
		{1 - 2*(yy+zz), 2 * (xy + wz), 2 * (xz - wy)},
		{2 * (xy - wz), 1 - 2*(xx+zz), 2 * (yz + wx)},
		{2 * (xz + wy), 2 * (yz - wx), 1 - 2*(xx+yy)},
	}
}

func MatrixToEuler(axes [3]Vec3) (x, y, z float64) {
	m00 := axes[0].X
	m10, m11, m12 := axes[1].X, axes[1].Y, axes[1].Z // вторая
	m20, m21, m22 := axes[2].X, axes[2].Y, axes[2].Z // третья

	// Проверка на сингулярность
	sy := math.Sqrt(m00*m00 + m10*m10)
	if sy > 1e-6 {
		x = math.Atan2(m21, m22)
		y = math.Atan2(-m20, sy)
		z = math.Atan2(m10, m00)
	} else {
		x = math.Atan2(-m12, m11)
		y = math.Atan2(-m20, sy)
		z = 0
	}
	return
}
