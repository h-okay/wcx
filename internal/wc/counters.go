package wc

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

type Counts struct {
	Lines         int `json:"lines"`
	Words         int `json:"words"`
	Chars         int `json:"chars"`
	Bytes         int `json:"bytes"`
	MaxLineLength int `json:"maxLineLength"`
}

func CountAll(values []byte) Counts {
	selection := CountSelection{Lines: true, Words: true, Chars: true, Bytes: true, MaxLineLength: true}
	counts, _ := CountReader(bytes.NewReader(values), selection)
	return counts
}

func CountBytes(values []byte) int {
	return len(values)
}

func CountLines(values []byte) int {
	selection := CountSelection{Lines: true}
	counts, _ := CountReader(bytes.NewReader(values), selection)
	return counts.Lines
}

func CountWords(values []byte) int {
	selection := CountSelection{Words: true}
	counts, _ := CountReader(bytes.NewReader(values), selection)
	return counts.Words
}

func CountChars(values []byte) int {
	selection := CountSelection{Chars: true}
	counts, _ := CountReader(bytes.NewReader(values), selection)
	return counts.Chars
}

func CountMaxLineLength(values []byte) int {
	selection := CountSelection{MaxLineLength: true}
	counts, _ := CountReader(bytes.NewReader(values), selection)
	return counts.MaxLineLength
}

// CountReader computes all requested metrics in one pass over reader.
// Invalid UTF-8 bytes are counted as bytes, treated as non-whitespace for
// words, skipped for chars, and contribute zero display width.
func CountReader(reader io.Reader, selection CountSelection) (Counts, error) {
	if selection.Bytes && !selection.Lines && !selection.Words && !selection.Chars && !selection.MaxLineLength {
		size, err := io.Copy(io.Discard, reader)
		if err != nil {
			return Counts{}, err
		}
		return Counts{Bytes: int(size)}, nil
	}

	counts := Counts{}
	inWord := false
	currentLineWidth := 0
	posixMode := os.Getenv("POSIXLY_CORRECT") != ""

	buf := bufio.NewReaderSize(reader, 64*1024)
	for {
		r, size, err := buf.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Counts{}, err
		}

		if selection.Bytes {
			counts.Bytes += size
		}

		isEncodingError := r == utf8.RuneError && size == 1

		if selection.Lines && r == '\n' {
			counts.Lines++
		}

		if selection.MaxLineLength {
			switch r {
			case '\n':
				if currentLineWidth > counts.MaxLineLength {
					counts.MaxLineLength = currentLineWidth
				}
				currentLineWidth = 0
			case '\t':
				currentLineWidth += 8 - (currentLineWidth % 8)
			default:
				if !isEncodingError {
					currentLineWidth += runeDisplayWidth(r)
				}
			}
		}

		if selection.Chars && !isEncodingError {
			counts.Chars++
		}

		if selection.Words {
			isWhitespace := !isEncodingError && IsWhitespace(r, posixMode)
			if isWhitespace {
				inWord = false
			} else if !inWord {
				counts.Words++
				inWord = true
			}
		}
	}

	if selection.MaxLineLength && currentLineWidth > counts.MaxLineLength {
		counts.MaxLineLength = currentLineWidth
	}

	return counts, nil
}

// IsWhitespace follows GNU wc behavior: in non-POSIX mode it also treats
// U+00A0, U+2007, U+202F, and U+2060 as whitespace.
func IsWhitespace(r rune, posixMode bool) bool {
	if unicode.IsSpace(r) {
		return true
	}

	if posixMode {
		return false
	}

	switch r {
	case '\u00A0', '\u2007', '\u202F', '\u2060':
		return true
	default:
		return false
	}
}
