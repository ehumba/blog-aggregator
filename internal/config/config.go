package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DB_URL          string `json: db_url`
	CurrentUserName string `json: current_user_name,omitempty`
}

func Read() (Config, error) {
	// Get and read the gatorconfig file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Print("no home directory")
		return Config{}, err
	}
	configFile := homeDir + "/" + "gatorconfig.json"

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

func (c Config) SetUser(username string) {
	//make this function
}
