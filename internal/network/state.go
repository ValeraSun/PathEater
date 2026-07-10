package network

import (
	"encoding/json"
	"fmt"
)

type State interface {
	Name() string
	HandleCommand(client *Client, cmd string, payload json.RawMessage) error
}

// BaseState содержит общую реализацию State.
type BaseState struct {
	name     string
	commands map[string]Command
}

func NewBaseState(name string) *BaseState {
	return &BaseState{
		name:     name,
		commands: make(map[string]Command),
	}
}

func (s *BaseState) Name() string { return s.name }

func (s *BaseState) registerCommand(cmd Command) {
	s.commands[cmd.Name()] = cmd
}

func (s *BaseState) HandleCommand(client *Client, cmdName string, payload json.RawMessage) error {
	cmd, exists := s.commands[cmdName]
	if !exists {
		return fmt.Errorf("неизвестная команда в состоянии %q", s.name)
	}
	return cmd.Execute(client, payload)
}

func MainMenuState() State {
	s := NewBaseState("mainMenu")
	s.RegisterCommand(&createRoomCommand{})
	s.RegisterCommand(&removeRoomCommand{})
	s.RegisterCommand(&joinRoomCommand{})
	s.RegisterCommand(&exitMenuCommand{})
	return s
}

func GameRoomState() State {
	s := NewBaseState("gameRoom")
	s.RegisterCommand(&removeRoomCommand{})
	s.RegisterCommand(&exitRoomCommand{})
	s.RegisterCommand(&startGameCommand{})
	return s
}

func GameState() State {
	s := NewBaseState("game")
	s.RegisterCommand(&movementCommand{})
	s.RegisterCommand(&useItemCommand{})
	s.RegisterCommand(&attackCommand{})
	s.RegisterCommand(&exitGameCommand{})
	return s
}