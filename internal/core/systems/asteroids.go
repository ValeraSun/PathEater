package systems

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const (
	displaySize           = 200
	spawnZoneSize         = 100
	maxSpawnAttempts      = 50
	minAsteroidDistance   = 40
	minCosmoAlienDistance = 40
	asteroidsPerFrame     = 5
	maxAsteroids          = 30
	maxSpeed              = 20
	minSpeed              = 1
	maxRadius             = 20
	minRadius             = 1
)

type AsteroidSystem struct {
	getter     componentsGetter
	publisher  publisher
	eventQueue chan *events.MeteoriteZoneEvent
	rng        *rand.Rand
}

func NewAsteroidSystem(getter componentsGetter, publisher publisher, subscriber subscriber) *AsteroidSystem {
	s := &AsteroidSystem{
		getter:     getter,
		publisher:  publisher,
		eventQueue: make(chan *events.MeteoriteZoneEvent, 100),
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	subscriber.Subscribe("meteoriteZone", s.OnEvent)
	return s
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

	var spawnPos geometry.Vec3
	var targetPos geometry.Vec3

	halfDisplay := float64(displaySize) / 2
	halfSpawn := float64(displaySize+spawnZoneSize) / 2

	switch side {
	case 0:
		spawnPos.X = s.rng.Float64()*(halfSpawn*2) - halfSpawn
		spawnPos.Y = halfSpawn
		targetPos.X = s.rng.Float64()*(halfDisplay*2) - halfDisplay
		targetPos.Y = s.rng.Float64()*(halfDisplay*2) - halfDisplay
	case 1:
		spawnPos.X = s.rng.Float64()*(halfSpawn*2) - halfSpawn
		spawnPos.Y = -halfSpawn
		targetPos.X = s.rng.Float64()*(halfDisplay*2) - halfDisplay
		targetPos.Y = s.rng.Float64()*(halfDisplay*2) - halfDisplay
	case 2:
		spawnPos.X = -halfSpawn
		spawnPos.Y = s.rng.Float64()*(halfSpawn*2) - halfSpawn
		targetPos.X = s.rng.Float64()*(halfDisplay*2) - halfDisplay
		targetPos.Y = s.rng.Float64()*(halfDisplay*2) - halfDisplay
	case 3:
		spawnPos.X = halfSpawn
		spawnPos.Y = s.rng.Float64()*(halfSpawn*2) - halfSpawn
		targetPos.X = s.rng.Float64()*(halfDisplay*2) - halfDisplay
		targetPos.Y = s.rng.Float64()*(halfDisplay*2) - halfDisplay
	}

	direction := targetPos.Sub(spawnPos).Normalize()

	speed := minSpeed + s.rng.Float64()*(maxSpeed-minSpeed)

	radius := minRadius + s.rng.Float64()*(maxRadius-minRadius)

	return spawnPos, radius, direction, speed
}

func (s *AsteroidSystem) isPositionOccupied(pos geometry.Vec3) bool {
	aliens := s.getter.GetEntitiesByComponent("cosmoAlien")
	asteroids := s.getter.GetEntitiesByComponent("asteroid")

	for id := range aliens {
		c, exists := s.getter.GetComponent(id, "transform")
		if !exists {
			continue
		}
		transform := c.(*components.TransformComponent)

		distance := pos.Sub(transform.Position).Length()
		if distance < minCosmoAlienDistance {
			return true
		}
	}

	for id := range asteroids {
		c, exists := s.getter.GetComponent(id, "transform")
		if !exists {
			continue
		}
		transform := c.(*components.TransformComponent)

		distance := pos.Sub(transform.Position).Length()
		if distance < minAsteroidDistance {
			return true
		}
	}

	return false
}

func (s *AsteroidSystem) createAsteroid(position geometry.Vec3, radius float64, direction geometry.Vec3, speed float64) {
	e := events.NewCreateAsteroidEvent(position, radius, direction, speed)
	s.publisher.Publish(e)
}

var active bool

func (s *AsteroidSystem) Update(dt float32) error {
	if active {
		ships := s.getter.GetEntitiesByComponent("ship")
		var shipId types.Entity
		for id := range ships {
			shipId = id
		}

		c, _ := s.getter.GetComponent(shipId, "transform")
		trShip := c.(*components.TransformComponent)

		c, _ = s.getter.GetComponent(shipId, "ship")
		ship := c.(*components.ShipComponent)

		comps := s.getter.GetEntitiesByComponent("asteroid")
		for id, comp := range comps {
			aster := comp.(*components.AsteroidComponent)

			c, _ := s.getter.GetComponent(id, "transform")
			transform := c.(*components.TransformComponent)

			c, _ = s.getter.GetComponent(id, "externalVelocity")
			ext, _ := comp.(*components.ExternalVelocityComponent)

			ext.Direction = geometry.GetZeroVector().Sub(trShip.Direction.Scale(ship.Speed)).Normalize()

			x := transform.Position.X
			y := transform.Position.Y
			aster.Visible = x >= -displaySize/2 && x <= displaySize/2 && y >= -displaySize/2 && y <= displaySize/2
			aster.OnField = x >= -(displaySize+100)/2 && x <= (displaySize+100)/2 && y >= -(displaySize+100)/2 && y <= (displaySize+100)/2

			if !aster.OnField {
				e := events.NewDeleteAsteroidEvent(id)
				s.publisher.Publish(e)
			}
		}
		s.spawnAsteroids()
	}
	return nil
}

func (s *AsteroidSystem) OnEvent(event events.Event) error {
	mz, ok := event.(*events.MeteoriteZoneEvent)

	if !ok {
		return nil
	}

	select {
	case s.eventQueue <- mz:
		e := <-s.eventQueue
		active = !e.Active
	default:
		fmt.Printf("Переполена очередь %v\n", s)
	}

	return nil
}
