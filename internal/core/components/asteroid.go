package components

import (
	"math/rand"

	"github.com/ValeraSun/PathEater/internal/core/geometry"
)

const sizeField = 200
const speedAsteroidMax = 10.5
const speedAsteroidMin = 0
const sizeAsteroidMax = 25.5
const sizeAsteroidMin = 5.5

type AsteroidComponent struct {
	Position  geometry.Vec2
	Direction geometry.Vec2
	Velocity  geometry.Vec2
}

func (*AsteroidComponent) Type() string {
	return "asteroid"
}

func NewAsteroidComponent(pos, dir, vel geometry.Vec2) *AsteroidComponent {
	return &AsteroidComponent{
		Position:  pos,
		Direction: dir,
		Velocity:  vel,
	}
}

func RandomAsteroid() (geometry.Vec2, geometry.Vec2, geometry.Vec2, float64) {
	vel := RandomVelocity()
	return RandomPosition(), vel.Normalize(), vel, RandomRadius()
}

func RandomPosition() geometry.Vec2 {
	return geometry.Vec2{
		X: rand.Float64()*(sizeField) - sizeField/2,
		Y: rand.Float64()*(sizeField) - sizeField/2,
	}
}

func RandomVelocity() geometry.Vec2 {
	return geometry.Vec2{
		X: speedAsteroidMin + rand.Float64()*(speedAsteroidMax-speedAsteroidMin),
		Y: speedAsteroidMin + rand.Float64()*(speedAsteroidMax-speedAsteroidMin),
	}
}

func RandomRadius() float64 {
	return sizeAsteroidMin + rand.Float64()*(sizeAsteroidMax-sizeAsteroidMin)
}
