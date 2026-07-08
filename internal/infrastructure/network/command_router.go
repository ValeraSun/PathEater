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

type CommandMove struct {}

func (*CommandMove) Execute(ctx context.Context, client *ws.Client, payload json.RawMessage) error {

}

func (*CommandMove) Name() string {

}


type Req struct {
	Cmd     string          `json:"cmd"`
	Payload json.RawMessage `json:"payload"`
}

func (s *State) HandleCommand(ctx context.Context, client *ws.Client, cmdName string, payload json.RawMessage) error {
    cmd, exists := s.commands[cmdName]
    if !exists {
        return errors.New("unknown command in menu")
    }
    return cmd.Execute(ctx, client, payload)
}