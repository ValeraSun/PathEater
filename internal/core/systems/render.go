package systems

import (
	"sync"

	"github.com/ValeraSun/PathEater/internal/core/types"
)

type RenderSystem struct {
	componentsGetter
	broadcaster
}

func NewRenderSystem(getter componentsGetter, broadcaster broadcaster) *RenderSystem {
	return &RenderSystem{
		getter,
		broadcaster,
	}
}

func (s *RenderSystem) Update(dt float32) error {
	comps := s.GetEntitiesByComponent("update")

	var wg sync.WaitGroup
	errChan := make(chan error, len(comps))
	done := make(chan struct{})

	for id := range comps {
		var err error
		wg.Add(1)

		go func(entityID types.Entity) {
			defer wg.Done()

			s.Send(id, s.SendEntityUpdate)

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
