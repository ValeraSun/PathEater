package sendler

import (
	"errors"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type breakdownData struct {
	Position geometry.Vec3 `json:"position"`
}

func (s *sendler) sendBreakdown(id types.Entity, broadcaster func(EntityInfo) error) error {
	c, ok := s.GetComponent(id, "transform")
	if !ok {
		return errors.New("не найден компонент transform")
	}
	transform := c.(*components.TransformComponent)

	broadcaster(EntityInfo{
		ID:   id,
		Type: "breakdown",
		Data: breakdownData{
			Position: transform.Position,
		},
	})

	return nil
}
