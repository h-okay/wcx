package cli

import (
	"fmt"
	"strings"

	"cc/wcx/internal/wc"
)

type Config struct {
	Selection  wc.CountSelection
	TotalMode  wc.TotalMode
	Files0From string
	JSON       bool
	Help       bool
	Version    bool
	Args       []string
}

type parseFlags struct {
	lines         bool
	words         bool
	chars         bool
	bytes         bool
	maxLineLength bool
}

// Parse handles GNU-like short/long flags and keeps operands in original order.
// A lone "-" is treated as a file operand (stdin), not as an option prefix.
func Parse(args []string) (Config, error) {
	config := Config{TotalMode: wc.TotalAuto}
	flags := parseFlags{}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--" {
			config.Args = append(config.Args, args[i+1:]...)
			break
		}

		if strings.HasPrefix(arg, "--") {
			name, value, hasValue := splitLongOption(arg[2:])
			switch name {
			case "bytes":
				flags.bytes = true
			case "lines":
				flags.lines = true
			case "words":
				flags.words = true
			case "chars":
				flags.chars = true
			case "max-line-length":
				flags.maxLineLength = true
			case "files0-from":
				if !hasValue {
					if i+1 >= len(args) {
						return Config{}, fmt.Errorf("missing value for --files0-from")
					}
					i++
					value = args[i]
				}
				config.Files0From = value
			case "total":
				if !hasValue {
					if i+1 >= len(args) {
						return Config{}, fmt.Errorf("missing value for --total")
					}
					i++
					value = args[i]
				}
				mode, ok := wc.ParseTotalMode(value)
				if !ok {
					return Config{}, fmt.Errorf("invalid value for --total: use auto, always, only, or never")
				}
				config.TotalMode = mode
			case "json":
				config.JSON = true
			case "version":
				config.Version = true
			case "help":
				config.Help = true
			default:
				return Config{}, fmt.Errorf("unknown option: --%s", name)
			}
			continue
		}

		if strings.HasPrefix(arg, "-") && arg != "-" {
			for _, short := range arg[1:] {
				switch short {
				case 'c':
					flags.bytes = true
				case 'l':
					flags.lines = true
				case 'w':
					flags.words = true
				case 'm':
					flags.chars = true
				case 'L':
					flags.maxLineLength = true
				case 'h':
					config.Help = true
				default:
					return Config{}, fmt.Errorf("unknown option: -%c", short)
				}
			}
			continue
		}

		config.Args = append(config.Args, arg)
	}

	config.Selection = wc.SelectionFromFlags(
		flags.lines,
		flags.words,
		flags.chars,
		flags.bytes,
		flags.maxLineLength,
	)

	return config, nil
}

func splitLongOption(arg string) (name string, value string, hasValue bool) {
	parts := strings.SplitN(arg, "=", 2)
	if len(parts) == 2 {
		return parts[0], parts[1], true
	}

	return arg, "", false
}

// HelpText is static to keep output stable across Go versions and avoid
// differences from the standard flag package formatter.
func HelpText() string {
	return `NAME:
   wcx - print newline, word, and byte counts for each file

USAGE:
   wcx [OPTION]... [FILE]...

DESCRIPTION:
   Drop-in compatible wc replacement with optional JSON output.

OPTIONS:
   -c, --bytes             print the byte counts
   -l, --lines             print the newline counts
   -w, --words             print the word counts
   -m, --chars             print the character counts
   -L, --max-line-length   print the maximum display width
       --files0-from=F     read input from NUL-terminated names in file F
       --total=WHEN        WHEN to print total counts: auto, always, only, never
       --json              output counts as JSON (wcx extension)
       --version           output version information and exit
   -h, --help              show help
`
}
