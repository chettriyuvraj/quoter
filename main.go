package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/chettriyuvraj/quoter/cmd"
)

func main() {
	switch os.Args[1] {
	case "add":
		err := cmd.HandleAdd(os.Stdout, os.Args[2:])
		if err != nil {
			if errors.Is(err, cmd.ErrNoQuotesFound) {
				fmt.Fprint(os.Stderr, err.Error())
				fmt.Fprintln(os.Stderr)
			}
			printUsage()
		}
	case "quote":
		err := cmd.HandleQuote(os.Stdout, os.Args[2:])
		if err != nil {
			if errors.Is(err, cmd.ErrNoQuotesFound) || errors.Is(err, cmd.ErrNoGenreSpecificQuotesFound) {
				fmt.Fprint(os.Stderr, err.Error())
				fmt.Fprintln(os.Stderr)
				return
			}
			printUsage()
		}
	}
}

func printUsage() {
	usageString := `
usage: %s <add|quote> [-g genre] [quote]`
	fmt.Fprint(os.Stderr, usageString)
}
