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
