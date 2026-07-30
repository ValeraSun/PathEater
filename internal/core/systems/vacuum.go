package systems

import (
	"sort"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/events"
	"github.com/ValeraSun/PathEater/internal/core/types"
)

const leak = 25.0

type zone struct {
	breakdowns  int
	totalOxygen float64
	rooms       map[string]*components.RoomComponent
}

type VacuumSystem struct {
	getter       componentsGetter
	rooms        map[string]*components.RoomComponent
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

func (s *VacuumSystem) getRooms() map[string]*components.RoomComponent {
	rooms := make(map[string]*components.RoomComponent)
	rs := s.getter.GetEntitiesByComponent("room")
	for _, room := range rs {
		r, ok := room.(*components.RoomComponent)
		if ok && r != nil {
			rooms[r.Name] = r
		}
	}
	return rooms
}

func (s *VacuumSystem) getDoors() map[types.Entity]*components.DoorComponent {
	doors := make(map[types.Entity]*components.DoorComponent)
	ds := s.getter.GetEntitiesByComponent("door")
	for id, door := range ds {
		d, ok := door.(*components.DoorComponent)
		if ok && d != nil {
			doors[id] = d
		}
	}
	return doors
}

func (s *VacuumSystem) Update(dt float32) error {

	s.recalculateZones()

	for i := range s.zones {
		zone := &s.zones[i]
		leakAmount := float64(zone.breakdowns) * leak * float64(dt)
		maxOxygen := float64(len(zone.rooms) * 100)

		// Утечка, если есть поломки
		if zone.breakdowns > 0 && zone.totalOxygen > 0 {
			zone.totalOxygen -= leakAmount
			if zone.totalOxygen < 0 {
				zone.totalOxygen = 0
			}
		}

		// Восстановление, если поломок нет
		if zone.breakdowns == 0 && zone.totalOxygen < maxOxygen {
			zone.totalOxygen += leakAmount
			if zone.totalOxygen > maxOxygen {
				zone.totalOxygen = maxOxygen
			}
		}

		// Распределение кислорода по комнатам
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

func (s *VacuumSystem) makeGraph() map[string][]string {
	graph := make(map[string][]string)

	openDoors := 0
	for _, door := range s.doors {
		if door.IsOpen {
			openDoors++
			graph[door.RoomA] = append(graph[door.RoomA], door.RoomB)
			graph[door.RoomB] = append(graph[door.RoomB], door.RoomA)
		}
	}

	return graph
}

func (s *VacuumSystem) findZones(graph map[string][]string) []zone {
	s.hasBreakdown = false
	visited := make(map[string]bool)
	var zones []zone

	for roomID := range graph {
		if visited[roomID] {
			continue
		}

		zone := zone{
			breakdowns:  0,
			totalOxygen: 0,
			rooms:       make(map[string]*components.RoomComponent),
		}

		queue := []string{roomID}
		visited[roomID] = true
		zoneSize := 0

		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]

			room, exists := s.rooms[current]
			if !exists || room == nil {
				continue
			}

			zone.rooms[current] = room
			zoneSize++
			if room.HasBreakdown {
				s.hasBreakdown = true
				zone.breakdowns++
			}
			zone.totalOxygen += room.Oxygen

			for _, neighbor := range graph[current] {
				if !visited[neighbor] {
					visited[neighbor] = true
					queue = append(queue, neighbor)
				}
			}
		}

		if len(zone.rooms) > 0 {
			zones = append(zones, zone)
			var roomNames []string
			for name := range zone.rooms {
				roomNames = append(roomNames, name)
			}
			sort.Strings(roomNames)
		}
	}

	isolatedCount := 0
	for roomName := range s.rooms {
		if !visited[roomName] {
			room, exists := s.rooms[roomName]
			if !exists || room == nil {
				continue
			}

			isolatedCount++
			zone := zone{
				breakdowns:  0,
				totalOxygen: 0,
				rooms:       make(map[string]*components.RoomComponent),
			}

			zone.rooms[roomName] = room
			if room.HasBreakdown {
				s.hasBreakdown = true
				zone.breakdowns++
			}
			zone.totalOxygen = room.Oxygen

			zones = append(zones, zone)
		}
	}

	return zones
}
