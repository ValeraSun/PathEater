package systems

import (
	"math/rand"
	"time"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	displayHeight    = 512
	displayWidth     = 1024
	spawnZoneSize    = 100
	maxSpawnAttempts = 50
	minAsteroidDist  = 40
	asteroidsPerFrame = 5
	maxAsteroids     = 30
	maxSpeed         = 20
	minSpeed         = 1
	maxRadius        = 50
	minRadius        = 5
)

type AsteroidSystem struct {
	getter    componentsGetter
	publisher publisher
	subscriber subscriber // добавлен для подписки
	rng       *rand.Rand
	active    bool
}

func NewAsteroidSystem(getter componentsGetter, publisher publisher, subscriber subscriber) *AsteroidSystem {
	s := &AsteroidSystem{
		getter:     getter,
		publisher:  publisher,
		subscriber: subscriber,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
		active:     false, 
	}
	subscriber.Subscribe("meteoriteZone", s.OnEvent)
	return s
}

func (s *AsteroidSystem) Update(dt float32) error {
	if s.active {
		s.spawnAsteroids()
	}
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

func (s *AsteroidSystem) spawnAsteroids() {
	asteroids := s.getter.GetEntitiesByComponent("asteroid")
	if len(asteroids) >= maxAsteroids {
		return
	}

	spawnCount := min(asteroidsPerFrame, maxAsteroids-len(asteroids))
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

	halfDisplayHeight := float64(displayHeight) / 2
	halfSpawnHeight := float64(displayHeight+spawnZoneSize) / 2
	halfDisplayWidth := float64(displayWidth) / 2
	halfSpawnWidth := float64(displayWidth+spawnZoneSize) / 2

	switch side {
	case 0: 
		spawnPos.X = s.rng.Float64()*(halfSpawnWidth*2) - halfSpawnWidth
		spawnPos.Y = halfSpawnHeight
		targetPos.X = s.rng.Float64()*(halfDisplayWidth*2) - halfDisplayWidth
		targetPos.Y = s.rng.Float64()*(halfDisplayHeight*2) - halfDisplayHeight
	case 1: 
		spawnPos.X = s.rng.Float64()*(halfSpawnWidth*2) - halfSpawnWidth
		spawnPos.Y = -halfSpawnHeight
		targetPos.X = s.rng.Float64()*(halfDisplayWidth*2) - halfDisplayWidth
		targetPos.Y = s.rng.Float64()*(halfDisplayHeight*2) - halfDisplayHeight
	case 2: 
		spawnPos.X = -halfSpawnHeight
		spawnPos.Y = s.rng.Float64()*(halfSpawnWidth*2) - halfSpawnWidth
		targetPos.X = s.rng.Float64()*(halfDisplayWidth*2) - halfDisplayWidth
		targetPos.Y = s.rng.Float64()*(halfDisplayHeight*2) - halfDisplayHeight
	case 3: 
		spawnPos.X = halfSpawnHeight
		spawnPos.Y = s.rng.Float64()*(halfSpawnWidth*2) - halfSpawnWidth
		targetPos.X = s.rng.Float64()*(halfDisplayWidth*2) - halfDisplayWidth
		targetPos.Y = s.rng.Float64()*(halfDisplayHeight*2) - halfDisplayHeight
	}

	direction := targetPos.Sub(spawnPos).Normalize()
	speed := minSpeed + s.rng.Float64()*(maxSpeed-minSpeed)
	radius := minRadius + s.rng.Float64()*(maxRadius-minRadius)

	return spawnPos, radius, direction, speed
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