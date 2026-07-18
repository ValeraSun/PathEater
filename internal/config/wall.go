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

func getWallsConfigPath() string {
	jsonData, err := os.ReadFile(CONFIG_PATH)

	if err != nil {
		log.Fatal(err)
	}

	var data map[string]interface{}
	yaml.Unmarshal(jsonData, &data)

	path := data["walls_path"].(string)

	return path
}

func CreateWalls(adder entityAdder) {
	path := getWallsConfigPath()
	walls := getWalls(path)
	for _, wall := range walls {
		adder.AddEntity(
			components.NewColliderComponent(geometry.NewBoxCollider(
				wall.Center,
				wall.HalfExtents,
				wall.Quaternion)),
		)
	}
}
