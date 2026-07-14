package network

import (
	"encoding/json"

	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/game"
	"github.com/ValeraSun/PathEater/internal/core/transfer"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type Command interface {
	Name() string
	Execute(client *Client, payload json.RawMessage) error
}

//Команды меню игры

// создание комнаты
type createRoomCommand struct{}

func (c *createRoomCommand) Name() string { return "createRoom" }

func (c *createRoomCommand) Execute(client *Client, payload json.RawMessage) error {
	// Создаём комнату с уникальным ID
	roomID := generateID()
	hub := GetHub()
	room := hub.CreateGameRoom(roomID)
	room.AddClient(client)

	var info struct {
		ID string `json:"id"`
	}
	response, err := json.Marshal(info)
	if err != nil {
		return err
	}

	client.SendMessage("SuccessCreateRoom", response)

	// Меняем состояние
	client.SetState(GameRoomState())
	return nil
}

// удаление комнаты
type deleteRoomCommand struct{}

func (c *deleteRoomCommand) Name() string { return "deleteRoom" }

func (c *deleteRoomCommand) Execute(client *Client, payload json.RawMessage) error {
	var roomID struct {
		RoomID string `json:"roomID"`
	}
	if err := json.Unmarshal(payload, &roomID); err != nil {
		return client.SendError(err)
	}
	hub := GetHub()

	room, exists := hub.GetGameRoom(roomID.RoomID)
	if !exists {
		return client.SendText("error", "RoomIsNotExists")
	}
	room.Close()

	var info struct {
		ID string `json:"id"`
	}
	response, err := json.Marshal(info)
	if err != nil {
		return err
	}

	client.SendMessage("SuccessDeleteRoom", response)

	// Меняем состояние
	if client.state == GameRoomState() {
		client.SetState(MainMenuState())
	}
	return nil
}

// присоединение к комнате
type joinRoomCommand struct{}

func (c *joinRoomCommand) Name() string { return "joinRoom" }

func (c *joinRoomCommand) Execute(client *Client, payload json.RawMessage) error {
	var roomID struct {
		RoomID string `json:"roomID"`
	}
	if err := json.Unmarshal(payload, &roomID); err != nil {
		return client.SendError(err)
	}
	hub := GetHub()
	room, exists := hub.GetGameRoom(roomID.RoomID)
	if !exists {
		return client.SendText("error", "RoomIsNotExists")
	}
	room.AddClient(client)

	var info struct {
		ID string `json:"id"`
	}
	response, err := json.Marshal(info)
	if err != nil {
		return err
	}

	client.SendMessage("SuccessJoinRoom", response)

	// Меняем состояние
	client.SetState(GameRoomState())
	return nil
}

// выход из меню (выход из всей игры) - отключение клиента
type exitMenuCommand struct{}

func (c *exitMenuCommand) Name() string { return "exitMenu" }

func (c *exitMenuCommand) Execute(client *Client, payload json.RawMessage) error {
	hub := GetHub()
	hub.UnregisterClient(client)
	client.Close()

	return nil
}

//Команды комнат

// Выход из комнаты
type exitRoomCommand struct{}

func (c *exitRoomCommand) Name() string { return "exitRoom" }

func (c *exitRoomCommand) Execute(client *Client, payload json.RawMessage) error {
	client.room.RemoveClient(client)

	return nil
}

type Sendler struct {
	room *GameRoom
}

func (s Sendler) SendEntityCreate(EntityInfo ecs.EntityInfo) error {
	var info struct {
		ID   types.Entity `json:"id"`
		Type string       `json:"type"`
		Data []byte       `json:"data"`
	}
	info.ID = EntityInfo.Id
	info.Type = EntityInfo.Type
	info.Data = EntityInfo.Data
	payload, err := json.Marshal(info)
	if err != nil {
		return err
	}

	s.room.SendToAll("CreateEntity", payload)
	return nil
}

func (s Sendler) SendEntityUpdate(EntityInfo ecs.EntityInfo) error {
	var info struct {
		ID   types.Entity `json:"id"`
		Type string       `json:"type"`
		Data []byte       `json:"data"`
	}
	info.ID = EntityInfo.Id
	info.Type = EntityInfo.Type
	info.Data = EntityInfo.Data
	payload, err := json.Marshal(info)
	if err != nil {
		return err
	}

	s.room.SendToAll("UpdateEntity", payload)
	return nil
}

func (s Sendler) SendEntityDelete(EntityInfo ecs.EntityInfo) error {
	var info struct {
		ID types.Entity `json:"id"`
	}
	info.ID = EntityInfo.Id
	payload, err := json.Marshal(info)
	if err != nil {
		return err
	}

	s.room.SendToAll("DeleteEntity", payload)
	return nil
}

func (s Sendler) SendSnapshotToAll(entities []ecs.EntityInfo) error {
	var info struct {
		Entities []ecs.EntityInfo `json:"entities"`
	}

	info.Entities = entities

	payload, err := json.Marshal(info)
	if err != nil {
		return err
	}

	s.room.SendToAll("Snapshot", payload)
	return nil
}

// начало игры
type createGameSessionCommand struct{}

func (c *createGameSessionCommand) Name() string { return "startGame" }

func (c *createGameSessionCommand) Execute(client *Client, payload json.RawMessage) error {
	broadcaster := Sendler{room: client.room}
	game.CreateGame(broadcaster)

	for _, client := range client.room.Clients {
		client.SetState(PlayerControlState())
		transfer.CreatePlayer(broadcaster, client.ID)
	}

	return nil
}

//Команды игры

type playerStateCommand struct{}

func (c *playerStateCommand) Name() string { return "playerState" }

func (c *playerStateCommand) Execute(client *Client, payload json.RawMessage) error {
	var state events.PlayerState

	if err := json.Unmarshal(payload, &state); err != nil {
		return client.SendError(err)
	}

	transfer.SendPlayerState(client.room.World.EventBus, state, types.Entity(client.ID))

	return nil
}

type exitGameCommand struct{}

func (c *exitGameCommand) Name() string { return "exitGame" }

func (c *exitGameCommand) Execute(client *Client, payload json.RawMessage) error {
	return nil
}
