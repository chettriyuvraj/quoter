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
	USAGE_STRING     = `
add: add quotes
			
usage: add`
)

type AddConfig struct {
	genre string
	quote string
}

func HandleAdd(w io.Writer, args []string) error {
	config, err := parseAddArgs(w, args)
	if err != nil {
		return err
	}

	err = runAddCmd(w, config)
	if err != nil {
		return err
	}

	return nil
}

func runAddCmd(w io.Writer, config AddConfig) error {
	/* Read current quotes file, or create one if it doesn't exist */
	f, err := os.OpenFile(PERSIST_FILENAME, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(f)
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
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = f.Write(writeData)
	if err != nil {
		return err
	}

	return nil
}

func parseAddArgs(w io.Writer, args []string) (AddConfig, error) {
	var config AddConfig

	/* Setup */
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.StringVar(&config.genre, "g", "misc", "genre to which the quote belongs")

	fs.SetOutput(w)
	fs.Usage = func() {
		fmt.Fprint(w, USAGE_STRING)
		fmt.Fprintln(w)
		fs.PrintDefaults()
	}

	/* Parse */
	err := fs.Parse(args)
	if err != nil {
		return config, err
	}

	/* First positional arg treated as quote and others ignored */
	if fs.NArg() == 0 {
		return config, ErrNoPositionalArgs
	}
	config.quote = fs.Arg(0)

	return config, nil
}

// func validateFlagSet(fs *flag.FlagSet) error {

// }
