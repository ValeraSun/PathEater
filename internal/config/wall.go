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

type wall struct {
	shape       string
	center      geometry.Vec3
	halfExtents geometry.Vec3
	rotation    [3]geometry.Vec3
}

func getWalls(path string) []wall {
	jsonData, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	type rootDTO struct {
		Colliders []struct {
			Shape       string              `json:"shape"`
			Center      geometry.Vec3       `json:"center"`
			HalfExtents geometry.Vec3       `json:"halfExtents"`
			Quaternion  geometry.Quaternion `json:"rotation"`
		} `json:"colliders"`
	}

	var DTO rootDTO
	err = json.Unmarshal(jsonData, &DTO)

	if err != nil {
		log.Fatal("Не распарсились стены")
	}

	collidersDTO := DTO.Colliders
	colliders := make([]wall, 0, len(collidersDTO))

	for _, col := range collidersDTO {
		colliders = append(colliders,
			wall{
				shape:       col.Shape,
				center:      col.Center,
				halfExtents: col.HalfExtents,
				rotation:    col.Quaternion.ToRotationMatrix(),
			})

	}
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
				wall.center,
				wall.halfExtents,
				wall.rotation)),
		)
	}
}
