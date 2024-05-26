package cmd

import (
	"bytes"
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

/* TODO: Have we added too much complexity for testing one simple thing? Is there something pre-existing that makes this easier? */

/* Utility to mimic a readwriteseeker in-place of an actual file */
type ReadWriteSeekerUtil struct {
	ReadSeeker *bytes.Reader
}

func (rws *ReadWriteSeekerUtil) Read(b []byte) (n int, err error) {
	return rws.ReadSeeker.Read(b)
}

func (rws *ReadWriteSeekerUtil) Seek(offset int64, whence int) (int64, error) {
	return rws.ReadSeeker.Seek(offset, whence)
}

func (rws *ReadWriteSeekerUtil) Write(b []byte) (n int, err error) {

	/* We have to write from the current seek position - get the current seek position */
	seekPosition, err := rws.Seek(0, io.SeekCurrent)
	if err != nil {
		return -1, err
	}

	/* Seek to 0 and grab data till current seek position */
	_, err = rws.Seek(0, 0)
	if err != nil {
		return -1, err
	}
	prevData := make([]byte, seekPosition)
	if len(prevData) > 0 {
		_, err = rws.Read(prevData)
		if err != nil {
			return -1, err
		}
	}

	newData := append(prevData, b...)
	rws.ReadSeeker = bytes.NewReader(newData)

	/* Set to correct seek position */
	_, err = rws.Seek(seekPosition+int64(len(b)), 0)
	if err != nil {
		return -1, err
	}

	return len(b), nil
}

func HandleError(w io.Writer, err error) { /* TODO: should this be unexported? */
	fmt.Fprint(w, err.Error())
}
