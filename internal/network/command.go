package network

import (
	"encoding/json"
	"errors"

	"github.com/ValeraSun/PathEater/internal/core/ecs"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/game"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/transfer"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

type Command interface {
	Name() string
	Execute(client *Client, payload json.RawMessage) error
}

var defaultSpawn = map[string]float64{"x": 0, "y": 1, "z": 0}

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

	client.SendMessage("SuccessCreateRoom", map[string]string{
		"roomId":   roomID,
		"playerId": client.ID,
	})

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
		return client.SendMessage("SuccessDeleteRoom", map[string]string{
			"error": "комната не найдена",
		})
	}
	room.Close()

	client.SendMessage("SuccessDeleteRoom", map[string]string{
		"roomId": roomID.RoomID,
	})

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
		return client.SendMessage("SuccessJoinRoom", map[string]string{
			"error": "комната не найдена",
		})
	}
	room.AddClient(client)

	client.SendMessage("SuccessJoinRoom", map[string]string{
		"roomId":   roomID.RoomID,
		"playerId": client.ID,
	})

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
	if client.room == nil {
		return client.SendError(errors.New("клиент не находится в комнате"))
	}

	client.room.RemoveClient(client)
	client.SetState(MainMenuState())

	return nil
}

type Sendler struct {
	room *GameRoom
}

func (s Sendler) SendEntityCreate(EntityInfo ecs.EntityCreateInfo) error {
	var info struct {
		ID        types.Entity  `json:"id"`
		Type      string        `json:"type"`
		Mesh      string        `json:"mesh"`
		Position  geometry.Vec3 `json:"position"`
		Direction geometry.Vec3 `json:"direction"`
	}

	info.ID = EntityInfo.ID
	info.Type = "object"
	info.Mesh = EntityInfo.Mesh
	info.Position = EntityInfo.Position
	info.Direction = EntityInfo.Direction

	payload, err := json.Marshal(info)
	if err != nil {
		return err
	}

	s.room.SendToAll("CrateEntity", payload)

	return nil
}

// func toEntityInfoJSON(info ecs.EntityInfo) entityInfoJSON {
// 	return entityInfoJSON{
// 		ID:   info.Id,
// 		Type: info.Type,
// 		Data: json.RawMessage(info.Data),
// 	}
// }

// func (s Sendler) SendEntityCreate(EntityInfo ecs.EntityInfo) error {
// 	s.room.SendToAll("CreateEntity", toEntityInfoJSON(EntityInfo))
// 	return nil
// }

func (s Sendler) SendEntityUpdate(entityInfo ecs.EntityUpdateInfo) error {
	var info struct {
		ID        types.Entity  `json:"id"`
		Type      string        `json:"type"`
		Data      []byte        `json:"data"`
		Position  geometry.Vec3 `json:"position"`
		Direction geometry.Vec3 `json:"direction"`
	}

	info.ID = entityInfo.ID
	info.Type = "object"
	info.Position = entityInfo.Position
	info.Direction = entityInfo.Direction

	payload, err := json.Marshal(info)
	if err != nil {
		return err
	}

	s.room.SendToAll("UpdateEntity", payload)
	return nil
}

func (s Sendler) SendEntityDelete(entityInfo ecs.EntityUpdateInfo) error {
	var info struct {
		ID types.Entity `json:"id"`
	}
	info.ID = entityInfo.ID
	payload, err := json.Marshal(info)
	if err != nil {
		return err
	}

	s.room.SendToAll("DeleteEntity", payload)
	return nil
}

func (s Sendler) CreateCameraForPlayer(id string, idEntity types.Entity) error {
	var info struct {
		ID   types.Entity `json:"id"`
		Type string       `json:"type"`
	}

	info.ID = idEntity
	info.Type = "camera"
	payload, err := json.Marshal(info)
	if err != nil {
		return err
	}

	var massange struct {
		cmd     string
		payload any
	}

	massange.cmd = "CreateEntity"
	massange.payload = payload

	j, _ := json.Marshal(massange)
	s.room.Clients[id].Send(j)
	return nil
}

func (s Sendler) SendSnapshotToAll(entities []ecs.EntityCreateInfo) error {
	// var info struct {
	// 	Entities []ecs.EntityCreateInfo `json:"entities"`
	// }

	// s.room.SendToAll("Snapshot", map[string]interface{}{
	// 	"entities": converted,
	// })
	return nil
}

// начало игры
type createGameSessionCommand struct{}

func (c *createGameSessionCommand) Name() string { return "createGameSession" }

func (c *createGameSessionCommand) Execute(client *Client, payload json.RawMessage) error {
	broadcaster := Sendler{room: client.room}
	w := game.CreateGame(broadcaster)

	for id := range client.room.Clients {
		transfer.CreatePlayer(w.EventBus, id)
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
