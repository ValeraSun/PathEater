package transfer

import (
	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

func SendPlayerState(w *ecs.World, ps events.PlayerState, id types.Entity) {
	e := events.NewSetPlayerStateEvent(ps, id)
	w.EventBus.Publish(e)
}
