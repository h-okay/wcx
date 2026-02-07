package wc

import (
	"bytes"
	"testing"
)

func FuzzCountReader(f *testing.F) {
	f.Add([]byte("hello world\n"))
	f.Add([]byte{0xff, 0xfe, '\n', 'x'})
	f.Add([]byte("\twide🙂\n"))

	selection := CountSelection{Lines: true, Words: true, Chars: true, Bytes: true, MaxLineLength: true}

	f.Fuzz(func(t *testing.T, input []byte) {
		counts, err := CountReader(bytes.NewReader(input), selection)
		if err != nil {
			t.Fatalf("CountReader returned unexpected error: %v", err)
		}

		if counts.Bytes < 0 || counts.Words < 0 || counts.Lines < 0 || counts.Chars < 0 || counts.MaxLineLength < 0 {
			t.Fatalf("negative count produced: %+v", counts)
		}

		if counts.Chars > counts.Bytes {
			t.Fatalf("chars cannot exceed bytes: %+v", counts)
		}

		if counts.Lines > counts.Bytes {
			t.Fatalf("lines cannot exceed bytes: %+v", counts)
		}
	})
}
