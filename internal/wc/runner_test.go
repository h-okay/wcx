package wc_test

import (
	"testing"

	"cc/wcx/internal/wc"
)

func TestParseTotalMode(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantMode  wc.TotalMode
		wantValid bool
	}{
		{name: "always", input: "always", wantMode: wc.TotalAlways, wantValid: true},
		{name: "invalid", input: "bogus", wantValid: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mode, ok := wc.ParseTotalMode(test.input)
			if ok != test.wantValid {
				t.Fatalf("validity mismatch: got %v want %v", ok, test.wantValid)
			}
			if ok && mode != test.wantMode {
				t.Fatalf("mode mismatch: got %q want %q", mode, test.wantMode)
			}
		})
	}
}

func TestRenderWithTotalOnly(t *testing.T) {
	selection := wc.DefaultSelection()
	result := wc.RunResult{
		Rows: []wc.OutputRow{
			{Name: "a.txt", Counts: wc.Counts{Lines: 10, Words: 20, Bytes: 300}},
			{Name: "b.txt", Counts: wc.Counts{Lines: 40, Words: 500, Bytes: 6000}},
		},
		Total:     wc.Counts{Lines: 50, Words: 520, Bytes: 6300},
		ShowTotal: true,
	}

	out, err := wc.Render(result, wc.RunOptions{Selection: selection, TotalMode: wc.TotalOnly})
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if out != "50 520 6300" {
		t.Fatalf("render output mismatch: got %q want %q", out, "50 520 6300")
	}
}
