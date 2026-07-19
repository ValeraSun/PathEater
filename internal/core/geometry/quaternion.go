package geometry

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
