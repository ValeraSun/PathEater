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
	position geometry.Vec2
	velocity geometry.Vec2
	radius   float64
}

func (*AsteroidComponent) Type() string {
	return "asteroid"
}

func NewAsteroidComponent() *AsteroidComponent {
	return &AsteroidComponent{}
}

func (c *AsteroidComponent) Random() {
	c.RandomPosition()
	c.RandomVelocity()
	c.RandomRadius()
}

func (c *AsteroidComponent) RandomPosition() {
	c.position.X = rand.Float64()*(sizeField) - sizeField/2
	c.position.Y = rand.Float64()*(sizeField) - sizeField/2
}

func (c *AsteroidComponent) RandomVelocity() {
	c.position.X = speedAsteroidMin + rand.Float64()*(speedAsteroidMax-speedAsteroidMin)
	c.position.Y = speedAsteroidMin + rand.Float64()*(speedAsteroidMax-speedAsteroidMin)
}

func (c *AsteroidComponent) RandomRadius() {
	c.radius = sizeAsteroidMin + rand.Float64()*(sizeAsteroidMax-sizeAsteroidMin)
}
