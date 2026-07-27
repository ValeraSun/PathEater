package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ValeraSun/PathEater/internal/core/components"
	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"github.com/ValeraSun/PathEater/internal/core/types"
	"gopkg.in/yaml.v3"
)

type entityAdder interface {
	AddEntity(components ...types.Component) (types.Entity, error)
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

	colliders := DTO.Colliders

	return colliders
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

	colliders := DTO.Colliders

	return colliders
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

var ExternalWallEntities []types.Entity
var Rooms []types.Entity
var DoorsEntities []types.Entity

func CreateWalls(adder entityAdder) {
	path := getWallsConfigPath()
	walls := getWalls(path)
	for _, wall := range walls {
		w, err := adder.AddEntity(
			components.NewColliderComponent(geometry.NewBoxCollider(
				wall.Center,
				wall.HalfExtents,
				wall.Quaternion.ToRotationMatrix())),
		)
		if err != nil {
			log.Printf("Ошибка создания стены %s: %v", wall.ID, err)
			continue
		}
		if wall.Room != "" {
			r, err := adder.AddEntity(
				components.NewRoomComponent(),
			)
			if err != nil {
				log.Printf("Ошибка создания комнаты %s: %v", wall.ID, err)
				continue
			}
			NewRoom(r)
		}
		if wall.Type == "hull" {
			ExternalWallEntities = append(ExternalWallEntities, w)
		}
	}
}

func CreateDoors(adder entityAdder) {
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
			components.NewDoorComponent(types.Entity(door.RoomA), types.Entity(door.RoomB), tminx, tmaxx, tminy, tmaxy),
		)
		if err != nil {
			log.Printf("Ошибка создания двери %s: %v", door.ID, err)
			continue
		}
		DoorsEntities = append(DoorsEntities, d)
	}
}

func GetExternalWalls() []types.Entity {
	return ExternalWallEntities
}

func GetRooms() []types.Entity {
	return ExternalWallEntities
}

func GetDoors() []types.Entity {
	return DoorsEntities
}

func NewRoom(roomID types.Entity) {
	if contains(Rooms, roomID) {
		return
	}
	Rooms = append(Rooms, roomID)
}

func contains(slice []types.Entity, item types.Entity) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
