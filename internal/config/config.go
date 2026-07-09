package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

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
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	return cfg.Server.World.TickPerSecond
}

func GetMillisecondPerTick() time.Duration {
	tickPerSecond := getTickPerSecond()
	millisecondPerTick := 1000 / tickPerSecond

	return time.Duration(millisecondPerTick)
}
