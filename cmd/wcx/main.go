package main

import (
	"errors"
	"fmt"
	"os"

	appcli "cc/wcx/internal/cli"
	"cc/wcx/internal/wc"
)

var version = "dev"

var errPartialFailure = errors.New("one or more inputs failed")

func main() {
	if err := run(os.Args[1:]); err != nil {
		if !errors.Is(err, errPartialFailure) {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}

func run(args []string) error {
	config, err := appcli.Parse(args)
	if err != nil {
		return err
	}

	if config.Help {
		fmt.Print(appcli.HelpText())
		return nil
	}

	if config.Version {
		fmt.Printf("wcx %s\n", version)
		return nil
	}

	inputs, err := wc.ResolveInputs(config.Args, config.Files0From)
	if err != nil {
		return err
	}

	options := wc.RunOptions{
		Selection: config.Selection,
		TotalMode: config.TotalMode,
		JSON:      config.JSON,
	}

	runResult := wc.Run(inputs, options)

	for _, row := range runResult.Rows {
		if row.Error != nil {
			name := row.Name
			if name == "" {
				name = "-"
			}
			_, _ = fmt.Fprintf(os.Stderr, "wcx: %s: %v\n", name, row.Error)
		}
	}

	output, err := wc.Render(runResult, options)
	if err != nil {
		return err
	}

	if output != "" {
		fmt.Println(output)
	}

	if runResult.HadErrors {
		return errPartialFailure
	}

	return nil
}
