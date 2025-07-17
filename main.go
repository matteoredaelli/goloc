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
	"encoding/json"
	"flag"
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

	// Define a flag: -o csv
	outputFormat := flag.String("o", "table", "output format (e.g. csv, json)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] file1 file2 dir1...\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\nOptions:")
		flag.PrintDefaults()
	}
	// Parse the flags
	flag.Parse()

	// Remaining args after flags (e.g. file1, file2)
	input_files := flag.Args()
	log.Info().Msgf("Command line params: files or dirs: %v", input_files)

	files := listFiles(input_files)
	
	if len(files) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	counter := FileStatsMap{}

	counter = parseFiles(files, *config)

	summary := BuildSummaryStats(counter)
	
	switch *outputFormat {
	case "json":
		jsonBytes, err := json.Marshal(summary)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		jsonStr := string(jsonBytes)
		fmt.Println(jsonStr)
	case "table":
		PrintSummaryStatsTable(summary)
	default:
		log.Error().Msgf("Unknown output format (-o) '%s'", *outputFormat)
	}

}
