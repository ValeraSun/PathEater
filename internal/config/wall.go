package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/sendler"
	"github.com/ValeraSun/PathEater/internal/core/types"
	"gopkg.in/yaml.v3"
)

type entityAdder interface {
	AddEntity(components ...types.Component) (types.Entity, error)
}

type componentsGetter interface {
	GetEntitiesByComponent(componentType string) map[types.Entity]types.Component
	HasComponents(entity types.Entity, componentTypes ...string) bool
	GetComponent(entity types.Entity, componentType string) (types.Component, bool)
}

type entitySendler interface {
	Send(types.Entity, func(sendler.EntityInfo) error) error
}

type broadcasterFunc interface {
	SendEntityCreate(entityInfo sendler.EntityInfo) error
	SendEntityUpdate(entityInfo sendler.EntityInfo) error
	SendEntityDelete(entityInfo sendler.EntityInfo) error
	SendGameOverState(sendler.GameOverInfo) error
	SendTime(sendler.TimeInfo) error
}

type broadcaster interface {
	entitySendler
	broadcasterFunc
}

type Wall struct {
	ID          string              `json:"id"`
	Room        string              `json:"room"`
	Type        string              `json:"type"`
	Shape       string              `json:"shape"`
	Center      geometry.Vec3       `json:"center"`
	HalfExtents geometry.Vec3       `json:"halfExtents"`
	Quaternion  geometry.Quaternion `json:"rotation"`
}

type Door struct {
	ID          string              `json:"id"`
	RoomA       string              `json:"roomA"`
	RoomB       string              `json:"roomB"`
	Shape       string              `json:"shape"`
	Center      geometry.Vec3       `json:"center"`
	HalfExtents geometry.Vec3       `json:"halfExtents"`
	Quaternion  geometry.Quaternion `json:"rotation"`
}

type Room struct {
	ID   string  `json:"id"`
	MinX float64 `json:"minX"`
	MaxX float64 `json:"maxX"`
	MinZ float64 `json:"minZ"`
	MaxZ float64 `json:"maxZ"`
}

func getWalls(path string) []Wall {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	type rootDTO struct {
		Colliders []Wall `json:"colliders"`
	}

	var DTO rootDTO
	err = json.Unmarshal(jsonData, &DTO)
	if err != nil {
		log.Fatal("Не распарсились стены")
	}

	return DTO.Colliders
}

func getDoors(path string) []Door {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	type rootDTO struct {
		Colliders []Door `json:"colliders"`
	}

	var DTO rootDTO
	err = json.Unmarshal(jsonData, &DTO)
	if err != nil {
		log.Fatal("Не распарсились двери")
	}

	return DTO.Colliders
}

func getRooms(path string) []Room {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	type rootDTO struct {
		Rooms []Room `json:"rooms"`
	}

	var DTO rootDTO
	err = json.Unmarshal(jsonData, &DTO)
	if err != nil {
		log.Fatal("Не распарсились комнаты")
	}

	return DTO.Rooms
}

func getWallsConfigPath() string {
	jsonData, err := os.ReadFile(CONFIG_PATH)
	if err != nil {
		log.Fatal(err)
	}

	var data map[string]interface{}
	yaml.Unmarshal(jsonData, &data)

	path := data["walls_path_dev"].(string)
	_, err = os.ReadFile(path)
	if err != nil {
		path = data["walls_path"].(string)
	}

	return path
}

