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
	cosmoSpawnZoneSize    = 100
	cosmoMaxSpawnAttempts = 50
	maxCosmoAliens        = 5
	cosmoSpawnChance      = 0.005
)

type CosmoAlienSystem struct {
	getter    componentsGetter
	publisher publisher
	rng       *rand.Rand
	shipID    types.Entity
}

func NewCosmoAlienSystem(getter componentsGetter, publisher publisher) *CosmoAlienSystem {
	s := &CosmoAlienSystem{
		getter:    getter,
		publisher: publisher,
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	s.findShipID()
	return s
}

func (s *CosmoAlienSystem) findShipID() bool {
	if s.shipID != "" {
		return true
	}

	ships := s.getter.GetEntitiesByComponent("ship")
	for id := range ships {
		s.shipID = id
		return true
	}

	return false
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
			s.createCosmoAlien(pos, s.shipID)
			return
		}
	}
}

func (s *CosmoAlienSystem) generateCosmoAlien() geometry.Vec3 {
	side := s.rng.Intn(4)

	var spawnPos geometry.Vec3

	halfSpawnWidth := float64(displayWidth+cosmoSpawnZoneSize) / 2
	halfSpawnHeight := float64(displayHeight+cosmoSpawnZoneSize) / 2

	switch side {
	case 0:
		spawnPos.X = s.rng.Float64()*(halfSpawnWidth*2) - halfSpawnWidth
		spawnPos.Y = halfSpawnHeight
	case 1:
		spawnPos.X = s.rng.Float64()*(halfSpawnWidth*2) - halfSpawnWidth
		spawnPos.Y = -halfSpawnHeight
	case 2:
		spawnPos.X = -halfSpawnWidth
		spawnPos.Y = s.rng.Float64()*(halfSpawnHeight*2) - halfSpawnHeight
	case 3:
		spawnPos.X = halfSpawnWidth
		spawnPos.Y = s.rng.Float64()*(halfSpawnHeight*2) - halfSpawnHeight
	}

	return spawnPos
}

func (s *CosmoAlienSystem) isPositionOccupied(pos geometry.Vec3) bool {
	aliens := s.getter.GetEntitiesByComponent("cosmoAlien")
	asteroids := s.getter.GetEntitiesByComponent("asteroid")

	for id := range aliens {
		c, ok := s.getter.GetComponent(id, "transform")
		if !ok {
			continue
		}
		transform, ok := c.(*components.TransformComponent)
		if !ok {
			continue
		}

		distance := pos.Sub(transform.Position).Length()
		if distance < minCosmoAlienDistance {
			return true
		}
	}
	for id := range asteroids {
		c, ok := s.getter.GetComponent(id, "transform")
		if !ok {
			continue
		}
		transform, ok := c.(*components.TransformComponent)
		if !ok {
			continue
		}

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
	if !s.findShipID() {
		return nil
	}

	c, ok := s.getter.GetComponent(s.shipID, "transform")
	if !ok {
		return nil
	}
	trShip, ok := c.(*components.TransformComponent)
	if !ok {
		return nil
	}

	c, ok = s.getter.GetComponent(s.shipID, "ship")
	if !ok {
		return nil
	}
	ship, ok := c.(*components.ShipComponent)
	if !ok {
		return nil
	}

	comps := s.getter.GetEntitiesByComponent("cosmoAlien")
	for id, comp := range comps {
		alien, ok := comp.(*components.CosmoAlienComponent)
		if !ok {
			continue
		}

		c, ok := s.getter.GetComponent(id, "transform")
		if !ok {
			continue
		}
		transform, ok := c.(*components.TransformComponent)
		if !ok {
			continue
		}

		c, ok = s.getter.GetComponent(id, "movement")
		if !ok {
			continue
		}
		mov, ok := c.(*components.MovementComponent)
		if !ok {
			continue
		}

		c, ok = s.getter.GetComponent(id, "externalVelocity")
		if !ok {
			continue
		}
		ext, ok := c.(*components.ExternalVelocityComponent)
		if !ok {
			continue
		}

		ext.Direction = geometry.GetZeroVector().Sub(trShip.Direction.Scale(ship.Speed)).Normalize()

		mov.Direction = trShip.Position.Sub(transform.Position)

		x := transform.Position.X
		y := transform.Position.Y
		alien.Visible = x >= -displayWidth/2 && x <= displayWidth/2 && y >= -displayHeight/2 && y <= displayHeight/2
	}

	s.spawnCosmoAliens()

	return nil
}
