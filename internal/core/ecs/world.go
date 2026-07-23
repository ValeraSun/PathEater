package ecs

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ValeraSun/PathEater/internal/config"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type World struct {
	entities map[types.Entity]map[string]types.Component

	componentIndex map[string]map[types.Entity]struct{}

	systems []types.System

	EventBus *events.EventBus

	mu sync.RWMutex

	entityCount  int64
	systemTimers map[string]time.Duration

	Broadcaster Broadcaster

	done chan struct{}
}

type Broadcaster interface {
	SendEntityCreate(EntityInfo) error
	SendEntityUpdate(EntityInfo) error
	SendEntityDelete(EntityInfo) error
	SendSnapshotToAll([]EntityInfo) error
}

type EntityInfo struct {
	ID   types.Entity
	Type string
	Data any
}

func newWorld(eventBus *events.EventBus, room Broadcaster) *World {
	return &World{
		entities:       make(map[types.Entity]map[string]types.Component),
		componentIndex: make(map[string]map[types.Entity]struct{}),
		systems:        make([]types.System, 0),
		EventBus:       eventBus,
		systemTimers:   make(map[string]time.Duration),
		Broadcaster:    room,
	}
}

func CreateWorld(room Broadcaster) *World {
	eb := events.NewEventBus(100)
	w := newWorld(eb, room)

	return w
}

func (w *World) Close() {
	w.done <- struct{}{}
}

func HandleWorld(world *World) {
	ticker := time.NewTicker(config.GetNanosecondPerTick())
	defer ticker.Stop()

	lastTick := time.Now()

	for {
		select {
		case <-world.done:
			world.EventBus.Close()
			return
		case now := <-ticker.C:
			dt := time.Since(lastTick).Seconds()
			lastTick = now

			world.Update(float32(dt))
		}
	}
}

func (w *World) AddSystem(system types.System) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.systems = append(w.systems, system)
}

func (w *World) RemoveSystem(system types.System) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for i, s := range w.systems {
		if s == system {
			w.systems = append(w.systems[:i], w.systems[i+1:]...)
			return
		}
	}
}

// GetSystems возвращает копию списка систем (безопасно для чтения)
func (w *World) GetSystems() []types.System {
	w.mu.RLock()
	defer w.mu.RUnlock()

	systems := make([]types.System, len(w.systems))
	copy(systems, w.systems)
	return systems
}

func (w *World) AddEntity(components ...types.Component) (types.Entity, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	entity := types.NewEntity()

	w.entities[entity] = make(map[string]types.Component)

	for _, comp := range components {
		if err := w.addComponentToEntity(entity, comp); err != nil {
			return "", fmt.Errorf("failed to add component: %w", err)
		}
	}

	w.entityCount++
	return entity, nil
}

func (w *World) AddEntityByID(entity types.Entity, components ...types.Component) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, exists := w.entities[entity]; exists {
		return errors.New("уже существует entity с таким id")
	}

	w.entities[entity] = make(map[string]types.Component)
	for _, comp := range components {
		if err := w.addComponentToEntity(entity, comp); err != nil {
			return fmt.Errorf("failed to add component: %w", err)
		}
	}
	w.entityCount++
	return nil
}

func (w *World) addComponentToEntity(entity types.Entity, comp types.Component) error {
	compType := comp.Type()

	w.entities[entity][compType] = comp

	if w.componentIndex[compType] == nil {
		w.componentIndex[compType] = make(map[types.Entity]struct{})
	}
	w.componentIndex[compType][entity] = struct{}{}

	return nil
}

func (w *World) GetEntity(entity types.Entity) (map[string]types.Component, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	comps, exists := w.entities[entity]
	return comps, exists
}

func (w *World) GetComponent(entity types.Entity, componentType string) (types.Component, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if comps, exists := w.entities[entity]; exists {
		comp, ok := comps[componentType]
		return comp, ok
	}
	return nil, false
}

func (w *World) GetEntitiesByComponent(componentType string) map[types.Entity]types.Component {
	w.mu.RLock()
	defer w.mu.RUnlock()

	result := make(map[types.Entity]types.Component)

	if entities, exists := w.componentIndex[componentType]; exists {
		for entity := range entities {
			if comp, ok := w.entities[entity][componentType]; ok {
				result[entity] = comp
			}
		}
	}

	return result
}

func (w *World) Update(dt float32) error {
	w.mu.RLock()
	systems := make([]types.System, len(w.systems))
	copy(systems, w.systems)
	w.mu.RUnlock()

	for _, system := range systems {
		start := time.Now()

		if err := system.Update(dt); err != nil {
			return fmt.Errorf("system %T failed: %w", system, err)
		}

		w.mu.Lock()
		w.systemTimers[fmt.Sprintf("%T", system)] = time.Since(start)
		w.mu.Unlock()
	}

	return nil
}

func (w *World) RemoveEntity(entity types.Entity) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if comps, exists := w.entities[entity]; exists {
		for compType := range comps {
			delete(w.componentIndex[compType], entity)
		}
	}

	delete(w.entities, entity)

	w.entityCount--
}

func (w *World) HasComponents(entity types.Entity, componentTypes ...string) bool {
	for _, comp := range componentTypes {
		_, exists := w.entities[entity][comp]
		if !exists {
			return false
		}
	}
	return true
}
