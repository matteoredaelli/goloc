package main

import (
	"fmt"
	"os"
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

type StatsMap map[string]FileStats

// Add adds values from another Stats to this one
func (s *FileStats) Add(other FileStats) {
	s.Files += other.Files
	s.Skipped += other.Skipped
	s.Lines += other.Lines
	s.Code += other.Code
	s.Comments += other.Comments
	s.Blanks += other.Blanks
}

func (sm StatsMap) Merge(other StatsMap) {
	for k, v2 := range other {
		if v1, exists := sm[k]; exists {
			v1.Add(v2)
			sm[k] = v1
		} else {
			sm[k] = v2
		}
	}
}

func PrintStatsMapTable(data StatsMap) {
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

	table.Render()
}