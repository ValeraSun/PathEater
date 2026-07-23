package components

import (
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type NavigationDisplayComponent struct {
	AvailableID types.Entity
}

func (*NavigationDisplayComponent) Type() string {
	return "movement"
}

func NewNavigationDisplayComponent(id types.Entity) *NavigationDisplayComponent {
	return &NavigationDisplayComponent{
		AvailableID: id,
	}
}
