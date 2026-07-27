package systems

import (
	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const leak = 1.0

type zone struct {
	breakdowns  int
	totalOxygen float64
	rooms       map[types.Entity]*components.RoomComponent
}

type VacuumSystem struct {
	getter       componentsGetter
	rooms        map[types.Entity]*components.RoomComponent
	doors        map[types.Entity]*components.DoorComponent
	zones        []zone
	hasBreakdown bool
}

func NewVacuumSystem(getter componentsGetter, subscriber subscriber) *VacuumSystem {
	s := &VacuumSystem{
		getter:       getter,
		hasBreakdown: false,
	}
	s.rooms = s.getRooms()
	s.doors = s.getDoors()
	s.recalculateZones()
	subscriber.Subscribe("vacuum_recalculate", s.OnEvent)
	return s
}

func (s *VacuumSystem) getRooms() map[types.Entity]*components.RoomComponent {
	rooms := make(map[types.Entity]*components.RoomComponent)
	rs := s.getter.GetEntitiesByComponent("room")
	for id, room := range rs {
		rooms[id] = room.(*components.RoomComponent)
	}
	return rooms
}

func (s *VacuumSystem) getDoors() map[types.Entity]*components.DoorComponent {
	doors := make(map[types.Entity]*components.DoorComponent)
	ds := s.getter.GetEntitiesByComponent("door")
	for id, door := range ds {
		doors[id] = door.(*components.DoorComponent)
	}
	return doors
}

func (s *VacuumSystem) Update(dt float32) error {
	for i := range s.zones {
		zone := &s.zones[i]
		leakAmount := float64(zone.breakdowns) * leak * float64(dt)
		maxOxygen := float64(len(zone.rooms) * 100)

		if zone.breakdowns > 0 && zone.totalOxygen > 0 {
			zone.totalOxygen -= leakAmount
			if zone.totalOxygen < 0 {
				zone.totalOxygen = 0
			}
		}

		if !s.hasBreakdown && zone.totalOxygen < maxOxygen {
			zone.totalOxygen += leakAmount
			if zone.totalOxygen > maxOxygen {
				zone.totalOxygen = maxOxygen
			}
		}

		if len(zone.rooms) > 0 {
			localOxygen := zone.totalOxygen / float64(len(zone.rooms))
			for _, room := range zone.rooms {
				room.Oxygen = localOxygen
				room.Vacuum = (room.Oxygen < 0.001)
			}
		}
	}
	return nil
}

func (s *VacuumSystem) OnEvent(event events.Event) error {
	s.recalculateZones()
	return nil
}

func (s *VacuumSystem) recalculateZones() {
	graph := s.makeGraph()
	s.zones = s.findZones(graph)
}

func (s *VacuumSystem) makeGraph() map[types.Entity][]types.Entity {
	graph := make(map[types.Entity][]types.Entity)

	for _, door := range s.doors {
		if door.IsOpen {
			graph[door.RoomA] = append(graph[door.RoomA], door.RoomB)
			graph[door.RoomB] = append(graph[door.RoomB], door.RoomA)
		}
	}
	return graph
}

func (s *VacuumSystem) findZones(graph map[types.Entity][]types.Entity) []zone {
	s.hasBreakdown = false
	visited := make(map[types.Entity]bool)
	var zones []zone

	for roomID := range graph {
		zone := zone{
			breakdowns:  0,
			totalOxygen: 0,
			rooms:       make(map[types.Entity]*components.RoomComponent),
		}

		if visited[roomID] {
			continue
		}

		queue := []types.Entity{roomID}
		visited[roomID] = true

		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]
			zone.rooms[current] = s.rooms[current]
			if zone.rooms[current].HasBreakdown {
				s.hasBreakdown = true
				zone.breakdowns++
			}

			zone.totalOxygen = zone.totalOxygen + zone.rooms[current].Oxygen

			for _, neighbor := range graph[current] {
				if !visited[neighbor] {
					visited[neighbor] = true
					queue = append(queue, neighbor)
				}
			}
		}

		zones = append(zones, zone)
	}

	return zones
}
