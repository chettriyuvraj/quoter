package cmd

import (
	"flag"
	"fmt"
	"io"
)

const (
	PERSISTFILENAME = "quotes"
	USAGE_STRING    = `
add: add quotes
			
usage: add`
)

type AddConfig struct {
	genre string
}

// func HandleAdd(w io.Writer, args []string) error {
// 	config :=
// }

func parseAddArgs(w io.Writer, args []string) (AddConfig, error) {
	var config AddConfig

	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.StringVar(&config.genre, "g", "misc", "genre to which the quote belongs")

	fs.SetOutput(w)
	fs.Usage = func() {
		fmt.Fprint(w, USAGE_STRING)
		fmt.Fprintln(w)
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return config, err
	}

	return config, nil
}

// func validateFlagSet(fs *flag.FlagSet) error {

// }
