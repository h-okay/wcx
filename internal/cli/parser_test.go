package cli

import (
	"reflect"
	"testing"

	"cc/wcx/internal/wc"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		check     func(*testing.T, Config)
		wantError bool
	}{
		{
			name: "default selection and args",
			args: []string{"file.txt"},
			check: func(t *testing.T, config Config) {
				if config.Selection != wc.DefaultSelection() {
					t.Fatalf("default selection mismatch: %+v", config.Selection)
				}
				if config.TotalMode != wc.TotalAuto {
					t.Fatalf("default total mode mismatch: got %q", config.TotalMode)
				}
				if !reflect.DeepEqual(config.Args, []string{"file.txt"}) {
					t.Fatalf("args mismatch: got %#v", config.Args)
				}
			},
		},
		{
			name: "short flags are parsed",
			args: []string{"-lm", "file.txt"},
			check: func(t *testing.T, config Config) {
				want := wc.SelectionFromFlags(true, false, true, false, false)
				if config.Selection != want {
					t.Fatalf("selection mismatch: got %+v want %+v", config.Selection, want)
				}
			},
		},
		{
			name: "long options with equals",
			args: []string{"--total=never", "--files0-from=list.bin"},
			check: func(t *testing.T, config Config) {
				if config.TotalMode != wc.TotalNever {
					t.Fatalf("total mode mismatch: got %q want %q", config.TotalMode, wc.TotalNever)
				}
				if config.Files0From != "list.bin" {
					t.Fatalf("files0-from mismatch: got %q", config.Files0From)
				}
			},
		},
		{
			name:      "invalid total value returns error",
			args:      []string{"--total=bad"},
			wantError: true,
		},
		{
			name:      "unknown option returns error",
			args:      []string{"--whoops"},
			wantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, err := Parse(test.args)
			if test.wantError {
				if err == nil {
					t.Fatalf("expected parse error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected parse error: %v", err)
			}

			test.check(t, config)
		})
	}
}
