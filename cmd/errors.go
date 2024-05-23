package cmd

import "errors"

var ErrNoPositionalArgs = errors.New("no positional arg provided")
var ErrNoQuotesFound = errors.New("no quotes found in database")
var ErrNoGenreSpecificQuotesFound = errors.New("no quotes found for specified genre")
