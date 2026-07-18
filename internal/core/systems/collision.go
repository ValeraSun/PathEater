package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type CollisionSystem struct {
	getter     componentsGetter
	entitiesID []types.Entity
}

func NewCollisionSystem(getter componentsGetter) *CollisionSystem {
	return &CollisionSystem{
		getter:     getter,
		entitiesID: make([]types.Entity, 0, 128),
	}
}
func (s *CollisionSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("collider")

	// Структура для кэширования данных сущности
	type entityData struct {
		id        types.Entity // или ваш тип ID
		collider  *components.ColliderComponent
		transform *components.TransformComponent
	}

	// Предварительно выделяем память
	entities := make([]entityData, 0, len(comps))

	// Первый проход: собираем все данные и обновляем позиции коллайдеров
	for id, comp := range comps {
		collider, ok := comp.(*components.ColliderComponent)
		if !ok || collider == nil {
			continue
		}

		ed := entityData{
			id:       id,
			collider: collider,
		}

		// Получаем трансформ, если есть
		if s.getter.HasComponents(id, "transform") {
			c, _ := s.getter.GetComponent(id, "transform")
			if transform, ok := c.(*components.TransformComponent); ok {
				ed.transform = transform
				// Обновляем позицию коллайдера
				collider.Collider.ChangeCenter(transform.Position)
			}
		}

		entities = append(entities, ed)
	}

	// Второй проход: проверка коллизий
	for i := 0; i < len(entities); i++ {
		for j := i + 1; j < len(entities); j++ {
			e1, e2 := entities[i], entities[j]

			// Проверяем коллизию
			mtv, isColliding := e1.collider.Collide(e2.collider)
			if !isColliding {
				continue
			}

			// Обрабатываем разрешение коллизии
			switch {
			case e1.transform != nil && e2.transform != nil:
				// Оба двигаются
				e1.transform.Position.Add(mtv.Scale(0.5))
				e2.transform.Position.Add(mtv.Scale(-0.5))

			case e1.transform != nil:
				// Двигается только первый
				e1.transform.Position = e1.transform.Position.Add(mtv)

			case e2.transform != nil:
				// Двигается только второй
				e2.transform.Position = e2.transform.Position.Add(mtv.Scale(-1))
			}
		}
	}

	return nil
}
