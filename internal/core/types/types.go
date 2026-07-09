package types

import (
	"github.com/google/uuid"
)

type Entity string

type Component interface {
	Type() string
}

type System interface {
	Update(dt float32) error
}

func NewEntity() Entity {
	return Entity(uuid.New().String())
}
