package wc

import (
	"bytes"
	"testing"
)

func BenchmarkCountReader(b *testing.B) {
	line := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit.\n")
	input := bytes.Repeat(line, 50_000)
	selection := CountSelection{Lines: true, Words: true, Chars: true, Bytes: true, MaxLineLength: true}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := CountReader(bytes.NewReader(input), selection); err != nil {
			b.Fatalf("CountReader failed: %v", err)
		}
	}
}
