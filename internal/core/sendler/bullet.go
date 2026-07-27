package sendler

import (
	"errors"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type bulletData struct {
	Position geometry.Vec3 `json:"position"`
	Success  bool          `json:"success"`
}

func (s *sendler) sendBullet(id types.Entity, broadcaster func(EntityInfo) error) error {
	shipPos, err := getShipPosition(s)
	if err != nil {
		return err
	}

	c, ok := s.GetComponent(id, "transform")
	if !ok {
		return errors.New("не найден компонент transform")
	}
	transform := c.(*components.TransformComponent)

	relPos := geometry.Vec3{
		X: transform.Position.X - shipPos.X,
		Y: transform.Position.Y - shipPos.Y,
		Z: transform.Position.Z - shipPos.Z,
	}

	broadcaster(EntityInfo{
		ID:   id,
		Type: "bullet",
		Data: bulletData{
			Position: relPos,
		},
	})

	return nil
}

func getShipPosition(getter componentsGetter) (geometry.Vec3, error) {
	entities := getter.GetEntitiesByComponent("ship")
	for id := range entities {
		if comp, ok := getter.GetComponent(id, "transform"); ok {
			trans := comp.(*components.TransformComponent)
			return trans.Position, nil
		}
	}
	return geometry.Vec3{}, errors.New("ship not found")
}
