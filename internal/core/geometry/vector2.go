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

func (v1 Vec3) CosOfAngleBetweenVec2(v2 Vec3) float64 {
	return v1.Dot(v2) / (v1.Length() * v2.Length())
}

func (v *Vec3) RotateBySinCos(sin, cos float64) {
	x := v.X
	y := v.Y
	v.X = x*cos - y*sin
	v.Y = y*cos + x*sin
}

func (v Vec3) Rotate(angle float64) {
	sin, cos := math.Sincos(angle)
	v.RotateBySinCos(sin, cos)
}

func (v *Vec3) RotateAroundPoint(pos *Vec3, center Vec3, cos float64) {
	sin := cosToSin(cos)
	v.RotateBySinCos(sin, cos)
	dxA := pos.X - center.X
	dyA := pos.Y - center.Y
	pos.X = center.X + dxA*cos - dyA*sin
	pos.Y = center.Y + dxA*sin + dyA*cos
}

//расстояние от точки до отрезка
func pointSegmentDistance(p, a, b Vec3) (float64, Vec3) {
	ab := b.Sub(a)
	ap := p.Sub(a)

	t := (ap.X*ab.X + ap.Y*ab.Y + ap.Z*ab.Z) / (ab.X*ab.X + ab.Y*ab.Y + ab.Z*ab.Z)
	t = math.Max(0, math.Min(1, t))

	closest := a.Add(ab.Scale(t))
	diff := p.Sub(closest)
	return diff.Length(), closest
}

//расстояние от точки до треугольника
func pointTriangleDistance(p, v1, v2, v3 Vec3) (float64, Vec3, Vec3) {
	var minDist float64 = math.MaxFloat64
	var closestPoint Vec3
	var edgeNormal Vec3

	dist, cp := pointSegmentDistance(p, v1, v2)
	if dist < minDist {
		minDist = dist
		closestPoint = cp
		edge := v2.Sub(v1)
		edgeNormal = Vec3{X: -edge.Y, Y: edge.X, Z: 0}
		if edgeNormal.Length() > 0 {
			edgeNormal = edgeNormal.Normalize()
		}
	}

	dist, cp = pointSegmentDistance(p, v2, v3)
	if dist < minDist {
		minDist = dist
		closestPoint = cp
		edge := v3.Sub(v2)
		edgeNormal = Vec3{X: -edge.Y, Y: edge.X, Z: 0}
		if edgeNormal.Length() > 0 {
			edgeNormal = edgeNormal.Normalize()
		}
	}

	dist, cp = pointSegmentDistance(p, v3, v1)
	if dist < minDist {
		minDist = dist
		closestPoint = cp
		edge := v1.Sub(v3)
		edgeNormal = Vec3{X: -edge.Y, Y: edge.X, Z: 0}
		if edgeNormal.Length() > 0 {
			edgeNormal = edgeNormal.Normalize()
		}
	}

	if pointInTriangle(p, v1, v2, v3) {
		center := Vec3{
			X: (v1.X + v2.X + v3.X) / 3,
			Y: (v1.Y + v2.Y + v3.Y) / 3,
			Z: 0,
		}
		dir := p.Sub(center)
		if dir.Length() > 0 {
			edgeNormal = dir.Normalize()
		} else {
			edgeNormal = Vec3{X: 0, Y: 1, Z: 0}
		}
		return 0, p, edgeNormal
	}

	return minDist, closestPoint, edgeNormal
}

//проверяет, внутри треугольника ли точка
func pointInTriangle(p, v1, v2, v3 Vec3) bool {
	d1 := sign(p, v1, v2)
	d2 := sign(p, v2, v3)
	d3 := sign(p, v3, v1)

	hasNeg := (d1 < 0) || (d2 < 0) || (d3 < 0)
	hasPos := (d1 > 0) || (d2 > 0) || (d3 > 0)

	return !(hasNeg && hasPos)
}

//проверяет, с какой стороны от линии v1-v2 точка p. Отрицательное значение - слева, положительное - справа
func sign(p, v1, v2 Vec3) float64 {
	return (p.X-v2.X)*(v1.Y-v2.Y) - (v1.X-v2.X)*(p.Y-v2.Y)
}
