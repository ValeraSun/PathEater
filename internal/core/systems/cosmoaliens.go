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
	cosmoDisplaySize      = 200
	cosmoSpawnZoneSize    = 100
	cosmoMaxSpawnAttempts = 50
	maxCosmoAliens        = 5
	cosmoSpawnChance      = 0.05
)

type CosmoAlienSystem struct {
	getter    componentsGetter
	publisher publisher
	rng       *rand.Rand
}

var shipID types.Entity

func NewCosmoAlienSystem(getter componentsGetter, publisher publisher) *CosmoAlienSystem {
	s := &CosmoAlienSystem{
		getter:    getter,
		publisher: publisher,
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	ships := s.getter.GetEntitiesByComponent("ship")
	for id := range ships {
		shipID = id
	}
	return s
}

func (s *CosmoAlienSystem) spawnCosmoAliens() {
	aliens := s.getter.GetEntitiesByComponent("cosmoAlien")
	if len(aliens) >= maxCosmoAliens {
		return
	}

	if s.rng.Float64() > cosmoSpawnChance {
		return
	}

	s.spawnCosmoAlien()
}

func (s *CosmoAlienSystem) spawnCosmoAlien() {
	for attempt := 0; attempt < cosmoMaxSpawnAttempts; attempt++ {
		pos := s.generateCosmoAlien()

		if !s.isPositionOccupied(pos) {
			s.createCosmoAlien(pos, shipID)
			return
		}
	}
}

func (s *CosmoAlienSystem) generateCosmoAlien() geometry.Vec3 {
	side := s.rng.Intn(4)

	var spawnPos geometry.Vec3

	halfSpawn := float64(cosmoDisplaySize+cosmoSpawnZoneSize) / 2

	switch side {
	case 0:
		spawnPos.X = s.rng.Float64()*(halfSpawn*2) - halfSpawn
		spawnPos.Y = halfSpawn
	case 1:
		spawnPos.X = s.rng.Float64()*(halfSpawn*2) - halfSpawn
		spawnPos.Y = -halfSpawn
	case 2:
		spawnPos.X = -halfSpawn
		spawnPos.Y = s.rng.Float64()*(halfSpawn*2) - halfSpawn
	case 3:
		spawnPos.X = halfSpawn
		spawnPos.Y = s.rng.Float64()*(halfSpawn*2) - halfSpawn
	}

	return spawnPos
}

func (s *CosmoAlienSystem) isPositionOccupied(pos geometry.Vec3) bool {
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
		if distance < minCosmoAlienDistance {
			return true
		}
	}

	return false
}

func (s *CosmoAlienSystem) createCosmoAlien(position geometry.Vec3, shipID types.Entity) {
	e := events.NewCreateCosmoAlienEvent(position, shipID)
	s.publisher.Publish(e)
}

func (s *CosmoAlienSystem) Update(dt float32) error {
	c, _ := s.getter.GetComponent(shipID, "transform")
	trShip := c.(*components.TransformComponent)

	c, _ = s.getter.GetComponent(shipID, "ship")
	ship := c.(*components.ShipComponent)

	comps := s.getter.GetEntitiesByComponent("cosmoAlien")
	for id, comp := range comps {
		alien := comp.(*components.CosmoAlienComponent)

		c, _ := s.getter.GetComponent(id, "transform")
		transform := c.(*components.TransformComponent)

		c, _ = s.getter.GetComponent(id, "movement")
		mov := c.(*components.MovementComponent)

		c, _ = s.getter.GetComponent(id, "externalVelocity")
		ext, _ := comp.(*components.ExternalVelocityComponent)

		ext.Direction = geometry.GetZeroVector().Sub(trShip.Direction.Scale(ship.Speed)).Normalize()

		mov.Direction = trShip.Position.Sub(transform.Position)

		x := transform.Position.X
		y := transform.Position.Y
		alien.Visible = x >= -displaySize/2 && x <= displaySize/2 && y >= -displaySize/2 && y <= displaySize/2
	}

	s.spawnCosmoAliens()

	return nil
}
