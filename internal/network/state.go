package network

import (
	"encoding/json"
	"log"
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

func (s *BaseState) RegisterCommand(cmd Command) {
	s.commands[cmd.Name()] = cmd
}

func (s *BaseState) HandleCommand(client *Client, cmdName string, payload json.RawMessage) error {
	cmd, exists := s.commands[cmdName]
	if !exists {
		log.Println("неизвестная команда в состоянии ", s.name, cmdName, payload)
		return nil
	}
	return cmd.Execute(client, payload)
}

// Реализации State

func MainMenuState() State {
	s := NewBaseState("mainMenu")
	s.RegisterCommand(&createRoomCommand{})
	s.RegisterCommand(&deleteRoomCommand{})
	s.RegisterCommand(&joinRoomCommand{})
	s.RegisterCommand(&exitMenuCommand{})
	return s
}

func GameRoomState() State {
	s := NewBaseState("roomMenu")
	s.RegisterCommand(&deleteRoomCommand{})
	s.RegisterCommand(&exitRoomCommand{})
	s.RegisterCommand(&startGameCommand{})
	return s
}

func PlayerControlState() State {
	s := NewBaseState("playerControl")
	s.RegisterCommand(&playerStateCommand{})
	s.RegisterCommand(&exitGameCommand{})
	return s
}
