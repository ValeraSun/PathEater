package network

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type Command interface {
	Name() string
	Execute(ctx context.Context, client *Client, payload json.RawMessage) error
}

// -----------------------------------------------------------------------------------------------------------------------------------
// Команды старта игры
type createRoomCommand struct{}

func (c *createRoomCommand) Name() string { return "createRoom" }

func (c *createRoomCommand) Execute(ctx context.Context, client *Client, payload json.RawMessage) error {
	// Создаем комнату с уникальным ID
	roomID := generateID()
	hub := GetHub()
	room := hub.CreateGameRoom(roomID)
	room.AddClient(client)

	// Меняем состояние
	client.SetState(NewRoomMenuState())
	client.state.OnEnter(ctx, client)
	return nil
}

type createGameSessionCommand struct{}

func (c *createGameSessionCommand) Name() string { return "createGameSession" }

func (c *createGameSessionCommand) Execute(ctx context.Context, client *Client, payload json.RawMessage) error {

	w := ecs.CreateWorld()
	room := client.room

	if room == nil {
		return errors.New("Ошибка комната не найдена")
	}

	room.World = w

	client.SetState(NewPlayerControlState())
	client.state.OnEnter(ctx, client)

	return nil
}

// ------------------------------------------------------------------------------------------------------------------
// обработка команд движения
type MoveCommand struct{}

func (c *MoveCommand) Name() string { return "move" }

func (*MoveCommand) Execute(ctx context.Context, client *Client, payload json.RawMessage) error {
	var transform struct {
		Id        types.Entity
		Position  geometry.Vector3
		Direction geometry.Vector3
	}
	if err := json.Unmarshal(payload, &transform); err != nil {
		log.Println(err)
		return nil
	}

	//Передача события движения
	e := events.CreateEventMove(transform.Position, transform.Direction, transform.Id)
	client.room.World.EventBus.Publish(e)

	return nil
}

// -----------------------------------------------------------------------------------------------------------------------------------
// Команды взаимодействия с предметами
type UseItemCommand struct{}

func (*UseItemCommand) Name() string { return "useItem" }

func (*UseItemCommand) Execute(ctx context.Context, client *Client, payload json.RawMessage) error {
	var data struct {
		ItemID string `json:"itemID"`
	}
	if err := json.Unmarshal(payload, &data); err != nil {
		return client.SendError(err)
	}

	//Передача события использования
	// e := CreateEventItemUse(data)
	// client.room.World.EventBus.Publish(e)

	return nil
}

// ------------------------------------------------------------------------------------------------------------------
// обработка команд выхода
type ExitGameSeccionCommand struct{}

func (*ExitGameSeccionCommand) Name() string { return "exitGameSeccion" }

func (*ExitGameSeccionCommand) Execute(ctx context.Context, client *Client, payload json.RawMessage) error {
	client.state.OnExit(ctx, client)
	client.SetState(NewMainMenuState())
	client.state.OnEnter(ctx, client)

	//Передача события использования
	// e := CreateEventExit()
	// client.room.World.EventBus.Publish(e)

	return nil
}
