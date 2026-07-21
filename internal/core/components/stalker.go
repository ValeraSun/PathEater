package components

import "github.com/ValeraSun/PathEater/internal/core/types"

type StalkerComponent struct {
	Target types.Entity
}

func NewStalkerComponent(target types.Entity) *StalkerComponent {
	return &StalkerComponent{
		Target: target,
	}
}

func (*StalkerComponent) Type() string { return "stalker" }
