package configs

import (
	"Kudak/models"
	"encoding/json"
	"os"
)

type Config struct {
	Server   models.Server   `json:"server"`
	Database models.Database `json:"database"`
}

func InitConfigs() (*Config, error) {
	bytes, err := os.ReadFile("C:\\Users\\user\\Desktop\\work\\projects\\Kudak\\configs\\config.json")
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
