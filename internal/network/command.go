package network

import (
	"encoding/json"
	"errors"
	"fmt"

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

var defaultSpawn = map[string]float64{"x": 2, "y": 2, "z": -3}

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

func (s Sendler) SendEntityCreate(entityInfo ecs.EntityInfo) error {
	var info struct {
		ID   types.Entity `json:"id"`
		Type string       `json:"type"`
		Data any          `json:"data"`
	}

	info.ID = entityInfo.ID
	info.Type = entityInfo.Type
	info.Data = entityInfo.Data

	payload, err := json.Marshal(info)
	if err != nil {
		return err
	}

	fmt.Printf("Отправлено create в room %v сообщение: %+v\n", s.room.ID, info)

	s.room.SendToAll("CreateEntity", payload)

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

func (s Sendler) SendEntityUpdate(entityInfo ecs.EntityInfo) error {
	var info struct {
		ID   types.Entity `json:"id"`
		Type string       `json:"type"`
		Data any          `json:"data"`
	}

	info.ID = entityInfo.ID
	info.Type = entityInfo.Type
	info.Data = entityInfo.Data

	payload, err := json.Marshal(info)
	if err != nil {
		return err
	}

	s.room.SendToAll("UpdateEntity", payload)

	//fmt.Printf("Отправлено update в room %v сообщение: %+v\n", s.room.ID, info)
	return nil
}

func (s Sendler) SendEntityDelete(entityInfo ecs.EntityInfo) error {
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

func (s Sendler) SendSnapshotToAll(entities []ecs.EntityInfo) error {
	// var info struct {
	// 	Entities []ecs.EntityCreateInfo `json:"entities"`
	// }

	// s.room.SendToAll("Snapshot", map[string]interface{}{
	// 	"entities": converted,
	// })
	return nil
}

func (s Sendler) SendGameOverState(gameOverInfo ecs.GameOverInfo) error {
	var info struct {
		Data any `json:"data"`
	}
	info.Data = gameOverInfo.Data
	payload, err := json.Marshal(info)
	if err != nil {
		return err
	}

	s.room.SendToAll("GameOver", payload)
	return nil
}

// начало игры
type createGameSessionCommand struct{}

func (c *createGameSessionCommand) Name() string { return "createGameSession" }

func (c *createGameSessionCommand) Execute(client *Client, payload json.RawMessage) error {
	broadcaster := Sendler{room: client.room}
	w := game.CreateGame(broadcaster)
	client.room.World = w

	for id, roomClient := range client.room.Clients {
		transfer.CreatePlayer(w.EventBus, id)

		roomClient.SetState(PlayerControlState())

		roomClient.SendMessage("GameStarted", map[string]interface{}{
			"playerId": id,
			"spawn":    defaultSpawn,
		})
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

	transfer.SendPlayerState(client.room.World.EventBus, state, client.ID)

	return nil
}

type weaponStateCommand struct{}

func (c *weaponStateCommand) Name() string { return "weaponState" }

func (c *weaponStateCommand) Execute(client *Client, payload json.RawMessage) error {
	var state events.WeaponState

	if err := json.Unmarshal(payload, &state); err != nil {
		return client.SendError(err)
	}

	transfer.SendWeaponState(client.room.World.EventBus, state)

	return nil
}

type shipStateCommand struct{}

func (c *shipStateCommand) Name() string { return "shipState" }

func (c *shipStateCommand) Execute(client *Client, payload json.RawMessage) error {
	var state events.ShipState

	if err := json.Unmarshal(payload, &state); err != nil {
		return client.SendError(err)
	}

	transfer.SendShipState(client.room.World.EventBus, state)

	return nil
}

type exitGameCommand struct{}

func (c *exitGameCommand) Name() string { return "exitGame" }

func (c *exitGameCommand) Execute(client *Client, payload json.RawMessage) error {
	if client.room != nil {
		client.SetState(GameRoomState())
	} else {
		client.SetState(MainMenuState())
	}
	return nil
}
