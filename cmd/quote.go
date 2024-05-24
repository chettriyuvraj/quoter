package cmd

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"os"
)

type QuoteConfig struct {
	Genre string `json:"genre"`
}

const (
	QUOTE_USAGE_STRING = `
quote: returns a quote from stored list

Usage: quote [OPTIONS]` /* TODO: Add proper usage */
)

func HandleQuote(w io.Writer, args []string) error {
	/* Parse flags */
	config, err := parseQuoteArgs(w, args)
	if err != nil {
		return err
	}

	/* Open current quote file */
	f, err := os.Open(PERSIST_FILENAME)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return ErrNoQuotesFound
		}
		return err
	}
	defer f.Close()

	/* Run command */
	err = runQuoteCmd(w, f, config)
	if err != nil {
		return err
	}

	return nil
}

func runQuoteCmd(w io.Writer, quoteStorage io.ReadWriteSeeker, config QuoteConfig) error {
	var quotes []Quote

	/* Read entire contents of quoteStorage */
	_, err := quoteStorage.Seek(0, 0)
	if err != nil {
		return err
	}

	/* Unmarshal to slice of Quote */
	data, err := io.ReadAll(quoteStorage)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &quotes)
	if err != nil {
		return err
	}

	/* If no genre specified, write a random quote to writer */
	if config.Genre == "" {
		randIdx := rand.Intn(len(quotes))
		randQuote := quotes[randIdx]
		fmt.Fprint(w, randQuote.Text)
		fmt.Fprintln(w)
		return nil
	}

	/* If genre specified, find genre specific quotes */
	var genreSpecificQuotes []Quote
	for _, quote := range quotes {
		if quote.Genre == config.Genre {
			genreSpecificQuotes = append(genreSpecificQuotes, quote)
		}
	}
	if len(genreSpecificQuotes) == 0 {
		return ErrNoGenreSpecificQuotesFound
	}
	randIdx := rand.Intn(len(genreSpecificQuotes))
	randQuote := genreSpecificQuotes[randIdx]
	fmt.Fprint(w, randQuote.Text)
	fmt.Fprintln(w)

	return nil
}

func parseQuoteArgs(w io.Writer, args []string) (QuoteConfig, error) {
	var config QuoteConfig

	fs := flag.NewFlagSet("quote", flag.ContinueOnError)
	fs.StringVar(&config.Genre, "g", "", "genre from which we want a quote")

	fs.SetOutput(w)
	fs.Usage = func() {
		fmt.Fprint(w, QUOTE_USAGE_STRING)
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fmt.Fprint(w, "OPTIONS:")
		fmt.Fprintln(w)
		fmt.Fprintln(w)
		fs.PrintDefaults()
	}

	err := fs.Parse(args)
	if err != nil {
		return config, err
	}

	return config, nil
}
