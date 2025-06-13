// Copyright 2025 Matteo Redaelli
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	_ "embed"
	"encoding/json"
	"os"
)

//go:embed config.json
var configData []byte

type MultilineBlock struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`		
}

type LanguageConfig struct {
	MultilineComments []MultilineBlock `json:"multilineComments,omitempty"`
	SingleComments    []string `json:"singleComments,omitempty"`
	MultilineStrings   []MultilineBlock `json:"multilineString,omitempty"`
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
