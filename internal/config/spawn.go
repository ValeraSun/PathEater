package config

import (
	"log"
	"os"

	"github.com/ValeraSun/PathEater/internal/core/geometry"
	"gopkg.in/yaml.v3"
)

const avarageSpawns = 10
const defaultSpawnY = 2

func GetAlienSpawns() []geometry.Vec3 {
	type spawns struct {
		Alien []geometry.Vec3 `yaml:"alien_spawns"`
	}

	var config spawns

	data, err := os.ReadFile(CONFIG_PATH)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	for i := range config.Alien {
		config.Alien[i].Y = defaultSpawnY
	}

	return config.Alien

}
