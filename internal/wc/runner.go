package wc

import (
	"runtime"
	"strings"
	"sync"
)

type TotalMode string

const (
	TotalAuto   TotalMode = "auto"
	TotalAlways TotalMode = "always"
	TotalOnly   TotalMode = "only"
	TotalNever  TotalMode = "never"
)

type RunOptions struct {
	Selection CountSelection
	TotalMode TotalMode
	JSON      bool
}

type OutputRow struct {
	Name   string
	Counts Counts
	Error  error
}

type RunResult struct {
	Rows      []OutputRow
	Total     Counts
	ShowTotal bool
	HadErrors bool
}

func ParseTotalMode(value string) (TotalMode, bool) {
	mode := TotalMode(strings.ToLower(strings.TrimSpace(value)))
	switch mode {
	case TotalAuto, TotalAlways, TotalOnly, TotalNever:
		return mode, true
	default:
		return "", false
	}
}

// Run processes inputs, preserving input order in the returned rows even when
// file counting runs in parallel.
func Run(inputs []InputSource, options RunOptions) RunResult {
	rows := make([]OutputRow, len(inputs))

	if canRunInParallel(inputs) {
		runParallel(inputs, options.Selection, rows)
	} else {
		runSequential(inputs, options.Selection, rows)
	}

	total := Counts{}
	hadErrors := false
	successCount := 0

	for i := range rows {
		row := rows[i]
		if row.Error != nil {
			hadErrors = true
			continue
		}

		successCount++
		total.Lines += row.Counts.Lines
		total.Words += row.Counts.Words
		total.Chars += row.Counts.Chars
		total.Bytes += row.Counts.Bytes
		if row.Counts.MaxLineLength > total.MaxLineLength {
			total.MaxLineLength = row.Counts.MaxLineLength
		}
	}

	showTotal := shouldShowTotal(options.TotalMode, len(inputs), successCount)

	return RunResult{
		Rows:      rows,
		Total:     total,
		ShowTotal: showTotal,
		HadErrors: hadErrors,
	}
}

// canRunInParallel disables parallelism when stdin is present, since repeated
// reads from a shared stream are order-dependent.
func canRunInParallel(inputs []InputSource) bool {
	if len(inputs) < 2 {
		return false
	}

	for _, input := range inputs {
		if input.FromStdin {
			return false
		}
	}

	return true
}

func runSequential(inputs []InputSource, selection CountSelection, rows []OutputRow) {
	for i := range inputs {
		rows[i] = processInput(inputs[i], selection)
	}
}

func runParallel(inputs []InputSource, selection CountSelection, rows []OutputRow) {
	workerCount := min(runtime.GOMAXPROCS(0), len(inputs))

	jobs := make(chan int)
	var waitGroup sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			for index := range jobs {
				rows[index] = processInput(inputs[index], selection)
			}
		}()
	}

	for i := range inputs {
		jobs <- i
	}
	close(jobs)
	waitGroup.Wait()
}

func processInput(input InputSource, selection CountSelection) OutputRow {
	reader, err := OpenInput(input)
	if err != nil {
		return OutputRow{Name: input.DisplayName, Error: err}
	}
	defer reader.Close()

	counts, err := CountReader(reader, selection)
	if err != nil {
		return OutputRow{Name: input.DisplayName, Error: err}
	}

	return OutputRow{Name: input.DisplayName, Counts: counts}
}

func shouldShowTotal(mode TotalMode, inputCount int, successCount int) bool {
	if successCount == 0 {
		return false
	}

	switch mode {
	case TotalAlways, TotalOnly:
		return true
	case TotalNever:
		return false
	default:
		return inputCount > 1
	}
}

// Render applies --total and --json output policy to already computed rows.
func Render(result RunResult, options RunOptions) (string, error) {
	rows := make([]OutputRow, 0, len(result.Rows)+1)
	for _, row := range result.Rows {
		if row.Error != nil {
			continue
		}
		if options.TotalMode == TotalOnly {
			continue
		}
		rows = append(rows, row)
	}

	var totalForOutput *Counts
	if result.ShowTotal {
		totalForOutput = &result.Total
		totalRowName := "total"
		if options.TotalMode == TotalOnly {
			totalRowName = ""
		}
		if !options.JSON {
			rows = append(rows, OutputRow{Name: totalRowName, Counts: result.Total})
		}
	}

	if options.JSON {
		jsonRows := make([]OutputRow, 0, len(result.Rows))
		for _, row := range result.Rows {
			if options.TotalMode == TotalOnly && row.Error == nil {
				continue
			}
			jsonRows = append(jsonRows, row)
		}
		return FormatJSON(jsonRows, options.Selection, totalForOutput)
	}

	align := options.TotalMode != TotalOnly
	return FormatTextRowsWithAlignment(rows, options.Selection, align), nil
}
