package sendler

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type playerData struct {
	Position geometry.Vec3 `json:"position"`
	Rotation geometry.Vec3 `json:"rotation"`
	Health   int           `json:"health"`
}

func (s *sendler) sendPlayer(id types.Entity, broadcaster func(EntityInfo) error) error {
	c, _ := s.GetComponent(id, "transform")
	transform := c.(*components.TransformComponent)

	c, _ = s.GetComponent(id, "control")
	control := c.(*components.ControlComponent)

	c, _ = s.GetComponent(id, "health")
	hp := c.(*components.HealthComponent)

	broadcaster(EntityInfo{
		ID:   id,
		Type: "player",
		Data: playerData{
			Position: transform.Position,
			Rotation: control.Direction,
			Health:   int(hp.Health),
		},
	},
	)

	return nil
}
