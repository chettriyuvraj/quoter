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
		/* Parse errors already printed to 'w' by fs.Parse command + additional errors handled inside parseAddArgs() */
		return err
	}

	/* Read current quotes file, or create one if it doesn't exist */
	f, err := os.OpenFile(PERSIST_FILENAME, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		HandleError(stderr, err)
		return err
	}
	defer f.Close()

	err = addQuoteToStorage(f, config)
	if err != nil {
		HandleError(stderr, err)
		return err
	}
	fmt.Fprint(stdout, ADD_SUCCESS_MSG)
	fmt.Fprintln(stdout)

	return nil
}

func addQuoteToStorage(quoteStorage io.ReadWriteSeeker, config AddConfig) error {

	/* Read entire contents of quoteStorage */
	_, err := quoteStorage.Seek(0, 0)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(quoteStorage)
	if err != nil {
		return err
	}

	/* Quotes are stored as JSON, parse them */
	var quotes []Quote
	if len(data) > 0 {
		err = json.Unmarshal(data, &quotes)
		if err != nil {
			return err
		}
	}

	/* Append current quote to current set of quotes and rewrite entire file */
	quotes = append(quotes, Quote{Text: config.quote, Genre: config.genre})
	writeData, err := json.Marshal(quotes)
	if err != nil {
		return err
	}
	_, err = quoteStorage.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = quoteStorage.Write(writeData)
	if err != nil {
		return err
	}

	return nil
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
