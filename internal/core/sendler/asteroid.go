package sendler

import (
	"errors"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type asteroidData struct {
	Position  geometry.Vec3 `json:"position"`
	Radius    float64       `json:"radius"`
	Destroyed bool          `json:"destroyed"`
}

func (s *sendler) sendAsteroid(id types.Entity, broadcaster func(EntityInfo) error) error {
	c, ok := s.GetComponent(id, "asteroid")
	if !ok {
		return errors.New("не найден компонент asteroid")
	}
	aster := c.(*components.AsteroidComponent)

	c, ok = s.GetComponent(id, "transform")
	if !ok {
		return errors.New("не найден компонент transform")
	}
	transform := c.(*components.TransformComponent)

	c, ok = s.GetComponent(id, "collider")
	if !ok {
		return errors.New("не найден компонент collider")
	}
	col := c.(*components.ColliderComponent)

	broadcaster(EntityInfo{
		ID:   id,
		Type: "asteroid",
		Data: asteroidData{
			Position:  transform.Position,
			Radius:    col.Collider.(*geometry.CircleCollider).Radius,
			Destroyed: aster.Destroyed,
		},
	})

	return nil
}
