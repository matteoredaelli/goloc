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
	"bufio"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

type BlockType int
const (
    None BlockType = iota
    Comment
    String
)

type Block struct {
	blockType BlockType
	end_string string
}

func is_line_with_single_comment(line string, config LanguageConfig) bool {
	// line is already trimmed
	for _, s := range config.SingleComments {
		if strings.HasPrefix(line, s) {
			return true
		}
	}
	return false	
}


func find_end_block(line string, end_string string, index int) int {
	index_new := index + len(end_string)
	if len(line) > index_new {
		return strings.Index(line[index_new:], end_string)
	} else {
		return -1
	}
}

func find_start_block_comment(line string, config LanguageConfig, block *Block) error {
	idx_start := -1
	comment_idx := -1
	comment_end_block := ""
	string_idx := -1
	string_end_block := ""
	for _, b := range config.MultilineComments {
		i := strings.Index(line, b.Start)
		if i >= 0 && (comment_idx == -1 || i < comment_idx) {
			comment_idx       = i
			comment_end_block = b.End	
		}
	}
	for _, b := range config.MultilineStrings {
		i := strings.Index(line, b.Start)
		log.Printf("multiline string start=%v, u=%v", b.Start, i)
		if i >= 0 && (string_idx == -1 || i < string_idx) {
			string_idx       = i
			string_end_block = b.End	
		}
	}
	log.Printf("comment_idx=%v, string_idx=%v", comment_idx, string_idx)
	
	switch {
	case comment_idx >= 0 && (string_idx == -1 || comment_idx < string_idx):
		idx_start = comment_idx
		log.Printf("Multiline 'comment' startfound")
		block.blockType =  Comment
		block.end_string = comment_end_block
	case string_idx >= 0 && (comment_idx == -1 || string_idx < comment_idx):
		idx_start = string_idx
		log.Printf("Multiline 'string' start found")
		block.blockType = String
		block.end_string = string_end_block
	case comment_idx == string_idx:
	default:
		log.Printf("No multiline string/comment start found")
		block.blockType =  None
		block.end_string = ""
		return nil
	}
	log.Debug().Msgf("find_start_block_comment - block: %s", block)
	// TODO if multiline start and end patterns  are in the same row
	// checking if multine comments/docs start and end in teh same line
	if idx_start >= 0 && string_end_block != "" {
		idx_end := find_end_block(line, string_end_block, idx_start)
		if idx_end > -1 {
			log.Debug().Msg("Multirows comment/doc starts and ends in the same line")
			block.end_string = ""
		//		if idx_start == 0 and idx_end == (len(string) - len(string_end_block)) {
		}
	}
	
	//block.end_string = ""
	log.Debug().Msgf("find_start_block_comment - block: %s", block)
	return nil
}

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

func parseLine(line string, language string, config LanguageConfig, block *Block, stats *FileStats) {
	trimmed := strings.TrimSpace(line)
	log.Debug().Msgf("parseLine - Block: %s", block)
	switch block.blockType {
	case None:
		// not inside a multiline block
		if is_line_with_single_comment(trimmed, config) {
			stats.Comments++
			return
		}
		if trimmed == "" {
			stats.Blanks++
			return
		}
			
		find_start_block_comment(trimmed, config, block)
		log.Debug().Msgf("parseLine - Block: %s", block)
		switch block.blockType {
		case None: 
			stats.Code++
			return
		case Comment:
			stats.Comments++
			// starting a multiline block
			// end block could be in the same line
		case String:
			stats.Code++				
		}		
	case Comment:
		stats.Comments++
		
		// inside a multiline block
		if strings.Contains(trimmed, block.end_string) {
			log.Debug().Msg("Multirows comment end found")
			block.blockType = None
			block.end_string = ""
		}
	case String:
		stats.Code++
		
		// inside a multiline block
		if strings.Contains(trimmed, block.end_string) {
			log.Debug().Msg("Multirows string end found")
			block.blockType = None
			block.end_string = ""
		}
	// default:
	// 	//TODO Raise
		
	}
	log.Debug().Msgf("parseLine - block: %s", block)
}

func parseFile(filename string, config Config) StatsMap {
	var language string
	var languageConfig LanguageConfig
	ext := filepath.Ext(filename)
	if len(ext) > 1 {
		ext = ext[1:] // removes the dot
	}

	if lang, ok := config.Extensions[ext]; ok {
		language = lang
		languageConfig = config.Languages[language]
	} else {
		return StatsMap{ext: FileStats{Files: 1, Skipped: 1}}
	}

	file, err := os.Open(filename)
	if err != nil {
		return StatsMap{ext: FileStats{Files: 1, Skipped: 1}}
	}
	defer file.Close()

	block := Block{
		blockType:  None,
		end_string: "",
	}
	
	var stats FileStats
	stats.Files++
	
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		stats.Lines++
		log.Debug().Msgf("Line: '%s'", line)
		log.Debug().Msgf("parseFile - block: %s, stats: %s", block, stats)
		parseLine(line, language, languageConfig, &block, &stats)
		log.Debug().Msgf("parseFile - block: %s, stats: %s", block, stats)
	}
	
	resp := StatsMap{}
	resp[language] = stats
	return resp
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
