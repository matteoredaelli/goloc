package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
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

func parseFile(filename string, config Config) StatsMap {
	ext := filepath.Ext(filename)
	if len(ext) > 1 {
		ext = ext[1:] // removes the dot
	}

	language := ""
	startComment := ""
	endComment := ""
	singleComment := ""

	if lang, ok := config.Extensions[ext]; ok {
		language = lang
		startComment = config.Languages[lang].StartComment
		endComment = config.Languages[lang].EndComment
		singleComment = config.Languages[lang].SingleComment
	} else {
		return StatsMap{ext: FileStats{Files: 1, Skipped: 1}}
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
		return StatsMap{ext: FileStats{Files: 1, Skipped: 1}}
	}

	content := string(data)

	// Count original lines (add 1 if file doesn't end with newline)
	originalLines := strings.Count(content, "\n")
	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		originalLines++
	}

	// Remove block comments (non-greedy)
	if startComment != "" && endComment != "" {
		reBlockComments := regexp.MustCompile("(?s)" + regexp.QuoteMeta(startComment) + ".*?" + regexp.QuoteMeta(endComment))
		content = reBlockComments.ReplaceAllString(content, singleComment)
	}

	comments := originalLines - strings.Count(content, "\n")
	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		comments += originalLines - (strings.Count(content, "\n") + 1)
	}

	// Process line-by-line
	code, blanks := 0, 0
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			blanks++
		} else if singleComment != "" && strings.HasPrefix(trimmed, singleComment) {
			comments++
		} else {
			code++
		}
	}

	if blanks > 0 {
		blanks--
	}

	return StatsMap{language: FileStats{
		Files:    1,
		Skipped:  0,
		Lines:    originalLines,
		Code:     code,
		Comments: comments,
		Blanks:   blanks,
	}}
}

func parseFiles(files []string, config Config) StatsMap {
	var wg sync.WaitGroup
	results := make(chan StatsMap, len(files)) // buffered to avoid blocking

	for _, file := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			stats := parseFile(f, config)
			results <- stats
		}(file)
	}

	wg.Wait()
	close(results)

	// Collect and return counter
	counter := StatsMap{}
	for result := range results {
		counter.Merge(result)
	}
	return counter
}