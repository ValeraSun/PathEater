package network

import (
	"context"
	"encoding/json"
	"log"
)

type State interface {
	Name() string
	HandleCommand(ctx context.Context, client *Client, cmd string, payload json.RawMessage) error
	OnEnter(ctx context.Context, client *Client)
	OnExit(ctx context.Context, client *Client)
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

func (s *BaseState) HandleCommand(ctx context.Context, client *Client, cmdName string, payload json.RawMessage) error {
	cmd, exists := s.commands[cmdName]
	if !exists {
		log.Println("неизвестная команда в состоянии ", s.name, cmdName, payload)
		log.Println("Пришло", []byte(cmdName))
		crc := &createRoomCommand{}
		log.Println("Надо", []byte(crc.Name()))
		return nil
	}
	log.Println("команда прошла:", cmd.Name())
	return cmd.Execute(ctx, client, payload)
}

func (s *BaseState) OnEnter(ctx context.Context, client *Client) {}
func (s *BaseState) OnExit(ctx context.Context, client *Client)  {}

// Реализации State

func NewMainMenuState() State {
	s := NewBaseState("mainMenu")
	s.RegisterCommand(&createRoomCommand{})
	return s
}

func NewRoomMenuState() State {
	s := NewBaseState("roomMenu")
	s.RegisterCommand(&createGameSessionCommand{})
	return s
}

func NewPlayerControlState() State {
	s := NewBaseState("playerControl")
	s.RegisterCommand(&MoveCommand{})
	return s
}
