package events

import "github.com/ValeraSun/PathEater/internal/core/types"

type WeaponEvent struct {
	ID       types.Entity
	Left     bool
	Right    bool
	Shooting bool
}

func (*WeaponEvent) Type() string { return "weapon" }
