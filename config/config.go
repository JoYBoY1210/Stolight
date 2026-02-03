package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	ServerPort   int      `json:"server_port"`
	StorageNodes []string `json:"storage_nodes"`
	DBPath       string   `json:"db_path"`
}

// global var as i want to load conifg only once but want to use it in every API
var Cfg *Config

func LoadConfig() *Config {
	var cfg Config
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Println("config file not found using default values")
		return &Config{
			ServerPort:   6124,
			StorageNodes: []string{"./data_nodes/node1", "./data_nodes/node2", "./data_nodes/node3"},
			DBPath:       "stolight.db",
		}
	}
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatalf("failed to parse config file: %v", err)
		return nil
	}
	fmt.Println("config loaded successfully")
	return &cfg
}

func init() {
	Cfg = LoadConfig()
}
