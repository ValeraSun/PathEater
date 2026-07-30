package sendler

import (
	"errors"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type cosmoAlienData struct {
	Position geometry.Vec3 `json:"position"`
	Rotation geometry.Vec3 `json:"rotation"`
	Health   int           `json:"health"`
	Died     bool          `json:"died"`
}

func (s *sendler) sendCosmoAlien(id types.Entity, broadcaster func(EntityInfo) error) error {
	c, ok := s.GetComponent(id, "cosmoAlien")
	if !ok {
		return errors.New("не найден компонент cosmoAlien")
	}
	alien := c.(*components.CosmoAlienComponent)

	c, ok = s.GetComponent(id, "transform")
	if !ok {
		return errors.New("не найден компонент transform")
	}
	transform := c.(*components.TransformComponent)

	c, ok = s.GetComponent(id, "health")
	if !ok {
		return errors.New("не найден компонент health")
	}
	hp := c.(*components.HealthComponent)

	broadcaster(EntityInfo{
		ID:   id,
		Type: "cosmoAlien",
		Data: cosmoAlienData{
			Position: transform.Position,
			Rotation: transform.Direction,
			Health:   int(hp.Health),
			Died:     alien.Died,
		},
	})

	return nil
}
