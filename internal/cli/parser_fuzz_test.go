package cli

import (
	"strings"
	"testing"
)

func FuzzParse(f *testing.F) {
	seeds := []string{
		"",
		"-l --words file.txt",
		"--total=only --json a.txt b.txt",
		"--files0-from=list.bin",
		"-L -m -c -w -l",
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, raw string) {
		if len(raw) > 2048 {
			raw = raw[:2048]
		}

		args := strings.Fields(raw)
		_, _ = Parse(args)
	})
}
