package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const (
	config_path = "config/config.json"
)

type Config struct {
	Username string `json:"username"`
}

func (c *Config) SetUser(name string) error {
	c.Username = name

	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(config_path, bytes, 0666)
	if err != nil {
		return err
	}
	return nil
}

func ReadConfig() (Config, error) {
	cfg := Config{}

	file, err := os.Open(config_path)
	if err != nil {
		return Config{}, fmt.Errorf("os open err: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return Config{}, fmt.Errorf("io read all err: %w", err)
	}

	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("unmarshal err: %w", err)
	}

	return cfg, nil
}
