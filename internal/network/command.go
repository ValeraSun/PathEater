package network

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"

	"github.com/ValeraSun/PathEater/internal/core/events"
)

type Command interface {
    Name() string
    Execute(client *Client, payload json.RawMessage) error
}


//-----------------------------------------------------------------------------------------------------------------------------------
//Команды меню игры

//создание комнаты
type createRoomCommand struct {}

func (c *createRoomCommand) Name() string { return "createRoom" }

func (c *createRoomCommand) Execute(client *Client, payload json.RawMessage) error {
    // Создаём комнату с уникальным ID
    roomID := GenerateID()
    hub := GetHub()
    room := hub.NewGameRoom(roomID)
    room.AddClient(client)

    var info struct{
        ID   EntityInfo.Id   `json:"id"`
    }
	if payload, err := json.Marshal(info); err != nil {
		return err
	}

    s.room.SendToAll("SuccessCreateRoom", payload)

    // Меняем состояние
    client.SetState(GameRoomState())
    return nil
}

//удаление комнаты
type deleteRoomCommand struct {}

func (c *deleteRoomCommand) Name() string { return "deleteRoom" }

func (c *deleteRoomCommand) Execute(client *Client, payload json.RawMessage) error {
    var roomID struct {
		RoomID string `json:"roomID"`
	}
	if err := json.Unmarshal(payload, &roomID); err != nil {
		return client.SendError(err)
	}
    hub := GetHub()

    room, exists := hub.GetGameRoom(roomID)
    if !exists {
        return client.SendText("error", "RoomIsNotExists")
    }
    room.Close()

    var info struct{
        ID   EntityInfo.Id   `json:"id"`
    }
	if payload, err := json.Marshal(info); err != nil {
		return err
	}

    s.room.SendToAll("SuccessDeleteRoom", payload)

    // Меняем состояние
    if client.state == GameRoomState() {
      client.SetState(MainMenuState())
    } 
    return nil
}

//присоединение к комнате
type joinRoomCommand struct {}

func (c *joinRoomCommand) Name() string { return "joinRoom" }

func (c *joinRoomCommand) Execute(client *Client, payload json.RawMessage) error {
    var roomID struct {
		RoomID string `json:"roomID"`
	}
	if err := json.Unmarshal(payload, &roomID); err != nil {
		return client.SendError(err)
	}
    hub := GetHub()
    room, exists := hub.GetGameRoom(roomID)
    if !exists {
        return client.SendText("error", "RoomIsNotExists")
    }
    room.AddClient(client)

    var info struct{
        ID   EntityInfo.Id   `json:"id"`
    }
	if payload, err := json.Marshal(info); err != nil {
		return err
	}

    s.room.SendToAll("SuccessJoinRoom", payload)

    // Меняем состояние
    client.SetState(GameRoomState())
    return nil
}

//выход из меню (выход из всей игры) - отключение клиента
type exitMenuCommand struct {}

func (c *exitMenuCommand) Name() string { return "exitMenu" }

func (c *exitMenuCommand) Execute(client *Client, payload json.RawMessage) error {
    hub := GetHub()
    hub.UnregisterClient(client)
    client.Close()

    return nil
}


//-----------------------------------------------------------------------------------------------------------------------------------
//Команды комнат

//Выход из комнаты
type exitRoomCommand struct {}

func (c *exitRoomCommand) Name() string { return "exitRoom" }

func (c *exitRoomCommand) Execute(client *Client, payload json.RawMessage) error {
    hub := GetHub()
    client.room.RemoveClient(client)

    return nil
}

type Sendler struct {
    room GameRoom
}

func (s Sendler) SendEntityCreate(EntityInfo) error {
    var info struct{
        ID   EntityIndo.Id   `json:"id"`
	    Type EntityIndo.Type `json:"type"`
	    Data EntityIndo.Data `json:"data"`
    }
	if payload, err := json.Marshal(info); err != nil {
		return err
	}

    s.room.SendToAll("CreateEntity", payload)
}

func (s Sendler) SendEntityUpdate(EntityInfo) error {
    var info struct{
        ID   EntityInfo.Id   `json:"id"`
	    Type EntityInfo.Type `json:"type"`
	    Data EntityInfo.Data `json:"data"`
    }
	if payload, err := json.Marshal(info); err != nil {
		return err
	}

    s.room.SendToAll("UpdateEntity", payload)
}

func (s Sendler) SendEntityDelete(EntityInfo) error {
    var info struct{
        ID   EntityInfo.Id   `json:"id"`
    }
	if payload, err := json.Marshal(info); err != nil {
		return err
	}

    s.room.SendToAll("DeleteEntity", payload)
}

func (s Sendler) SendSnapshot(types.Entity, []EntityInfo) error {
    var info struct{
        ID   EntityInfo.Id   `json:"id"`
	    Type EntityInfo.Type `json:"type"`
	    Data EntityInfo.Data `json:"data"`
    }
	if payload, err := json.Marshal(info); err != nil {
		return err
	}

    s.room.SendToAll("Snapshot", payload)
}

//начало игры
type startGameCommand struct {}

func (c *startGameCommand) Name() string { return "startGame" }

func (c *startGameCommand) Execute(client *Client, payload json.RawMessage) error {
    
    CreateWorld(Sendler)

    return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------
//Команды игры

type playerStateCommand struct {}

func (c *playerStateCommand) Name() string { return "startGame" }

func (c *playerStateCommand) Execute(client *Client, payload json.RawMessage) error {
    var state events.PlayerState

    if err := json.Unmarshal(payload, &state); err != nil {
		return client.SendError(err)
	}
    
    events.SendPlayerState(client.room.World.EventBus, state, client.ID)

    return nil
}