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
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/sabhiram/go-gitignore"
)

func isTextFile(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return false
	}

	contentType := http.DetectContentType(buffer[:n])
	return strings.HasPrefix(contentType, "text/")
}

func removeDuplicates(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range input {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func listFiles(paths []string) []string {
	var result []string

	for _, path := range(paths)  {
		log.Info().Msgf("Processing input file/dir '%s'", path)
		
		info, err := os.Stat(path)
		if err != nil {
			log.Error().Msgf("%s: error: %v\n", path, err)
			continue
		}

		if info.IsDir() {
			log.Debug().Msgf("%s is a directory\n", path)
			files, err := listDirFiles(path)
			if err != nil {
				log.Error().Msgf("%s: error: %v\n", path, err)
			} else {
				result = slices.Concat(result, files)
			}
		} else if isTextFile(path) {
			log.Debug().Msgf("%s is a text file\n", path)
			result = append(result, path)
		} else {
			log.Error().Msgf("%s is not a text file\n", path)
		}
	}
	return removeDuplicates(result)
}

func listDirFiles(root string) ([]string, error) {
	var files []string
	ig, err := ignore.CompileIgnoreFile(filepath.Join(root, ".gitignore"))
	
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Compute relative path for .gitignore matching
		relPath, _ := filepath.Rel(root, path)

		log.Debug().Msgf("path = '%s', relPath = '%s'", path, relPath)
		
		// Skip ignored files or directories
		if ig != nil && ig.MatchesPath(relPath) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if d.IsDir() && strings.HasPrefix(relPath, ".git") {
			log.Debug().Msgf("git dir '%s' will be skipped", path)
			return fs.SkipDir
		}
		if !d.IsDir() {
			log.Debug().Msgf("file '%s' will be parsed", path)
			files = append(files, path)
		}
		return nil
	})

	return files, err
}


func parseDirGitignore_old(root string, config Config) (StatsMap, error) {
	ig, err := ignore.CompileIgnoreFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse .gitignore: %w", err)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 8) // Or any concurrency limit
	results := make(chan StatsMap)
	counter := StatsMap{}

	// Walk directory
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing path %s: %v\n", path, err)
			return err
		}

		// Skip directories/files ignored by .gitignore
		relPath, _ := filepath.Rel(root, path)
		if ig != nil && ig.MatchesPath(relPath) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		// Process file in a goroutine
		sem <- struct{}{}
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			defer func() { <-sem }()
			result := parseFile(p, config)
			results <- result
		}(path)

		return nil
	})

	// Wait for all workers to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		counter.Merge(r)
	}

	return counter, nil
}
