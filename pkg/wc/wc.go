package wc

import core "cc/wcx/internal/wc"

type (
	Counts         = core.Counts
	CountSelection = core.CountSelection
	InputSource    = core.InputSource
	OutputRow      = core.OutputRow
	RunOptions     = core.RunOptions
	RunResult      = core.RunResult
	TotalMode      = core.TotalMode
)

const (
	TotalAuto   = core.TotalAuto
	TotalAlways = core.TotalAlways
	TotalOnly   = core.TotalOnly
	TotalNever  = core.TotalNever
)

func DefaultSelection() CountSelection {
	return core.DefaultSelection()
}

func SelectionFromFlags(lines, words, chars, bytes, maxLineLength bool) CountSelection {
	return core.SelectionFromFlags(lines, words, chars, bytes, maxLineLength)
}

func ParseTotalMode(value string) (TotalMode, bool) {
	return core.ParseTotalMode(value)
}

func ResolveInputs(args []string, files0From string) ([]InputSource, error) {
	return core.ResolveInputs(args, files0From)
}

func Run(inputs []InputSource, options RunOptions) RunResult {
	return core.Run(inputs, options)
}

func Render(result RunResult, options RunOptions) (string, error) {
	return core.Render(result, options)
}
