package wc_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"cc/wcx/internal/wc"
)

func TestReadInputSources(t *testing.T) {
	fileBytes, err := wc.ReadFile(testFileName)
	if err != nil {
		t.Fatalf("unable to read fixture file: %v", err)
	}

	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "read file",
			run: func(t *testing.T) {
				values, err := wc.ReadFile(testFileName)
				if err != nil {
					t.Fatalf("ReadFile failed: %v", err)
				}
				if len(values) != 3735 {
					t.Fatalf("ReadFile length mismatch: got %d want 3735", len(values))
				}
			},
		},
		{
			name: "read stdin",
			run: func(t *testing.T) {
				r, w, err := os.Pipe()
				if err != nil {
					t.Fatalf("unable to create os.Pipe: %v", err)
				}

				if _, err := w.Write(fileBytes); err != nil {
					t.Fatalf("unable to write to pipe: %v", err)
				}
				if err := w.Close(); err != nil {
					t.Fatalf("unable to close writer pipe: %v", err)
				}

				defer func(original *os.File) { os.Stdin = original }(os.Stdin)
				os.Stdin = r

				values, err := wc.ReadFromStdin()
				if err != nil {
					t.Fatalf("ReadFromStdin failed: %v", err)
				}
				if !bytes.Equal(values, fileBytes) {
					t.Fatalf("ReadFromStdin result mismatch")
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.run(t)
		})
	}
}

func TestResolveInputs(t *testing.T) {
	tmp := t.TempDir()
	listPath := filepath.Join(tmp, "files0.list")
	if err := os.WriteFile(listPath, []byte("one.txt\x00-\x00two.txt\x00"), 0o644); err != nil {
		t.Fatalf("unable to write files0 list: %v", err)
	}

	tests := []struct {
		name    string
		args    []string
		files0  string
		check   func(*testing.T, []wc.InputSource)
		wantErr bool
	}{
		{
			name:   "implicit stdin when no args",
			args:   nil,
			files0: "",
			check: func(t *testing.T, inputs []wc.InputSource) {
				if len(inputs) != 1 {
					t.Fatalf("input count mismatch: got %d want 1", len(inputs))
				}
				if !inputs[0].FromStdin || inputs[0].DisplayName != "" {
					t.Fatalf("unexpected implicit stdin input: %+v", inputs[0])
				}
			},
		},
		{
			name:   "files0 resolves inputs",
			args:   nil,
			files0: listPath,
			check: func(t *testing.T, inputs []wc.InputSource) {
				if len(inputs) != 3 {
					t.Fatalf("input count mismatch: got %d want 3", len(inputs))
				}
				if inputs[0].DisplayName != "one.txt" {
					t.Fatalf("unexpected first input: %+v", inputs[0])
				}
				if !inputs[1].FromStdin || inputs[1].DisplayName != "-" {
					t.Fatalf("unexpected stdin placeholder input: %+v", inputs[1])
				}
				if inputs[2].DisplayName != "two.txt" {
					t.Fatalf("unexpected last input: %+v", inputs[2])
				}
			},
		},
		{
			name:    "reject file operands with files0",
			args:    []string{"a.txt"},
			files0:  "list.txt",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputs, err := wc.ResolveInputs(test.args, test.files0)
			if test.wantErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("ResolveInputs failed: %v", err)
			}
			test.check(t, inputs)
		})
	}
}
