package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/chettriyuvraj/quoter/cmd"
)

var ErrInvalidSubCmd = errors.New("invalid subcommand")
var ErrNoPositionalArgs = errors.New("no positional args provided")

const (
	USAGE_STRING = `
Usage: quoter <COMMAND> [OPTIONS]` /* TODO add dynamic command name */
)

func main() {
	err := handleCmd(os.Stdout, os.Stderr, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
}

/* No abstractions here - parsing and running simultaneously */
func handleCmd(stdout, stderr io.Writer, args []string) error {
	var err error

	fs := flag.NewFlagSet("quoter", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		printUsage(stdout, stderr)
	}

	/* Parse */
	err = fs.Parse(args)
	if err != nil {
		return err
	}

	/* First positional arg treated as sub-command */
	if fs.NArg() == 0 {
		fmt.Fprint(stderr, ErrNoPositionalArgs.Error())
		fmt.Fprintln(stderr)
		printUsage(stdout, stderr)
		return ErrNoPositionalArgs
	}
	switch fs.Arg(0) {
	case "add":
		err = cmd.HandleAdd(stdout, stderr, args[1:])
	case "quote":
		err = cmd.HandleQuote(stdout, stderr, args[1:])
	default:
		fmt.Fprint(stderr, ErrInvalidSubCmd.Error())
		fmt.Fprintln(stderr)
		printUsage(stdout, stderr)
		err = ErrInvalidSubCmd
	}

	return err

}

func printUsage(stdout, stderr io.Writer) {
	fmt.Fprint(stderr, USAGE_STRING)
	fmt.Fprintln(stderr)
	fmt.Fprintln(stderr)
	fmt.Fprint(stderr, "COMMANDS:")
	fmt.Fprintln(stderr)
	cmd.HandleAdd(stdout, stderr, []string{"-h"})
	cmd.HandleQuote(stdout, stderr, []string{"-h"})
}
