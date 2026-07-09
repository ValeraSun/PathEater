package ecs

import (
	"fmt"
	"sync"
	"time"

	"github.com/ValeraSun/PathEater/internal/config"
	"github.com/ValeraSun/PathEater/internal/core/events"
)

type World struct {
	// entities хранит все сущности и их компоненты
	entities map[Entity]map[string]Component

	// componentIndex позволяет быстро искать сущности по типу компонента
	componentIndex map[string]map[Entity]struct{}

	// systems выполняются в порядке добавления
	systems []System

	// eventBus для межсистемной коммуникации
	EventBus *events.EventBus

	// entityPool для переиспользования удаленных сущностей
	entityPool sync.Pool

	mu sync.RWMutex

	// Метрики
	entityCount  int64
	systemTimers map[string]time.Duration
}

func newWorld(eventBus *events.EventBus) *World {
	return &World{
		entities:       make(map[Entity]map[string]Component),
		componentIndex: make(map[string]map[Entity]struct{}),
		systems:        make([]System, 0),
		EventBus:       eventBus,
		systemTimers:   make(map[string]time.Duration),
		entityPool: sync.Pool{
			New: func() interface{} {
				return NewEntity()
			},
		},
	}
}

func CreateWorld() *World {
	eb := events.NewEventBus(100)
	w := newWorld(eb)
	go HandleWorld(w)
	return w
}

func HandleWorld(world *World) {
	ticker := time.NewTicker(config.GetMillisecondPerTick() * time.Millisecond)
	defer ticker.Stop() // Важно: останавливаем тикер при выходе

	// Бесконечный цикл для обработки тиков
	for range ticker.C {
		go world.Update(float64(config.GetMillisecondPerTick()))
	}
}

// AddEntity добавляет сущность с компонентами
func (w *World) AddEntity(components ...Component) (Entity, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	entity := w.entityPool.Get().(Entity)

	w.entities[entity] = make(map[string]Component)

	for _, comp := range components {
		if err := w.addComponentToEntity(entity, comp); err != nil {
			return "", fmt.Errorf("failed to add component: %w", err)
		}
	}

	w.entityCount++
	return entity, nil
}

// addComponentToEntity добавляет компонент и обновляет индексы
func (w *World) addComponentToEntity(entity Entity, comp Component) error {
	compType := comp.Type()

	// Добавляем в сущность
	w.entities[entity][compType] = comp

	// Обновляем индекс
	if w.componentIndex[compType] == nil {
		w.componentIndex[compType] = make(map[Entity]struct{})
	}
	w.componentIndex[compType][entity] = struct{}{}

	return nil
}

// GetEntity получает все компоненты сущности
func (w *World) GetEntity(entity Entity) (map[string]Component, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	comps, exists := w.entities[entity]
	return comps, exists
}

// GetComponent получает конкретный компонент сущности
func (w *World) GetComponent(entity Entity, componentType string) (Component, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if comps, exists := w.entities[entity]; exists {
		comp, ok := comps[componentType]
		return comp, ok
	}
	return nil, false
}

// GetEntitiesByComponent возвращает все сущности с указанным компонентом
func (w *World) GetEntitiesByComponent(componentType string) map[Entity]Component {
	w.mu.RLock()
	defer w.mu.RUnlock()

	result := make(map[Entity]Component)

	if entities, exists := w.componentIndex[componentType]; exists {
		for entity := range entities {
			if comp, ok := w.entities[entity][componentType]; ok {
				result[entity] = comp
			}
		}
	}

	return result
}

// Update выполняет все системы
func (w *World) Update(dt float64) error {
	w.mu.RLock()
	systems := make([]System, len(w.systems))
	copy(systems, w.systems)
	w.mu.RUnlock()

	for _, system := range systems {
		start := time.Now()

		if err := system.Update(w, dt); err != nil {
			return fmt.Errorf("system %T failed: %w", system, err)
		}

		w.mu.Lock()
		w.systemTimers[fmt.Sprintf("%T", system)] = time.Since(start)
		w.mu.Unlock()
	}

	return nil
}

// RemoveEntity удаляет сущность и возвращает её в пул
func (w *World) RemoveEntity(entity Entity) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Удаляем из индексов
	if comps, exists := w.entities[entity]; exists {
		for compType := range comps {
			delete(w.componentIndex[compType], entity)
		}
	}

	// Удаляем сущность
	delete(w.entities, entity)

	// Возвращаем ID в пул
	w.entityPool.Put(entity)
	w.entityCount--
}
