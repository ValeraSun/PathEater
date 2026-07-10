// internal/core/ecs/types.go
package ecs

import (
	"github.com/google/uuid"
)

// Entity - уникальный идентификатор сущности
type Entity string

// Component - интерфейс для всех компонентов
type Component interface {
	Type() string
}

// System - интерфейс для всех систем
type System interface {
	Update(world *World, dt float64) error
}

// NewEntity создает новый уникальный идентификатор
func NewEntity() Entity {
	return Entity(uuid.New().String())
}
