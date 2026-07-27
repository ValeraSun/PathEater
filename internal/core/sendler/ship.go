package sendler

import (
	"errors"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type shipData struct {
	Collider        geometry.Collider `json:"collider"`
	WeaponDirection geometry.Vec3     `json:"weapon_direction"`
	ShootSuccess    bool              `json:"shoot_success"`
	BaggageStatus   int               `json:"baggage_status"`
	Health          int               `json:"health"`
}

func (s *sendler) sendShip(id types.Entity, broadcaster func(EntityInfo) error) error {

	c, ok := s.GetComponent(id, "ship")
	ship := c.(*components.ShipComponent)

	if !ok {
		return errors.New("Не найден необходимый копнонент")
	}

	c, ok = s.GetComponent(id, "health")
	hp := c.(*components.HealthComponent)

	if !ok {
		return errors.New("Не найден необходимый компнонент")
	}

	c, ok = s.GetComponent(id, "collider")
	col := c.(*components.ColliderComponent)

	if !ok {
		return errors.New("Не найден необходимый компнонент")
	}

	c, ok = s.GetComponent(id, "weapon")
	weap := c.(*components.WeaponComponent)

	if !ok {
		return errors.New("Не найден необходимый компнонент")
	}

	broadcaster(EntityInfo{
		ID:   id,
		Type: "ship",
		Data: shipData{
			Collider:        col.Collider,
			WeaponDirection: weap.Direction,
			ShootSuccess:    weap.ShootSuccess,
			BaggageStatus:   ship.BaggageStatus,
			Health:          hp.Health,
		}})

	return nil
}
