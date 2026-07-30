package sendler

import (
	"errors"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

var shipCenter geometry.Vec3 = geometry.Vec3{X: 24.67, Y: 2.8, Z: -9.0}

type breakdownData struct {
	Position geometry.Vec3 `json:"position"`
	Normal   geometry.Vec3 `json:"normal,omitempty"`
}

func (s *sendler) sendBreakdown(id types.Entity, broadcaster func(EntityInfo) error) error {
	c, ok := s.GetComponent(id, "transform")
	if !ok {
		return errors.New("не найден компонент transform")
	}
	transform := c.(*components.TransformComponent)

	// Получаем нормаль стены (если есть)
	var normal geometry.Vec3
	if n, ok := s.GetComponent(id, "wallNormal"); ok {
		if colComp, ok := n.(*components.ColliderComponent); ok {
			box, _ := colComp.Collider.(*geometry.BoxCollider)
			normal = box.GetInwardNormal(shipCenter)
		}
	}

	broadcaster(EntityInfo{
		ID:   id,
		Type: "breakdown",
		Data: breakdownData{
			Position: transform.Position,
			Normal:   normal,
		},
	})

	return nil
}
