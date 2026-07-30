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
	cosmoSpawnChance      = 0.012
	minCosmoAlienDistance = 40
)

type CosmoAlienSystem struct {
	componentsGetter
	publisher
	rng    *rand.Rand
	shipID types.Entity
}

func NewCosmoAlienSystem(getter componentsGetter, publisher publisher) *CosmoAlienSystem {
	s := &CosmoAlienSystem{
		getter,
		publisher,
		rand.New(rand.NewSource(time.Now().UnixNano())),
		"",
	}
	s.findShipID()
	return s
}

func (s *CosmoAlienSystem) findShipID() bool {
	if s.shipID != "" {
		return true
	}
	s.shipID, _ = getShip(s)

	if s.shipID == "" {
		return false
	}

	return true
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
	aliens := s.GetEntitiesByComponent("cosmoAlien")
	for id := range aliens {
		c, ok := s.GetComponent(id, "transform")
		if !ok {
			continue
		}
		transform := c.(*components.TransformComponent)
		if pos.Sub(transform.Position).Length() < minCosmoAlienDistance {
			return true
		}
	}
	asteroids := s.GetEntitiesByComponent("asteroid")
	for id := range asteroids {
		c, ok := s.GetComponent(id, "transform")
		if !ok {
			continue
		}
		transform := c.(*components.TransformComponent)
		if pos.Sub(transform.Position).Length() < minCosmoAlienDistance {
			return true
		}
	}
	return false
}

func (s *CosmoAlienSystem) Update(dt float32) error {
	if !s.findShipID() {
		return nil
	}

	aliens := s.GetEntitiesByComponent("cosmoAlien")
	for id := range aliens {
		transformRaw, _ := s.GetComponent(id, "transform")
		transform := transformRaw.(*components.TransformComponent)

		moveRaw, _ := s.GetComponent(id, "movement")
		move := moveRaw.(*components.MovementComponent)

		toCenter := geometry.Vec3{X: 0, Y: 0, Z: 0}.Sub(transform.Position)
		if toCenter.Length() > 1 {
			move.Direction = toCenter.Normalize()
		}
	}

	s.spawnCosmoAliens()
	return nil
}

func (s *CosmoAlienSystem) spawnCosmoAliens() {
	aliens := s.GetEntitiesByComponent("cosmoAlien")
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

func (s *CosmoAlienSystem) createCosmoAlien(position geometry.Vec3, shipID types.Entity) {
	e := events.NewCreateCosmoAlienEvent(position, shipID)
	s.publisher.Publish(e)
}
