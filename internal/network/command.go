package network

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ValeraSun/PathEater/internal/core/ecs"
)

type Command interface {
	Name() string
	Execute(ctx context.Context, client *Client, payload json.RawMessage) error
}

// -----------------------------------------------------------------------------------------------------------------------------------
// Команды старта игры
type StartRoomCommand struct{}

func (c *StartRoomCommand) Name() string { return "StartRoom" }

func (c *StartRoomCommand) Execute(ctx context.Context, client *Client, payload json.RawMessage) error {
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

type StartWorldCommand struct{}

func (c *StartWorldCommand) Name() string { return "StartGame" }

func (c *StartWorldCommand) Execute(ctx context.Context, client *Client, payload json.RawMessage) error {

	w := ecs.CreateWorld()
	room := client.room

	if room == nil {
		return errors.New("Ошибка комната не найдена")
	}

	room.World = w
	return nil
}

// ------------------------------------------------------------------------------------------------------------------
// обработка команд движения
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
	return nil

	//Передача события движения
	e := CreateEventMove(coord)
	client.room.World.EventBus.Publish(e)

	return nil
}

// -----------------------------------------------------------------------------------------------------------------------------------
// Команды взаимодействия с предметами
type UseItemCommand struct{}

func (*UseItemCommand) Name() string { return "UseItem" }

func (*UseItemCommand) Execute(ctx context.Context, client *Client, payload json.RawMessage) error {
	var data struct {
		ItemID string `json:"itemID"`
	}
	if err := json.Unmarshal(payload, &data); err != nil {
		return client.SendError(err)
	}

	//Передача события использования
	e := CreateEventItemUse(data)
	client.room.World.EventBus.Publish(e)

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

	//Передача события использования
	e := CreateEventExit()
	client.room.World.EventBus.Publish(e)

	return nil
}
