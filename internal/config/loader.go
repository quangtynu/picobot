package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// LoadConfig loads config from ~/.picobot/config.json if present, otherwise returns defaults.
func LoadConfig() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	path := filepath.Join(home, ".picobot", "config.json")
	var cfg Config
	f, err := os.Open(path)
	if err != nil {
		// return empty config (not an error)
		return Config{}, nil
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
