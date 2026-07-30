package sendler

import (
	"errors"
	"log"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type doorData struct {
	Position geometry.Vec3 `json:"position"`
	Rotation geometry.Vec3 `json:"rotation"`
	IsOpen   bool          `json:"isOpen"`
}

func (s *sendler) sendDoor(id types.Entity, broadcaster func(EntityInfo) error) error {

	c, ok := s.GetComponent(id, "collider")
	if !ok {
		return errors.New("не найден компонент collider")
	}
	col := c.(*components.ColliderComponent)
	box := col.Collider.(*geometry.BoxCollider)

	eulerX, eulerY, eulerZ := geometry.MatrixToEuler(box.Axes)
	rotation := geometry.Vec3{X: eulerX, Y: eulerY, Z: eulerZ}

	c, ok = s.GetComponent(id, "door")
	if !ok {
		return errors.New("не найден компонент door")
	}
	door := c.(*components.DoorComponent)

	log.Println("ПОСЛАНА ДВЕРЬ: ", door.IsOpen)
	broadcaster(EntityInfo{
		ID:   id,
		Type: "door",
		Data: doorData{
			Position: col.Collider.GetCenter(),
			Rotation: rotation,
			IsOpen:   door.IsOpen,
		},
	})

	return nil
}
