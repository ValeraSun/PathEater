package network

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
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
    
    //действия на сервере

    // Меняем состояние
    client.SetState(GameRoomState())
    return nil
}

//удаление комнаты
type removeRoomCommand struct {}

func (c *removeRoomCommand) Name() string { return "removeRoom" }

func (c *removeRoomCommand) Execute(client *Client, payload json.RawMessage) error {
    var roomID struct {
		RoomID string `json:"roomID"`
	}
	if err := json.Unmarshal(payload, &roomID); err != nil {
		return client.SendError(err)
	}
    hub := GetHub()
    room, exists := hub.GetGameRoom(roomID)
    if !exists {
        return fmt.Errorf("Комнаты с таким ID не существует")
    }
    room.Close()
    
    //действия на сервере

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
        return fmt.Errorf("Комнаты с таким ID не существует")
    }
    room.AddClient(client)
    
    //действия на сервере

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
    
    //действия на сервере

    return nil
}


//-----------------------------------------------------------------------------------------------------------------------------------
//Команды комнат

//Выход из комнаты
type exitRoomCommand struct {}

func (c *exitRoomCommand) Name() string { return "exitRoom" }

func (c *exitRoomCommand) Execute(client *Client, payload json.RawMessage) error {
    var roomID struct {
		RoomID string `json:"roomID"`
	}
	if err := json.Unmarshal(payload, &roomID); err != nil {
		return client.SendError(err)
	}
    hub := GetHub()
    room := hub.GetGameRoom(roomID)
    room.RemoveClient(client)
    
    //действия на сервере

    return nil
}

//начало игры
type startGameCommand struct {}

func (c *startGameCommand) Name() string { return "startGame" }

func (c *startGameCommand) Execute(client *Client, payload json.RawMessage) error {
    
    //действия на сервере

    return nil
}

//-----------------------------------------------------------------------------------------------------------------------------------
//Команды игры

//движение
type movementCommand struct {}

func (c *movementCommand) Name() string { return "movement" }

func (c *movementCommand) Execute(client *Client, payload json.RawMessage) error {
    var coord struct {
		X float32 `json:"X"`
		Y float32 `json:"Y"`
		Z float32 `json:"Z"`
	}
	if err := json.Unmarshal(payload, &coord); err != nil {
		return client.SendError(err)
	}
    
    //действия на сервере

    e := events.NewMoveEvent(coord)
    client.room.world.EventBus.Publish(e)

    return nil
}

//использование предмета
type useItemCommand struct {}

func (c *useItemCommand) Name() string { return "useItem" }

func (c *useItemCommand) Execute(client *Client, payload json.RawMessage) error {
    var itemID struct {
		ItemID float32 `json:"itemID"`
	}
	if err := json.Unmarshal(payload, &itemID); err != nil {
		return client.SendError(err)
	}
    
    //действия на сервере

    return nil
}

//атака ближнего боя
type attackCommand struct {}

func (c *attackCommand) Name() string { return "attack" }

func (c *attackCommand) Execute(client *Client, payload json.RawMessage) error {
    
    //действия на сервере

    return nil
}

//выход в комнату (в лобби)
type exitGameCommand struct {}

func (c *exitGameCommand) Name() string { return "exitGame" }

func (c *exitGameCommand) Execute(client *Client, payload json.RawMessage) error {
    
    //действия на сервере

    // Меняем состояние
    client.SetState(GameRoomState())
    return nil
}
