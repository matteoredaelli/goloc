package main

import (
	"fmt"
	"os"
)

func main() {
	config, err := LoadEmbeddedConfig()
	if err != nil {
		panic(fmt.Errorf("failed to parse embedded config: %w", err))
	}

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