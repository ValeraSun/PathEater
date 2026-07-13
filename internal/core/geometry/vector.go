package geometry

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
