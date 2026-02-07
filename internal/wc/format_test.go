package wc_test

import (
	"testing"

	"cc/wcx/internal/wc"
)

func TestFormatTextRows(t *testing.T) {
	tests := []struct {
		name      string
		rows      []wc.OutputRow
		selection wc.CountSelection
		want      string
	}{
		{
			name: "multi-field alignment",
			rows: []wc.OutputRow{{
				Name: testFileName,
				Counts: wc.Counts{
					Lines: 2,
					Words: 10,
					Bytes: 5,
				},
			}},
			selection: wc.DefaultSelection(),
			want:      " 2 10  5 " + testFileName,
		},
		{
			name: "single-field no leading spaces",
			rows: []wc.OutputRow{{
				Name:   testFileName,
				Counts: wc.Counts{Lines: 2},
			}},
			selection: wc.SelectionFromFlags(true, false, false, false, false),
			want:      "2 " + testFileName,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := wc.FormatTextRows(test.rows, test.selection)
			if got != test.want {
				t.Fatalf("formatted output mismatch:\n got: %q\nwant: %q", got, test.want)
			}
		})
	}
}
