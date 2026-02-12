package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	ServerURL string `json:"server_url"`
	Token     string `json:"token"`
	Role      string `json:"role"`
}

const ConfigFileName = "stolight_config.json"

func SaveConfig(token, role, url string) error {
	cfg := Config{ServerURL: url,
		Token: token,
		Role:  role,
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, ConfigFileName)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(cfg)

}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(home, ConfigFileName)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("not logged in. Please run 'sto login' first")
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
