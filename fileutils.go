package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/sabhiram/go-gitignore"
)

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func listFiles(root string) ([]string, error) {
	var files []string
	ig, err := ignore.CompileIgnoreFile(filepath.Join(root, ".gitignore"))
	
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Compute relative path for .gitignore matching
		relPath, _ := filepath.Rel(root, path)

		// Skip ignored files or directories
		if ig != nil && ig.MatchesPath(relPath) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func parseDir(root string, config Config) StatsMap {
	files, err := listFiles(root)
	
	if err != nil {
		panic(err)
	}

	return parseFiles(files, config)
}

func parseDirGitignore(root string, config Config) (StatsMap, error) {
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