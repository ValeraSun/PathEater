package usecase

import (
    "context"
    "my-game-server/internal/domain/entity"
    "my-game-server/internal/domain/service"
)

// Broadcaster - интерфейс, который определяет, КАК отправлять данные (зависимость инвертирована)
type Broadcaster interface {
    SendToPlayer(ctx context.Context, playerID string, data []byte) error
    SendToRoom(ctx context.Context, roomID string, data []byte) error
}

// BroadcastUseCase - сценарий рассылки
type BroadcastUseCase struct {
    broadcaster Broadcaster
    playerRepo  PlayerRepository // интерфейс из domain/repository
}

func NewBroadcastUseCase(b Broadcaster, vis *service.VisibilityService, pRepo PlayerRepository) *BroadcastUseCase {
    return &BroadcastUseCase{broadcaster: b, playerRepo: pRepo}
}

// BroadcastToVisiblePlayers - рассылает данные только тем, кто видит игрока (AOI)
func (uc *BroadcastUseCase) BroadcastToVisiblePlayers(ctx context.Context, roomID string, sourcePlayer *entity.Player, data []byte) error {
    allPlayers, err := uc.playerRepo.FindByRoom(ctx, roomID)
    if err != nil {
        return err
    }

    for _, p := range visiblePlayers {
        go func(player *entity.Player) {
            uc.broadcaster.SendToPlayer(ctx, player.ID, data)
        }(p)
    }
    return nil
}