package wc_test

import (
	"testing"

	"cc/wcx/internal/wc"
)

func TestCountsForFixture(t *testing.T) {
	values, err := wc.ReadFile(testFileName)
	if err != nil {
		t.Fatalf("unable to read fixture file: %v", err)
	}

	tests := []struct {
		name string
		got  int
		want int
	}{
		{name: "bytes", got: wc.CountBytes(values), want: 3735},
		{name: "lines", got: wc.CountLines(values), want: 9},
		{name: "words", got: wc.CountWords(values), want: 551},
		{name: "chars", got: wc.CountChars(values), want: 3735},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.want {
				t.Fatalf("%s mismatch: got %d want %d", test.name, test.got, test.want)
			}
		})
	}
}

func TestCounterEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		check func(*testing.T, []byte)
	}{
		{
			name:  "chars skip invalid utf8",
			input: []byte{0xff, 'a'},
			check: func(t *testing.T, input []byte) {
				if got, want := wc.CountChars(input), 1; got != want {
					t.Fatalf("char count mismatch: got %d want %d", got, want)
				}
			},
		},
		{
			name:  "words treat invalid bytes as content",
			input: []byte{0xff, ' ', 'a'},
			check: func(t *testing.T, input []byte) {
				if got, want := wc.CountWords(input), 2; got != want {
					t.Fatalf("word count mismatch: got %d want %d", got, want)
				}
			},
		},
		{
			name:  "max line length handles tabs",
			input: []byte("ab\tc\n1234\n"),
			check: func(t *testing.T, input []byte) {
				if got, want := wc.CountMaxLineLength(input), 9; got != want {
					t.Fatalf("max line length mismatch: got %d want %d", got, want)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.check(t, test.input)
		})
	}
}
