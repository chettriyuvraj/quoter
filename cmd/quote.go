package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
)

type QuoteConfig struct {
	Genre string `json:"genre"`
}

const (
	QUOTE_USAGE_STRING = `
quote: returns a random quote from stored list

Usage: quote [OPTIONS]` /* TODO: Add proper usage */
)

/* TODO: Pass errwriter and stdoutwriter separately? */

func HandleQuote(stdout, stderr io.Writer, args []string) error {
	/* Parse flags */
	config, err, errStr := parseQuoteArgs(args)
	if err != nil {
		HandleError(stderr, errors.New(errStr))
		return err
	}

	/* Open current quote file */
	f, err := os.Open(PERSIST_FILENAME)
	if err != nil {
		HandleError(stderr, err)
		return err
	}
	defer f.Close()

	/* Run command */
	quote, err := getRandomQuote(f, config)
	if err != nil {
		HandleError(stderr, err)
		return err
	}
	fmt.Fprint(stdout, quote.Text)
	fmt.Fprintln(stdout)

	return nil
}

func getRandomQuote(quoteStorage io.ReadSeeker, config QuoteConfig) (Quote, error) {
	var quotes []Quote

	/* Read entire contents of quoteStorage */
	_, err := quoteStorage.Seek(0, 0)
	if err != nil {
		return Quote{}, err
	}

	/* Unmarshal to slice of Quote */
	data, err := io.ReadAll(quoteStorage)
	if err != nil {
		return Quote{}, err
	}
	err = json.Unmarshal(data, &quotes)
	if err != nil {
		return Quote{}, err
	}

	/* If no genre specified, write a random quote to writer */
	if config.Genre == "" {
		randIdx := rand.Intn(len(quotes))
		return quotes[randIdx], err
	}

	/* If genre specified, find genre specific quotes */
	var genreSpecificQuotes []Quote
	for _, quote := range quotes {
		if quote.Genre == config.Genre {
			genreSpecificQuotes = append(genreSpecificQuotes, quote)
		}
	}
	if len(genreSpecificQuotes) == 0 {
		return Quote{}, ErrNoGenreSpecificQuotesFound
	}
	randIdx := rand.Intn(len(genreSpecificQuotes))
	return genreSpecificQuotes[randIdx], nil
}

func parseQuoteArgs(args []string) (config QuoteConfig, err error, errStr string) {
	var errBuf bytes.Buffer

	fs := flag.NewFlagSet("quote", flag.ContinueOnError)
	fs.StringVar(&config.Genre, "g", "", "genre from which we want a quote")

	fs.SetOutput(&errBuf)
	fs.Usage = func() {
		fmt.Fprint(&errBuf, QUOTE_USAGE_STRING)
		fmt.Fprintln(&errBuf)
		fmt.Fprintln(&errBuf)
		fmt.Fprint(&errBuf, "OPTIONS:")
		fmt.Fprintln(&errBuf)
		fmt.Fprintln(&errBuf)
		fs.PrintDefaults()
	}

	err = fs.Parse(args)
	if err != nil {
		return config, err, errBuf.String()
	}

	return config, nil, errBuf.String()
}
