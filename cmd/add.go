package cmd

import (
	"flag"
)

type AddConfig struct {
	persist bool
	genre   string
}

// func HandleAdd(w io.Writer, args []string) error {

// }

func parseAddArgs(fs *flag.FlagSet, args []string) (AddConfig, error) {
	var config AddConfig
	fs.BoolVar(&config.persist, "p", false, "persist quote")
	fs.StringVar(&config.genre, "g", "misc", "genre to which the quote belongs")

	err := fs.Parse(args)
	if err != nil {
		return config, err
	}

	return config, nil
}

// func validateFlagSet(fs *flag.FlagSet) error {

// }
