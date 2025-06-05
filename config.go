package main

import (
	_ "embed"
	"encoding/json"
	"os"
)

//go:embed config.json
var configData []byte

type LanguageConfig struct {
	StartComment  string `json:"startComment,omitempty"`
	EndComment    string `json:"endComment,omitempty"`
	SingleComment string `json:"singleComment,omitempty"`
}

type Config struct {
	Languages  map[string]LanguageConfig `json:"languages"`
	Extensions map[string]string         `json:"extensions"`
	Filenames  map[string]string         `json:"filenames"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func LoadEmbeddedConfig() (*Config, error) {
	var config Config
	err := json.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}