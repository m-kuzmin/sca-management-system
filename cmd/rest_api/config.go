package main

import (
	"os"

	"github.com/pelletier/go-toml"
)

// Config represents the server configuration
type Config struct {
	Database DatabaseConfig `toml:"database"`
	Server   ServerConfig   `toml:"server"`
}

// DatabaseConfig holds database related configuration
type DatabaseConfig struct {
	Driver           string `toml:"driver"`
	Address          string `toml:"address"`
	DBName           string `toml:"dbName"`
	PingRetries      uint   `toml:"pingRetries"`
	PingIntervalSecs uint   `toml:"pingIntervalSecs"`
	MigrationsSource string `toml:"migrationsSource"`
}

// ServerConfig holds server related configuration
type ServerConfig struct {
	Port          uint16    `toml:"port"`
	ReadTimeoutMs uint      `toml:"readTimeoutMs"`
	Gin           GinConfig `toml:"gin"`
}

type GinConfig struct {
	Mode string `toml:"mode"`
}

// NewConfigFromFile loads configuration from a TOML file
func NewConfigFromFile(filename string) (Config, error) {
	configText, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err = toml.Unmarshal(configText, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
