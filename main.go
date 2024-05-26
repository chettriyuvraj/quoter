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
	err := runCmd(os.Stderr, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
}

/* No abstractions here - parsing and running simultaneously */
func runCmd(w io.Writer, args []string) error {
	var err error

	fs := flag.NewFlagSet("quoter", flag.ContinueOnError)
	fs.SetOutput(w)
	fs.Usage = func() {
		printUsage(w)
	}

	/* Parse */
	err = fs.Parse(args)
	if err != nil {
		return err
	}

	/* First positional arg treated as sub-command */
	if fs.NArg() == 0 {
		fmt.Fprint(w, ErrNoPositionalArgs.Error())
		fmt.Fprintln(w)
		printUsage(w)
		return ErrNoPositionalArgs
	}
	switch fs.Arg(0) {
	case "add":
		err = cmd.HandleAdd(w, args[1:])
	case "quote":
		err = cmd.HandleQuote(w, args[1:])
	default:
		fmt.Fprint(w, ErrInvalidSubCmd.Error())
		fmt.Fprintln(w)
		printUsage(w)
		err = ErrInvalidSubCmd
	}

	return err

}

func printUsage(w io.Writer) {
	fmt.Fprint(w, USAGE_STRING)
	fmt.Fprintln(w)
	fmt.Fprintln(w)
	fmt.Fprint(w, "COMMANDS:")
	fmt.Fprintln(w)
	cmd.HandleAdd(w, []string{"-h"})
	cmd.HandleQuote(w, []string{"-h"})
}
