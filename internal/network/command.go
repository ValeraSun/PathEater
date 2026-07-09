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
<<<<<<< HEAD
	// Создаем комнату с уникальным ID
	roomID := generateID()
	hub := GetHub()
	room := hub.CreateGameRoom(roomID)
	room.AddClient(client)

	// Меняем состояние
	client.SetState(NewGameState(roomID))
	client.state.OnEnter(ctx, client)
	return nil
=======
    // Создаем комнату с уникальным ID
    roomID := generateID()
    hub := GetHub()
    room := hub.CreateGameRoom(roomID)
    room.AddClient(client)
    
    //Передача события старта
    e := CreateEventItemUse(itemID)
    client.room.World.eventBus.Publish(e)

    // Меняем состояние
    client.SetState(NewGameState(roomID))
    client.state.OnEnter(ctx, client)
    return nil
>>>>>>> origin/Egor
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
	var coord struct {
		ItemID string `json:"itemID"`
	}
	if err := json.Unmarshal(payload, &coord); err != nil {
		return client.SendError(err)
	}

	//Передача события использования
	e := CreateEventItemUse(coord)
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

<<<<<<< HEAD
	//Передача события использования
	e := CreateEventExit()
	client.room.World.EventBus.Publish(e)
=======
    //Передача события выхода
    e := CreateEventExit()
    client.room.World.eventBus.Publish(e)
>>>>>>> origin/Egor

	return nil
}
