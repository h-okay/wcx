package wc

import (
	"encoding/json"
	"fmt"
	"strings"
)

type CountSelection struct {
	Lines         bool
	Words         bool
	Chars         bool
	Bytes         bool
	MaxLineLength bool
}

func (s CountSelection) Fields() []string {
	fields := make([]string, 0, 5)
	if s.Lines {
		fields = append(fields, "lines")
	}
	if s.Words {
		fields = append(fields, "words")
	}
	if s.Chars {
		fields = append(fields, "chars")
	}
	if s.Bytes {
		fields = append(fields, "bytes")
	}
	if s.MaxLineLength {
		fields = append(fields, "maxLineLength")
	}

	return fields
}

func (s CountSelection) Metrics(counts Counts) []int {
	metrics := make([]int, 0, 5)
	if s.Lines {
		metrics = append(metrics, counts.Lines)
	}
	if s.Words {
		metrics = append(metrics, counts.Words)
	}
	if s.Chars {
		metrics = append(metrics, counts.Chars)
	}
	if s.Bytes {
		metrics = append(metrics, counts.Bytes)
	}
	if s.MaxLineLength {
		metrics = append(metrics, counts.MaxLineLength)
	}

	return metrics
}

func DefaultSelection() CountSelection {
	return CountSelection{Lines: true, Words: true, Bytes: true}
}

func SelectionFromFlags(lines, words, chars, bytes, maxLineLength bool) CountSelection {
	selection := CountSelection{
		Lines:         lines,
		Words:         words,
		Chars:         chars,
		Bytes:         bytes,
		MaxLineLength: maxLineLength,
	}

	if len(selection.Fields()) == 0 {
		return DefaultSelection()
	}

	return selection
}

func maxIntWidth(rows [][]int) int {
	width := 1
	for _, row := range rows {
		for _, value := range row {
			valueWidth := len(fmt.Sprint(value))
			if valueWidth > width {
				width = valueWidth
			}
		}
	}
	return width
}

func FormatTextRows(rows []OutputRow, selection CountSelection) string {
	return FormatTextRowsWithAlignment(rows, selection, true)
}

func FormatTextRowsWithAlignment(rows []OutputRow, selection CountSelection, align bool) string {
	if len(rows) == 0 {
		return ""
	}

	metricsPerRow := make([][]int, 0, len(rows))
	for _, row := range rows {
		metricsPerRow = append(metricsPerRow, selection.Metrics(row.Counts))
	}

	fieldCount := len(selection.Fields())
	width := maxIntWidth(metricsPerRow)

	lines := make([]string, 0, len(rows))
	for i, row := range rows {
		metrics := metricsPerRow[i]
		line := ""

		if fieldCount == 1 || !align {
			parts := make([]string, 0, len(metrics))
			for _, metric := range metrics {
				parts = append(parts, fmt.Sprint(metric))
			}
			line = strings.Join(parts, " ")
		} else {
			parts := make([]string, 0, fieldCount)
			for _, metric := range metrics {
				parts = append(parts, fmt.Sprintf("%*d", width, metric))
			}
			line = strings.Join(parts, " ")
		}

		if row.Name != "" {
			line += " " + row.Name
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

type JSONFileResult struct {
	File   string         `json:"file"`
	Counts map[string]int `json:"counts,omitempty"`
	Error  string         `json:"error,omitempty"`
}

type JSONOutput struct {
	Metrics []string         `json:"metrics"`
	Files   []JSONFileResult `json:"files,omitempty"`
	Total   map[string]int   `json:"total,omitempty"`
}

func BuildSelectedMetricsMap(selection CountSelection, counts Counts) map[string]int {
	selected := make(map[string]int)
	if selection.Lines {
		selected["lines"] = counts.Lines
	}
	if selection.Words {
		selected["words"] = counts.Words
	}
	if selection.Chars {
		selected["chars"] = counts.Chars
	}
	if selection.Bytes {
		selected["bytes"] = counts.Bytes
	}
	if selection.MaxLineLength {
		selected["maxLineLength"] = counts.MaxLineLength
	}

	return selected
}

func FormatJSON(rows []OutputRow, selection CountSelection, total *Counts) (string, error) {
	files := make([]JSONFileResult, 0, len(rows))
	for _, row := range rows {
		file := row.Name
		if file == "" {
			file = "stdin"
		}

		entry := JSONFileResult{File: file}
		if row.Error != nil {
			entry.Error = row.Error.Error()
		} else {
			entry.Counts = BuildSelectedMetricsMap(selection, row.Counts)
		}

		files = append(files, entry)
	}

	out := JSONOutput{
		Metrics: selection.Fields(),
		Files:   files,
	}
	if total != nil {
		out.Total = BuildSelectedMetricsMap(selection, *total)
	}

	raw, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}

	return string(raw), nil
}
