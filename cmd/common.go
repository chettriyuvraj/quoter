package cmd

import (
	"fmt"
	"io"
)

/* TODO: Must modify whenever flags or usage message changes  */
/* The strings below contain tabs and newlines in weird places to match the output for fs.Parse(), can we escape strings and test just the text? */
var completeAddUsageString string = `
add: add a new quote with an optional genre
			
Usage: add [OPTIONS] <quote>

OPTIONS:

  -g string
    	genre to which the quote belongs (default "misc")`

var completeQuoteUsageString string = `
quote: returns a random quote from stored list

Usage: quote [OPTIONS]

OPTIONS:

  -g string
    	genre from which we want a quote`

type Quote struct {
	Text  string `json:"text"`
	Genre string `json:"genre"`
}

func HandleError(stderr io.Writer, err error) { /* TODO: should this be unexported? */
	fmt.Fprint(stderr, err.Error())
}
