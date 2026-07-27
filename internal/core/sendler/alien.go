package sendler

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type alienData struct {
	Position geometry.Vec3 `json:"position"`
	Rotation geometry.Vec3 `json:"rotation"`
	Health   int           `json:"health"`
	Dead     bool          `json:"dead"`
}

func (s *sendler) sendAlien(id types.Entity, broadcaster func(EntityInfo) error) error {

	c, _ := s.GetComponent(id, "transform")
	transform := c.(*components.TransformComponent)

	c, _ = s.GetComponent(id, "ai")
	ai := c.(*components.AIComponent)

	c, _ = s.GetComponent(id, "health")
	hp := c.(*components.HealthComponent)

	broadcaster(EntityInfo{
		ID:   id,
		Type: "alien",
		Data: alienData{
			Position: transform.Position,
			Rotation: ai.Direction,
			Health:   hp.Health,
			Dead:     hp.Health == 0,
		},
	},
	)
	return nil
}
