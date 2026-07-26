package systems

import (
	"sync"

	"github.com/ValeraSun/PathEater/internal/core/events"
)

type RoomState struct {
	Pressure       float64
	BreakdownCount int
	IsSealed       bool
	mu             sync.Mutex
}

type VacuumSystem struct {
	rooms      map[string]*RoomState
	getter     componentsGetter
	eventQueue chan events.Event
}

func NewVacuumSystem(getter componentsGetter, subscriber subscriber) *VacuumSystem {
	s := &VacuumSystem{
		rooms:      make(map[string]*RoomState),
		getter:     getter,
		eventQueue: make(chan events.Event, 100),
	}
	subscriber.Subscribe("breach_created", s.OnBreachCreated)
	subscriber.Subscribe("breach_repaired", s.OnBreachRepaired)
	return s
}
func (s *VacuumSystem) Update(dt float64) {
	// Обрабатываем события из очереди
	for {
		select {
		case e := <-s.eventQueue:
			switch ev := e.(type) {
			case *events.BreachCreatedEvent:
				s.addBreach(ev.RoomID)
			case *events.BreachRepairedEvent:
				s.removeBreach(ev.RoomID)
			}
		default:
			s.updatePressure()
		}
	}
}

const LeakageRate = 0.05

func (s *VacuumSystem) updatePressure(dt float64) {
	for roomID, state := range s.rooms {
		state.mu.Lock()
		if state.BreakdownCount > 0 && !state.IsSealed {
			state.Pressure -= LeakageRate * dt
			if state.Pressure < 0 {
				state.Pressure = 0
			}
			if state.Pressure == 0 {
				//s.publisher.Publish(events.NewVacuumEvent(roomID))
			}
		}
		state.mu.Unlock()
	}
}

func (s *VacuumSystem) addBreach(roomID string) {
	s.getOrCreateRoom(roomID)
	s.rooms[roomID].mu.Lock()
	s.rooms[roomID].BreakdownCount++
	s.rooms[roomID].mu.Unlock()
}

func (vs *VacuumSystem) removeBreach(roomID string) {
	if _, ok := vs.rooms[roomID]; !ok {
		return
	}
	vs.rooms[roomID].mu.Lock()
	if vs.rooms[roomID].BreakdownCount > 0 {
		vs.rooms[roomID].BreakdownCount--
	}
	vs.rooms[roomID].mu.Unlock()
}

func (vs *VacuumSystem) getOrCreateRoom(roomID string) *RoomState {
	vs.rooms[roomID].mu.Lock()
	defer vs.rooms[roomID].mu.Unlock()
	if _, ok := vs.rooms[roomID]; !ok {
		vs.rooms[roomID] = &RoomState{
			Pressure:       1.0,
			BreakdownCount: 0,
			IsSealed:       false,
		}
	}
	return vs.rooms[roomID]
}

func (s *VacuumSystem) OnBreachCreated(event events.Event) error {
	ev := event.(*events.BreachCreatedEvent)
	s.eventQueue <- ev
	return nil
}

func (s *VacuumSystem) OnBreachRepaired(event events.Event) error {
	ev := event.(*events.BreachRepairedEvent)
	s.eventQueue <- ev
	return nil
}
