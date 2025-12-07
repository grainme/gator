package config

import (
	"encoding/json"
	"os"
)

const (
	configFileName = ".gatorconfig.json"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/" + configFileName, nil
}

func Read() (*Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var gatorConfig Config
	if err := json.Unmarshal(data, &gatorConfig); err != nil {
		return nil, err
	}

	return &gatorConfig, nil
}

func (cfg *Config) SetUser() error {
	// now that we have config struct already set
	// we need to update the json file.
	bytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	os.WriteFile(configFilePath, bytes, os.FileMode(0644))

	return nil
}
