package systems

import (
	"sync"

	"github.com/ValeraSun/PathEater/internal/core/entities"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type RenderSystem struct {
	getter      componentsGetter
	broadcaster Broadcaster
}

func NewRenderSystem(getter componentsGetter, broadcaster Broadcaster) *RenderSystem {
	return &RenderSystem{
		getter:      getter,
		broadcaster: broadcaster,
	}
}

func (s *RenderSystem) Update(dt float32) error {
	comps := s.getter.GetEntitiesByComponent("update")

	var wg sync.WaitGroup
	errChan := make(chan error, len(comps))
	done := make(chan struct{})

	for id := range comps {
		var err error
		wg.Add(1)
		go func(entityID types.Entity) {
			defer wg.Done()

			switch {
			case entities.IsPlayer(id, s.getter):
				err = entities.SendPlayer(id, s.getter, s.broadcaster.SendEntityUpdate)
			case entities.IsShip(id, s.getter):
				err = entities.SendShip(id, s.getter, s.broadcaster.SendEntityUpdate)
			case s.getter.HasComponents(id, "ai"):
				err = entities.SendAlien(id, s.getter, s.broadcaster.SendEntityUpdate)
			case entities.IsAsteroid(id, s.getter):
				err = entities.SendAsteroid(id, s.getter, s.broadcaster.SendEntityUpdate)
			case entities.IsCosmoAlien(id, s.getter):
				err = entities.SendCosmoAlien(id, s.getter, s.broadcaster.SendEntityUpdate)
			case entities.IsBullet(id, s.getter):
				err = entities.SendBullet(id, s.getter, s.broadcaster.SendEntityUpdate)
			}

			if err != nil {
				errChan <- err
			}

		}(id)

	}

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	for {
		select {
		case e := <-errChan:
			return e
		case <-done:
			close(errChan)
			return nil
		}

	}
}
