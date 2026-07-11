package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/types"
	"github.com/google/uuid"
)

func NewEntity() types.Entity {
	return types.Entity(uuid.New().String())
}