func getDoorsConfigPath() string {
	jsonData, err := os.ReadFile(CONFIG_PATH)
	if err != nil {
		log.Fatalf("Cannot read config file: %v", err)
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal(jsonData, &data); err != nil {
		log.Fatalf("Cannot parse config file: %v", err)
	}

	devPath, ok := data["doors_path_dev"].(string)
	if !ok || devPath == "" {
		log.Printf("doors_path_dev not found or empty in config")
	} else {
		log.Printf("Checking dev path: %s", devPath)
		if _, err := os.ReadFile(devPath); err == nil {
			log.Printf("Using dev path: %s", devPath)
			return devPath
		}
		log.Printf("Dev path failed: %v", err)
	}

	prodPath, ok := data["doors_path"].(string)
	if !ok || prodPath == "" {
		log.Fatal("doors_path not found in config")
	}

	log.Printf("Using production path: %s", prodPath)
	return prodPath
}

func getRoomsConfigPath() string {
	jsonData, err := os.ReadFile(CONFIG_PATH)
	if err != nil {
		log.Fatalf("Cannot read config file: %v", err)
	}

	var data map[string]interface{}
	if err := yaml.Unmarshal(jsonData, &data); err != nil {
		log.Fatalf("Cannot parse config file: %v", err)
	}

	devPath, ok := data["rooms_path_dev"].(string)
	if !ok || devPath == "" {
		log.Printf("rooms_path_dev not found or empty in config")
	} else {
		log.Printf("Checking dev path: %s", devPath)
		if _, err := os.ReadFile(devPath); err == nil {
			log.Printf("Using dev path: %s", devPath)
			return devPath
		}
		log.Printf("Dev path failed: %v", err)
	}

	prodPath, ok := data["rooms_path"].(string)
	if !ok || prodPath == "" {
		log.Fatal("rooms_path not found in config")
	}

	log.Printf("Using production path: %s", prodPath)
	return prodPath
}

func CreateWorldColliders(adder entityAdder, getter componentsGetter, broadcaster broadcaster) *WorldStructures {
	w := newWorldStructs()
	w.createWalls(adder)
	w.createRooms(adder, getter)
	w.createDoors(adder, broadcaster)
	return w
}

type WorldStructures struct {
	ExternalWallEntities []types.Entity
	ExternalWalls        map[types.Entity]Wall
	Rooms                map[types.Entity]types.Entity
	DoorsEntities        []types.Entity
	RoomEntities         []types.Entity
}

func newWorldStructs() *WorldStructures {
	return &WorldStructures{
		ExternalWalls: make(map[types.Entity]Wall),
		Rooms:         make(map[types.Entity]types.Entity),
	}
}

func (world *WorldStructures) createWalls(adder entityAdder) {
	path := getWallsConfigPath()
	walls := getWalls(path)

	for _, wall := range walls {
		w, err := adder.AddEntity(
			components.NewColliderComponent(geometry.NewBoxCollider(
				wall.Center,
				wall.HalfExtents,
				wall.Quaternion.ToRotationMatrix())),
			components.NewWallComponent(),
		)
		if err != nil {
			log.Printf("Ошибка создания стены %s: %v", wall.ID, err)
			continue
		}

		if wall.Type == "hull" {
			world.ExternalWalls[w] = wall
			world.ExternalWallEntities = append(world.ExternalWallEntities, w)
		}
	}
}

func (world *WorldStructures) createRooms(adder entityAdder, getter componentsGetter) {
	path := getRoomsConfigPath()
	rooms := getRooms(path)

	for _, room := range rooms {
		r, err := adder.AddEntity(
			components.NewRoomComponent(room.ID, room.MinX, room.MaxX, room.MinZ, room.MaxZ),
		)

		if err != nil {
			log.Printf("Ошибка создания комнаты %s: %v", room.ID, err)
			continue
		}

		for w, wall := range world.ExternalWalls {
			if wall.Room == room.ID {
				world.Rooms[w] = r
				c, _ := getter.GetComponent(w, "wall")
				wallComp := c.(*components.WallComponent)
				wallComp.Room = r
			}
		}
		world.RoomEntities = append(world.RoomEntities, r)
	}

	log.Printf("CreateRooms: created %d rooms", len(world.RoomEntities))
}

func (world *WorldStructures) createDoors(adder entityAdder, broadcaster broadcaster) {
	path := getDoorsConfigPath()
	doors := getDoors(path)
	for _, door := range doors {
		tminx := door.Center.X - door.HalfExtents.X
		tmaxx := door.Center.X + door.HalfExtents.X
		tminy := door.Center.Y - door.HalfExtents.Y
		tmaxy := door.Center.Y + door.HalfExtents.Y
		d, err := adder.AddEntity(
			components.NewColliderComponent(geometry.NewBoxCollider(
				door.Center,
				door.HalfExtents,
				door.Quaternion.ToRotationMatrix())),
			components.NewDoorComponent(door.RoomA, door.RoomB, tminx, tmaxx, tminy, tmaxy),
			components.NewUpdateComponent(),
			components.NewInteractableComponent(door.Center, door.HalfExtents, "down", components.InteractDoor),
		)
		if err != nil {
			log.Printf("Ошибка создания двери %s: %v", door.ID, err)
			continue
		}
		world.DoorsEntities = append(world.DoorsEntities, d)
		err = broadcaster.Send(d, broadcaster.SendEntityCreate)
	}
}

func (world *WorldStructures) GetExternalWalls() []types.Entity {
	return world.ExternalWallEntities
}

func (world *WorldStructures) GetRooms() map[types.Entity]types.Entity {
	return world.Rooms
}
