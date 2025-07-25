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
	"errors"
	"path/filepath"
	"strings"
	//	"os"
)

//go:embed config.json
var configData []byte

type LanguageConfig struct {
	Quotes            [][]string `json:"quotes,omitempty"`
	MultilineComments [][]string `json:"multi_line_comments,omitempty"`
	SingleComments    []string `json:"line_comment,omitempty"`
	MultilineStrings  [][]string `json:"doc_quotes,omitempty"`
	Extensions        []string `json:"extensions"`
	Filenames         []string `json:"filenames"`
}


type Options struct {
	CountFiles  bool
	UnknownFiles bool
}

type Config struct {
	Languages  map[string]LanguageConfig `json:"languages"`
	Extensions map[string]string `json:"extensions"`
	Filenames  map[string]string `json:"filenames"`
	Options    Options `json:"options"`
}

func LoadEmbeddedConfig() (*Config, error) {
	var config Config

	err := json.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}

	// Initialize Extensions and Options maps before using it
	config.Extensions = make(map[string]string)
	config.Filenames = make(map[string]string)
	// move extensions
	for lang, value := range config.Languages {
		for _, ext := range value.Extensions {	
			config.Extensions[ext] = lang
		}
		for _, file := range value.Filenames {
			config.Filenames[file] = lang
		}
	}
	return &config, nil
}

func findLanguage(filename string, config Config) (string, error) {
	filename = strings.ToLower(filepath.Base(filename)) // Windows system
	ext := filepath.Ext(filename)
	if len(ext) > 1 {
		ext = ext[1:] // removes the dot
	}
	if lang, ok := config.Extensions[ext]; ok {
		return lang, nil
	} else {
		if lang, ok := config.Filenames[filename]; ok {
			return lang, nil
		} else {
			return ext, errors.New("unknown_extension_or_filename")
		}
	}
}
