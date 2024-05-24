package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/chettriyuvraj/quoter/cmd"
)

var ErrInternal = errors.New("internal error")
var ErrInvalidSubCmd = errors.New("invalid subcommand")
var usageString = `
Usage: %s <COMMAND> [OPTIONS]`

func main() {
	handleCmd(os.Args)
}

func handleCmd(args []string) {
	var err error

	switch os.Args[1] {
	case "add":
		err = cmd.HandleAdd(os.Stdout, os.Args[2:])
	case "quote":
		err = cmd.HandleQuote(os.Stdout, os.Args[2:])
	case "-h":
		printUsage(os.Args[0])
	case "--help":
		printUsage(os.Args[0])
	default:
		err = ErrInvalidSubCmd
	}

	if err != nil {
		/* If not a known error, mark as internal error - don't expose */
		if !errors.Is(err, cmd.ErrNoQuotesFound) || !errors.Is(err, cmd.ErrNoGenreSpecificQuotesFound) {
			err = ErrInternal
		}
		fmt.Fprint(os.Stderr, "Error: ")
		fmt.Fprint(os.Stderr, err.Error())
		fmt.Fprintln(os.Stderr)
		printUsage(os.Args[0])
	}
}

func printUsage(cmdName string) {
	fmt.Fprintf(os.Stderr, usageString, cmdName)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr)
	fmt.Fprint(os.Stderr, "COMMANDS:")
	fmt.Fprintln(os.Stderr)
	cmd.HandleAdd(os.Stdout, []string{"-h"})
	cmd.HandleQuote(os.Stdout, []string{"-h"})
}
