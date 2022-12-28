package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Server struct {
		Host            string `json:"host"`
		Port            string `json:"port"`
		ReadTimeout     int    `json:"read_timeout"`
		WriteTimeout    int    `json:"write_timeout"`
		ShutDownTimeout int    `json:"shutdown_timeout"`
	} `json:"server"`
}

func LoadConfig(filename string) (Config, error) {
	config := Config{}
	configFile, err := os.Open(filename)
	if err != nil {
		return config, err
	}

	data, err := io.ReadAll(configFile)
	defer configFile.Close()
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
