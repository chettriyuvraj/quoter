package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	PERSIST_FILENAME = "quotes"
	ADD_USAGE_STRING = `
add: add a new quote with an optional genre
			
Usage: add [OPTIONS] <quote>` /* TODO: Write proper usage */
	ADD_SUCCESS_MSG = "Quote added successfully!"
)

/* TODO: Do you want the quote to be a part of the config itself here? If not, might require a redesign */
type AddConfig struct {
	genre string
	quote string
}

func HandleAdd(stdout, stderr io.Writer, args []string) error {
	config, err := parseAddArgs(stderr, args)
	if err != nil {
		/* Parse errors already printed to 'stderr' by fs.Parse command + additional errors handled inside parseAddArgs() */
		return err
	}

	/* Read current quotes file, or create one if it doesn't exist */
	f, err := os.OpenFile(PERSIST_FILENAME, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		HandleError(stderr, err)
		return err
	}
	defer f.Close()

	/* Parse all quotes in current file */
	quotes, err := parseQuotes(f)
	if err != nil {
		HandleError(stderr, err)
		return err
	}

	/* Add new quote to current list of quotes */
	quotes = addQuoteToList(config, quotes)
	if err != nil {
		HandleError(stderr, err)
		return err
	}

	/* Rewrite entire file */
	writeData, err := json.Marshal(quotes)
	if err != nil {
		HandleError(stderr, err)
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		HandleError(stderr, err)
		return err
	}
	_, err = f.Write(writeData)
	if err != nil {
		HandleError(stderr, err)
		return err
	}

	/* Write success to output */
	fmt.Fprint(stdout, ADD_SUCCESS_MSG)
	fmt.Fprintln(stdout)

	return nil
}

/* Parse raw json to quote struct */
func parseQuotes(quoteStorage io.ReadSeeker) ([]Quote, error) {
	var quotes []Quote

	/* Read entire contents of quoteStorage */
	_, err := quoteStorage.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(quoteStorage)
	if err != nil {
		return nil, err
	}

	/* Quotes are stored as JSON, parse them */
	if len(data) > 0 {
		err = json.Unmarshal(data, &quotes)
		if err != nil {
			return nil, err
		}
	}

	return quotes, nil
}

func addQuoteToList(config AddConfig, quotes []Quote) []Quote {
	return append(quotes, Quote{Text: config.quote, Genre: config.genre})
}

func parseAddArgs(stderr io.Writer, args []string) (AddConfig, error) {
	var config AddConfig

	/* Setup */
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.StringVar(&config.genre, "g", "misc", "genre to which the quote belongs")

	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprint(stderr, ADD_USAGE_STRING)
		fmt.Fprintln(stderr)
		fmt.Fprintln(stderr)
		fmt.Fprint(stderr, "OPTIONS:")
		fmt.Fprintln(stderr)
		fmt.Fprintln(stderr)
		fs.PrintDefaults()
	}

	/* Parse */
	err := fs.Parse(args)
	if err != nil {
		return config, err
	}

	/* First positional arg treated as quote and others ignored */
	if fs.NArg() == 0 {
		fmt.Fprint(stderr, ErrNoPositionalArgs)
		fmt.Fprintln(stderr)
		fs.Usage()
		return config, ErrNoPositionalArgs
	}
	config.quote = fs.Arg(0)

	return config, nil
}
