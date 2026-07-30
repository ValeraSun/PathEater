package entities

import (
	"github.com/ValeraSun/PathEater/internal/core/types"
)

func IsDoor(id types.Entity, getter componentsGetter) bool {
	return getter.HasComponents(id, "door")
}
