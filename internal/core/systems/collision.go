package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/transfer"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type CollisionSystem struct {
	getter     componentsGetter
	publisher  transfer.EventPublisher
	entitiesID []types.Entity
}

func NewCollisionSystem(getter componentsGetter, publisher transfer.EventPublisher) *CollisionSystem {
	return &CollisionSystem{
		getter:     getter,
		publisher:  publisher,
		entitiesID: make([]types.Entity, 0, 128),
	}
}
func (s *CollisionSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("collider")

	// Структура для кэширования данных сущности
	type entityData struct {
		id          types.Entity // или ваш тип ID
		collider    *components.ColliderComponent
		transform   *components.TransformComponent
		movable     *components.MovableComponent
		OnCollision func(id, other types.Entity)
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
			id:          id,
			collider:    collider,
			OnCollision: func(id, other types.Entity) {}
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

		if s.getter.HasComponents(id, "movable") {
			c, _ := s.getter.GetComponent(id, "movable")
			if movable, ok := c.(*components.MovableComponent); ok {
				ed.movable = movable
			}
		}

		if s.getter.HasComponents(id, "ship") {
			ed.OnCollision = func(id, other types.Entity) {
				switch {
				case s.getter.HasComponents(other, "asteroid"):
					e1, e2 := shipAsteroidCollision(id, other)

					s.publisher.Publish(e1)
					s.publisher.Publish(e2)
				case s.getter.HasComponents(other, "cosmoAlien"):
					e1, e2 := shipCosmoAlienCollision(id, other)

					s.publisher.Publish(e1)
					s.publisher.Publish(e2)
				}
			}

			if s.getter.HasComponents(id, "asteroid") {
				ed.OnCollision = func(id, other types.Entity) {
					comp, _ := s.getter.GetComponent(id, "asteroid")
					aster := comp.(*components.AsteroidComponent)
					aster.Destroyed = true
					switch {
					case s.getter.HasComponents(other, "ship"):
						e1, e2 := shipAsteroidCollision(other, id)

						s.publisher.Publish(e1)
						s.publisher.Publish(e2)
					case s.getter.HasComponents(other, "asteroid"):
						e1, e2 := asteroidAsteroidCollision(id, other)

						s.publisher.Publish(e1)
						s.publisher.Publish(e2)
					case s.getter.HasComponents(other, "cosmoAlien"):
						e1, e2 := asteroidCosmoAlienCollision(id, other)

						s.publisher.Publish(e1)
						s.publisher.Publish(e2)
					case s.getter.HasComponents(other, "bullet"):
						e1, e2 := asteroidBulletCollision(id, other)

						s.publisher.Publish(e1)
						s.publisher.Publish(e2)
					}
				}
			}

			if s.getter.HasComponents(id, "cosmoAlien") {
				ed.OnCollision = func(id, other types.Entity) {
					switch {
					case s.getter.HasComponents(other, "ship"):
						e1, e2 := shipCosmoAlienCollision(other, id)

						s.publisher.Publish(e1)
						s.publisher.Publish(e2)
					case s.getter.HasComponents(other, "asteroid"):
						e1, e2 := asteroidCosmoAlienCollision(other, id)

						s.publisher.Publish(e1)
						s.publisher.Publish(e2)
					case s.getter.HasComponents(other, "bullet"):
						e1, e2 := cosmoAlienBulletCollision(id, other)

						s.publisher.Publish(e1)
						s.publisher.Publish(e2)
					}
				}
			}

			if s.getter.HasComponents(id, "bullet") {
				ed.OnCollision = func(id, other types.Entity) {
					switch {
					case s.getter.HasComponents(other, "asteroid"):
						e1, e2 := asteroidBulletCollision(id, other)

						s.publisher.Publish(e1)
						s.publisher.Publish(e2)
					case s.getter.HasComponents(other, "cosmoAlien"):
						e1, e2 := cosmoAlienBulletCollision(other, id)

						s.publisher.Publish(e1)
						s.publisher.Publish(e2)
					}
				}
			}

			entities = append(entities, ed)
		}
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

			e1Movable := e1.transform != nil && s.getter.HasComponents(e1.id, "movable")
			e2Movable := e2.transform != nil && s.getter.HasComponents(e2.id, "movable")
			// Обрабатываем разрешение коллизии
			switch {
			case e1Movable && e2Movable:
				// Оба двигаются
				e1.transform.Position.Add(mtv.Scale(0.5))
				e2.transform.Position.Add(mtv.Scale(-0.5))

			case e1Movable:
				// Двигается только первый
				e1.transform.Position = e1.transform.Position.Add(mtv)

			case e2Movable:
				// Двигается только второй
				e2.transform.Position = e2.transform.Position.Add(mtv.Scale(-1))
			}
			e1.OnCollision(e1.id, e2.id)
		}
	}
	return nil
}

func shipAsteroidCollision(shipID, asteroidID types.Entity) (e1 *events.DamageShipEvent, e2 *events.DeleteAsteroidEvent) {
	return events.NewDamageShipEvent(shipID), events.NewDeleteAsteroidEvent(asteroidID)
}

func shipCosmoAlienCollision(shipID, alienID types.Entity) (e1 *events.DamageShipEvent, e2 *events.DeleteCosmoAlienEvent) {
	return events.NewDamageShipEvent(shipID), events.NewDeleteCosmoAlienEvent(alienID)
}

func asteroidAsteroidCollision(asteroid1ID, asteroid2ID types.Entity) (e1 *events.DeleteAsteroidEvent, e2 *events.DeleteAsteroidEvent) {
	return events.NewDeleteAsteroidEvent(asteroid1ID), events.NewDeleteAsteroidEvent(asteroid2ID)
}

func asteroidCosmoAlienCollision(asteroidID, alienID types.Entity) (e1 *events.DeleteAsteroidEvent, e2 *events.DeleteCosmoAlienEvent) {
	return events.NewDeleteAsteroidEvent(asteroidID) , events.NewDeleteAsteroidEvent(asteroid2ID)
}

func asteroidBulletCollision(asteroidID, bulletID types.Entity) (e1 *events.DeleteAsteroidEvent, e2 *events.DeleteBulletEvent) {
	return events.NewDeleteAsteroidEvent(asteroidID), events.NewDeleteBulletEvent(bulletID)
}

func cosmoAlienBulletCollision(alienID, bulletID types.Entity) ( e1 *events.DamageCosmoAienEvent, e2 *events.DeleteBulletEvent) {
	return events.NewDamageCosmoAlienEvent(shipID), events.NewDeleteBulletEvent(alienID)
}
