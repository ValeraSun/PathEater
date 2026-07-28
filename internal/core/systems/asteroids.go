package systems

import (
	"math"
	"math/rand"
	"time"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	displayHeight       = 512
	displayWidth        = 1024
	spawnZoneSize       = 100
	halfDisplayHeight   = displayHeight / 2
	halfDisplayWidth    = displayWidth / 2
	halfFieldHeight     = (displayHeight + spawnZoneSize) / 2
	halfFieldWidth      = (displayWidth + spawnZoneSize) / 2
	maxSpawnAttempts    = 50
	minAsteroidDist     = 40
	asteroidSpawnChance = 0.03
	maxAsteroids        = 50
	maxSpeed            = 50
	minSpeed            = 20
	maxRadius           = 30
	minRadius           = 5
)

type AsteroidSystem struct {
	getter    componentsGetter
	publisher publisher
	rng       *rand.Rand
	active    bool
}

func NewAsteroidSystem(getter componentsGetter, publisher publisher, subscriber subscriber) *AsteroidSystem {
	s := &AsteroidSystem{
		getter:    getter,
		publisher: publisher,
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
		active:    false,
	}
	subscriber.Subscribe("meteoriteZone", s.OnEvent)
	return s
}

func (s *AsteroidSystem) Update(dt float32) error {
	asteroids := s.getter.GetEntitiesByComponent("asteroid")
	if s.active {
		s.spawnAsteroids(asteroids)
	}
	s.deleteNotValid(asteroids)
	return nil
}

func (s *AsteroidSystem) OnEvent(event events.Event) error {
	_, ok := event.(*events.MeteoriteZoneEvent)
	if !ok {
		return nil
	}
	s.active = !s.active
	return nil
}

func (s *AsteroidSystem) spawnAsteroids(asteroids map[types.Entity]types.Component) {

	if len(asteroids) >= maxAsteroids {
		return
	}

	if s.rng.Float64() > asteroidSpawnChance {
		return
	}

	spawnCount := maxAsteroids - len(asteroids)
	for i := 0; i < spawnCount; i++ {
		s.spawnAsteroid()
	}
}

func (s *AsteroidSystem) spawnAsteroid() {
	for attempt := 0; attempt < maxSpawnAttempts; attempt++ {
		pos, rad, dir, speed := s.generateAsteroid()
		if !s.isPositionOccupied(pos) {
			s.createAsteroid(pos, rad, dir, speed)
			return
		}
	}
}

func (s *AsteroidSystem) generateAsteroid() (geometry.Vec3, float64, geometry.Vec3, float64) {
	side := s.rng.Intn(4)

	var spawnPos, targetPos geometry.Vec3

	switch side {
	case 0: //верх
		spawnPos.X = s.randomFloatInRange(-halfFieldWidth, halfFieldWidth)
		spawnPos.Y = s.randomFloatInRange(halfDisplayHeight, halfFieldHeight)
	case 1: //лево
		spawnPos.X = s.randomFloatInRange(-halfFieldWidth, -halfDisplayWidth)
		spawnPos.Y = s.randomFloatInRange(-halfFieldHeight, halfFieldHeight)
	case 2: //право
		spawnPos.X = s.randomFloatInRange(halfDisplayWidth, halfFieldWidth)
		spawnPos.Y = s.randomFloatInRange(-halfFieldHeight, halfFieldHeight)
	case 3: //низ
		spawnPos.X = s.randomFloatInRange(-halfFieldWidth, halfFieldWidth)
		spawnPos.Y = s.randomFloatInRange(-halfFieldHeight, -halfDisplayHeight)
	}

	targetPos.X = s.randomFloatInRange(-halfDisplayWidth, halfDisplayWidth)
	targetPos.Y = s.randomFloatInRange(-halfDisplayHeight, halfDisplayHeight)

	direction := targetPos.Sub(spawnPos).Normalize()
	speed := s.randomFloatInRange(minSpeed, maxSpeed)
	radius := s.randomFloatInRange(minRadius, maxRadius)

	return spawnPos, radius, direction, speed
}

func (s *AsteroidSystem) randomFloatInRange(min, max float64) float64 {
	return min + s.rng.Float64()*(max-min)
}

func (s *AsteroidSystem) isPositionOccupied(pos geometry.Vec3) bool {
	asteroids := s.getter.GetEntitiesByComponent("asteroid")
	for id := range asteroids {
		c, exists := s.getter.GetComponent(id, "transform")
		if !exists {
			continue
		}
		transform := c.(*components.TransformComponent)
		if pos.Sub(transform.Position).Length() < minAsteroidDist {
			return true
		}
	}
	return false
}

func (s *AsteroidSystem) createAsteroid(position geometry.Vec3, radius float64, direction geometry.Vec3, speed float64) {
	e := events.NewCreateAsteroidEvent(position, radius, direction, speed)
	s.publisher.Publish(e)
}

func (s *AsteroidSystem) deleteNotValid(asteroids map[types.Entity]types.Component) {

	for id := range asteroids {
		c, _ := s.getter.GetComponent(id, "transform")
		t := c.(*components.TransformComponent)

		if math.Abs(t.Position.X) > halfFieldWidth || math.Abs(t.Position.Y) > halfFieldHeight {
			e := events.NewDeleteEvent(id)
			s.publisher.Publish(e)
		}
	}

}
