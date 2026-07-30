package systems

import (
	"math/rand"
	"time"

	"github.com/ValeraSun/PathEater/internal/core/events"
)

type ZonesSystem struct {
	publisher publisher
	rng       *rand.Rand
}

const meteoriteZoneSpawnChance = 0.013

func NewZonesSystem(publisher publisher) *ZonesSystem {
	return &ZonesSystem{
		publisher: publisher,
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *ZonesSystem) Update(dt float32) error {
	if s.rng.Float64() > meteoriteZoneSpawnChance {
		return nil
	}
	s.publisher.Publish(events.NewMeteoriteZoneEvent())
	return nil
}
