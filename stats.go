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
	"sort"

	"github.com/olekukonko/tablewriter"
)

type FileStats struct {
	Files    int
	Skipped  int
	Lines    int
	Code     int
	Comments int
	Blanks   int
}

type FileStatsMap map[string]FileStats

type SummaryStatsMap struct {
	Totals            FileStats
	Stats             FileStatsMap
	MostUsedLanguage  string
}

// Add adds values from another Stats to this one
func (s *FileStats) Add(other FileStats) {
	s.Files += other.Files
	s.Skipped += other.Skipped
	s.Lines += other.Lines
	s.Code += other.Code
	s.Comments += other.Comments
	s.Blanks += other.Blanks
}

func (sm FileStatsMap) Merge(other FileStatsMap) {
	for k, v2 := range other {
		if v1, exists := sm[k]; exists {
			v1.Add(v2)
			sm[k] = v1
		} else {
			sm[k] = v2
		}
	}
}

func BuildSummaryStats(data FileStatsMap) SummaryStatsMap {
	var result SummaryStatsMap
	var maxLang string
	var maxFiles int
	var total FileStats
	
	for lang, stats := range data {
		total.Files += stats.Files
		total.Skipped += stats.Skipped
		total.Lines += stats.Lines
		total.Code += stats.Code
		total.Comments += stats.Comments
		total.Blanks += stats.Blanks
		
		if strings.HasPrefix(lang, "ext_") {
			continue
		}
		if stats.Files > maxFiles {
			maxFiles = stats.Files
			maxLang = lang
		}
	}
	result.Stats            = data
	result.Totals           = total
	result.MostUsedLanguage = maxLang
	
	return result	
}

func PrintSummaryStatsTable(summary SummaryStatsMap) {
	//var buf bytes.Buffer
	//table := tablewriter.NewWriter(&buf)
	data  := summary.Stats
	total := summary.Totals
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Lang", "Files", "Skipped", "Lines", "Code", "Comments", "Blanks"})

	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_LEFT,  // Lang
		tablewriter.ALIGN_RIGHT, // Files
		tablewriter.ALIGN_RIGHT, // Skipped
		tablewriter.ALIGN_RIGHT, // Lines
		tablewriter.ALIGN_RIGHT, // Code
		tablewriter.ALIGN_RIGHT, // Comments
		tablewriter.ALIGN_RIGHT, // Blanks
	})

	// Sort keys for consistent output
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := data[k]
		row := []string{
			k,
			fmt.Sprint(v.Files),
			fmt.Sprint(v.Skipped),
			fmt.Sprint(v.Lines),
			fmt.Sprint(v.Code),
			fmt.Sprint(v.Comments),
			fmt.Sprint(v.Blanks),
		}
		table.Append(row)
	}

	table.SetFooter([]string{
		"TOTAL",
		fmt.Sprint(total.Files),
		fmt.Sprint(total.Skipped),
		fmt.Sprint(total.Lines),
		fmt.Sprint(total.Code),
		fmt.Sprint(total.Comments),
		fmt.Sprint(total.Blanks),
	})
	table.SetFooterAlignment(tablewriter.ALIGN_RIGHT)
	
	table.Render()
}
