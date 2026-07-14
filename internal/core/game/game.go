package game

import (
	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/systems"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type systemAdder interface {
	AddSystem(types.System)
}

type componentsGetter interface {
	GetEntitiesByComponent(componentType string) map[types.Entity]types.Component
	HasComponents(entity types.Entity, componentTypes ...string) bool
	GetComponent(entity types.Entity, componentType string) (types.Component, bool)
}

func CreateGame(broadcaster ecs.Broadcaster) *ecs.World {
	w := ecs.CreateWorld(broadcaster)
	//eb := *w.EventBus
	initSystems(w, w)
	return w
}

func initSystems(adder systemAdder, getter componentsGetter) {
	adder.AddSystem(systems.NewPlayerControlSystem(getter))
	adder.AddSystem(systems.NewTransformSystem(getter))
}
