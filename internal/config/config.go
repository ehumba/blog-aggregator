package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name,omitempty"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("no home directory: %v", err)
		return "", err
	}
	configFilePath := homeDir + "/" + configFileName
	return configFilePath, nil
}

func Read() (Config, error) {
	// Get and read config file
	configFile, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("could not get config file path: %v", err)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}, fmt.Errorf("could not read gatorconfig: %v", err)
	}

	// Decode data into Config struct
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("could not unmarshal data: %v", err)
	}
	return config, nil
}

func (c *Config) SetUser(username string) {
	c.CurrentUserName = username
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Printf("Could not marshal Config: %v", err)
		return
	}

	configFilePath, err := getConfigFilePath()
	if err != nil {
		fmt.Printf("could not get config file path: %v", err)
		return
	}
	// Why do I have to call it e?
	e := os.WriteFile(configFilePath, bytes, 0644)
	if e != nil {
		fmt.Printf("could not write to config file: %v", e)
	}
}
