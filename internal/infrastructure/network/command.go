package network

import (
    "context"
    "encoding/json"
    "errors"
    "ws/internal/ws"
)

type Command interface {
    Name() string
    Execute(ctx context.Context, client *ws.Client, payload json.RawMessage) error
}

type StartGameCommand struct{}

func (c *StartGameCommand) Name() string { return "StartGame" }

func (c *StartGameCommand) Execute(ctx context.Context, client *ws.Client, payload json.RawMessage) error {
	
    client.SetState(NewGameState("game-123"))
    client.State.OnEnter(ctx, client)
    return nil
}

type MoveCommand struct{
	X, Y, Z float32
}

func (*MoveCommand) Name() string { return "Move"}

func(*MoveCommand) Execute(ctx context.Context, client Client, payload json.RawMessage) {	
	var coord struct {
		X string  `json:"X"`
		Y float32 `json:"Y"`
		Z float32 `json:"Z"`
	}
	if err := json.Unmarshal(message, &coord); err != nil {
		client.SendError("invalid_json")
		continue
	}
	//посылаем координаты на основную логику сервера
}
