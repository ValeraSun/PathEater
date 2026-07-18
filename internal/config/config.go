package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const CONFIG_PATH = "/app/internal/config/config.yaml"

type Config struct {
	Server ServerConfig `yaml:"server" json:"server"`
}

type ServerConfig struct {
	World WorldConfig `yaml:"world" json:"world"`
}

type WorldConfig struct {
	TickPerSecond int `yaml:"tick_per_second" json:"tick_per_second"`
}

func getTickPerSecond() int {
	data, err := os.ReadFile(CONFIG_PATH)
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	return cfg.Server.World.TickPerSecond
}

func GetNanosecondPerTick() time.Duration {
	tickPerSecond := getTickPerSecond()
	secondPerTick := 1.0 / float32(tickPerSecond)
	answer := time.Duration(secondPerTick * float32(time.Second))
	return answer
}
