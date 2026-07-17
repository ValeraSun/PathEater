package transfer

import (
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type EventPublisher interface {
	Publish(events.Event) error
}

func SendPlayerState(publisher EventPublisher, state events.PlayerState, id types.Entity) error {
	e := events.NewSetPlayerStateEvent(state, id)
	return publisher.Publish(e)
}

func CreatePlayer(publisher EventPublisher, id string) error {
	e := events.NewCreatePlayerEvent(id)
	return publisher.Publish(e)
}

func SendWeaponState(publisher EventPublisher, state events.WeaponState) error {
	e := events.NewSetWeaponStateEvent(state)
	return publisher.Publish(e)
}
