package network

import (
    "context"
    "encoding/json"
    "errors"
    "ws/internal/ws"
)

type State interface {
    Name() string
    HandleCommand(ctx context.Context, client *ws.Client, cmd string, payload json.RawMessage) error
    OnEnter(ctx context.Context, client *ws.Client)
    OnExit(ctx context.Context, client *ws.Client)
}

func NewGameState(gameID string) *State {
    s := &State{
        commands: make(map[string]Command),
        GameID:   gameID,
    }
    s.RegisterCommand(&MoveCommand{})
    s.RegisterCommand(&UseItemCommand{})
    s.RegisterCommand(&ExitCommand{})
    return s
}

type GameState struct {
    commands map[string]Command
	GameID string
}

func (s *GameState) Name() string { return "GAME" }

func (s *GameState) RegisterCommand(cmd Command) {
    s.commands[cmd.Name()] = cmd
}

func (s *GameState) HandleCommand(ctx context.Context, client Client, cmdName string, payload json.RawMessage) error {
    cmd, exists := s.commands[cmdName]
    if !exists {
        return errors.New("unknown command in game")
    }
    return cmd.Execute(ctx, client, payload)
}

func (s *GameState) OnEnter(ctx context.Context, client *ws.Client) {
    client.SendMessage("STATE_CHANGE", map[string]string{
        "state":  "GAME",
        "gameId": s.GameID,
    })
}

func (s *GameState) OnExit(ctx context.Context, client *ws.Client) {
    //Выход из игры
}