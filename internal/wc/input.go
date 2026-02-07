package wc

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type InputSource struct {
	Path        string
	DisplayName string
	FromStdin   bool
}

func ReadFile(filename string) ([]byte, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func ReadFromStdin() ([]byte, error) {
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	return stdin, nil
}

// OpenInput returns a stream for counting. For stdin operands this wraps the
// process stdin handle without taking ownership of it.
func OpenInput(input InputSource) (io.ReadCloser, error) {
	if input.FromStdin {
		return io.NopCloser(os.Stdin), nil
	}

	file, err := os.Open(input.Path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// ResolveInputs enforces GNU wc operand rules and normalizes input sources.
// When files0From is provided, positional operands are not allowed.
func ResolveInputs(args []string, files0From string) ([]InputSource, error) {
	if files0From != "" && len(args) > 0 {
		return nil, fmt.Errorf("file operands cannot be combined with --files0-from")
	}

	if files0From != "" {
		names, err := ReadFiles0From(files0From)
		if err != nil {
			return nil, err
		}
		return namesToInputs(names), nil
	}

	if len(args) == 0 {
		return []InputSource{{Path: "-", DisplayName: "", FromStdin: true}}, nil
	}

	return namesToInputs(args), nil
}

// ReadFiles0From parses a NUL-delimited file list used by --files0-from.
func ReadFiles0From(path string) ([]string, error) {
	var raw []byte
	var err error

	if path == "-" {
		raw, err = ReadFromStdin()
	} else {
		raw, err = ReadFile(path)
	}

	if err != nil {
		return nil, err
	}

	parts := bytes.Split(raw, []byte{0})
	names := make([]string, 0, len(parts))
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		names = append(names, string(part))
	}

	return names, nil
}

func namesToInputs(names []string) []InputSource {
	inputs := make([]InputSource, 0, len(names))
	for _, name := range names {
		if name == "-" {
			inputs = append(inputs, InputSource{Path: "-", DisplayName: "-", FromStdin: true})
			continue
		}

		inputs = append(inputs, InputSource{Path: name, DisplayName: name, FromStdin: false})
	}

	return inputs
}
