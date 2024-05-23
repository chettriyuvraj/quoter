package cmd

import (
	"flag"
	"fmt"
	"io"
)

type QuoteConfig struct {
	Genre string `json:"genre"`
}

const (
	QUOTE_USAGE_STRING = `
quote: returns a quote from stored list

usage: quote` /* TODO: Add proper usage */
)

func parseQuoteArgs(w io.Writer, args []string) (QuoteConfig, error) {
	var config QuoteConfig

	fs := flag.NewFlagSet("quote", flag.ContinueOnError)
	fs.StringVar(&config.Genre, "g", "", "genre from which we want a quote")

	fs.SetOutput(w)
	fs.Usage = func() {
		fmt.Fprint(w, QUOTE_USAGE_STRING)
		fmt.Fprintln(w)
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return config, err
	}

	return config, nil
}
