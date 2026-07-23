package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type NavigationDisplaySystem struct {
	getter componentsGetter
}

func NewNavigationDisplaySystem(getter componentsGetter, subscriber subscriber) *NavigationDisplaySystem {
	s := &NavigationDisplaySystem{
		getter: getter,
	}
	subscriber.Subscribe("setPlayerState", s.OnEvent)
	return s
}

func (s *NavigationDisplaySystem) Update(dt float32) error {
	return nil
}

func (s *NavigationDisplaySystem) OnEvent(event events.Event) error {
	ps, ok := event.(*events.SetPlayerStateEvent)

	if !ok {
		return nil
	}

	ships := s.getter.GetEntitiesByComponent("ship")
	var shipComp types.Component
	for _, sh := range ships {
		shipComp = sh
	}
	ship := shipComp.(*components.ShipComponent)

	if ps.PlayerState.Interact {
		switch ship.AvailableID {
		case "":
			ship.AvailableID = ps.ID
		case ps.ID:
			ship.AvailableID = ""
		}
		return nil
	}

	return nil
}
