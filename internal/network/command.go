package network

import (
	"context"
	"encoding/json"
)

type Command interface {
	Name() string
	Execute(ctx context.Context, client *Client, payload json.RawMessage) error
}

// -----------------------------------------------------------------------------------------------------------------------------------
// Команды старта игры
type StartGameCommand struct{}

func (c *StartGameCommand) Name() string { return "StartGame" }

func (c *StartGameCommand) Execute(ctx context.Context, client *Client, payload json.RawMessage) error {
	// Создаем комнату с уникальным ID
	roomID := generateID()
	hub := GetHub()
	room := hub.CreateGameRoom(roomID)
	room.AddClient(client)

	// Меняем состояние
	client.SetState(NewGameState(roomID))
	client.state.OnEnter(ctx, client)
	return nil
}

type MoveCommand struct{}

func (c *MoveCommand) Name() string { return "Move" }

func (*MoveCommand) Execute(ctx context.Context, client *Client, payload json.RawMessage) error {
	var coord struct {
		X float32 `json:"X"`
		Y float32 `json:"Y"`
		Z float32 `json:"Z"`
	}
	if err := json.Unmarshal(payload, &coord); err != nil {
		return client.SendError(err)
	}

	//посылаем координаты на основную логику сервера

	return nil
}

// -----------------------------------------------------------------------------------------------------------------------------------
// Команды взаимодействия с предметами
type UseItemCommand struct{}

func (*UseItemCommand) Name() string { return "UseItem" }

func (*UseItemCommand) Execute(ctx context.Context, client *Client, payload json.RawMessage) error {
	var coord struct {
		itemID string `json:"itemID"`
	}
	if err := json.Unmarshal(payload, &coord); err != nil {
		return client.SendError(err)
	}
	//посылаем ID предмета на основную логику сервера
	return nil
}

// ------------------------------------------------------------------------------------------------------------------
// обработка команд выхода
type ExitCommand struct{}

func (*ExitCommand) Name() string { return "ExitGame" }

func (*ExitCommand) Execute(ctx context.Context, client *Client, payload json.RawMessage) error {
	client.state.OnExit(ctx, client)
	client.SetState(NewMenuState())
	client.state.OnEnter(ctx, client)
	return nil
}
