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
	"fmt"
	"os"
	"strings"
	
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	config, err := LoadEmbeddedConfig()
	if err != nil {
		panic(fmt.Errorf("failed to parse embedded config: %w", err))
	}

	// Set output to human-friendly format (optional, for console debugging)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Read log level from environment variable
	levelStr := strings.ToLower(os.Getenv("LOGGING"))

	// Set default log level
	level := zerolog.InfoLevel

	// Parse level from environment
	switch strings.ToLower(levelStr) {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn", "warning":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	case "fatal":
		level = zerolog.FatalLevel
	case "panic":
		level = zerolog.PanicLevel
	case "trace":
		level = zerolog.TraceLevel
	default:
		// fallback, already set to InfoLevel
	}

	// Set global log level
	zerolog.SetGlobalLevel(level)
	
	if len(os.Args) < 2 {
		fmt.Println("Usage: tool <directory>")
		os.Exit(1)
	}

	directoryOrFile := os.Args[1]
	counter := StatsMap{}

	if info, err := os.Stat(directoryOrFile); err == nil && info.IsDir() {
		fmt.Println("Directory exists")
		counter = parseDir(directoryOrFile, *config)
	} else if os.IsNotExist(err) {
		fmt.Printf("Not existing Directory or file: %v\n", err)
		os.Exit(1)
	} else {
		counter = parseFile(directoryOrFile, *config)
	}

	PrintStatsMapTable(counter)
}
