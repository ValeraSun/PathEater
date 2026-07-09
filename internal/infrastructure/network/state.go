package network

import (
    "context"
    "encoding/json"
    "errors"
)

type State interface {
    Name() string
    HandleCommand(ctx context.Context, client *Client, cmd string, payload json.RawMessage) error
    OnEnter(ctx context.Context, client *Client)
    OnExit(ctx context.Context, client *Client)
}

type GameState struct {
    commands map[string]Command
	GameID string
}

func NewGameState(gameID string) *GameState {
    s := &GameState{
        commands: make(map[string]Command),
        GameID:   gameID,
    }
    s.RegisterCommand(&MoveCommand{})
    s.RegisterCommand(&UseItemCommand{})
    s.RegisterCommand(&ExitCommand{})
    return s
}


func (s *GameState) Name() string { return "GAME" }

func (s *GameState) RegisterCommand(cmd Command) {
    s.commands[cmd.Name()] = cmd
}

func (s *GameState) HandleCommand(ctx context.Context, client *Client, cmdName string, payload json.RawMessage) error {
    cmd, exists := s.commands[cmdName]
    if !exists {
        return errors.New("unknown command in game")
    }
    return cmd.Execute(ctx, client, payload)
}

func (s *GameState) OnEnter(ctx context.Context, client *Client) {
    client.SendMessage("STATE_CHANGE", map[string]string{
        "state":  "GAME",
        "gameId": s.GameID,
    })
}

func (s *GameState) OnExit(ctx context.Context, client *Client) {
    //Выход из игры
}