package network

import (
    "context"
    "encoding/json"
)

type Command interface {
    Name() string
    Execute(ctx context.Context, client *Client, payload json.RawMessage) error
}

//------------------------------------------------------------------------------------------------------------------
//обработка команд начала
type StartGameCommand struct{}

func (c *StartGameCommand) Name() string { return "StartGame" }

func (c *StartGameCommand) Execute(ctx context.Context, client *Client, payload json.RawMessage) error {
    if client.state != nil {
        client.state.OnExit(ctx, client)
    }

    client.SetState(NewGameState("game"))

    if client.state != nil {
        client.state.OnEnter(ctx, client)
    }
    return nil
}


//------------------------------------------------------------------------------------------------------------------
//обработка команд движения
type MoveCommand struct{}

func (*MoveCommand) Name() string { return "Move" }

func(*MoveCommand) Execute(ctx context.Context, client *Client, payload json.RawMessage) error {	
	var coord struct {
		X float32 `json:"X"`
		Y float32 `json:"Y"`
		Z float32 `json:"Z"`
	}
    
	err := json.Unmarshal(payload, &coord)
	if err != nil {
		client.SendError("invalid_json")
	}
	return err
	//посылаем координаты на основную логику сервера
}


//------------------------------------------------------------------------------------------------------------------
//обработка команд взаимодействия с предметами
type UseItemCommand struct{}

func (*UseItemCommand) Name() string { return "UseItem" }

func(*UseItemCommand) Execute(ctx context.Context, client *Client, payload json.RawMessage) error {	
	var item struct {
	itemID    string `json:"itemID"`
	typeUsing string `json:"typeUsing"`
	}

    err := json.Unmarshal(payload, &item)
	if err != nil {
		client.SendError("invalid_json")
	}
	return err
	//посылаем данные на основную логику сервера
}


//------------------------------------------------------------------------------------------------------------------
//обработка команд выхода
type ExitCommand struct{}

func (*ExitCommand) Name() string { return "ExitGame" }

func (c *ExitCommand) Execute(ctx context.Context, client *Client, payload json.RawMessage) error {
    if client.state != nil {
        client.state.OnExit(ctx, client)
    }
    client.SetState(NewGameState("menu"))
    return nil
}

