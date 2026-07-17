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

	var data map[string]interface{}
	json.Unmarshal(jsonData, &data)

	rawColliders := data["colliders"].([]interface{})

	colliders := make([]wall, 0)
	for _, col := range rawColliders {
		c := col.(map[string]interface{})

		colliders = append(colliders,
			wall{
				shape:       c["shape"].(string),
				center:      c["center"].(geometry.Vec3),
				halfExtents: c["halfExtents"].(geometry.Vec3),
				rotation:    c["rotation"].(geometry.Quaternion).ToRotationMatrix(),
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

func NewWall() {

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
