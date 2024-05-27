package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

/* TODO: test errors? */
func TestAddQuoteToStorageSingle(t *testing.T) {
	tcs := []struct {
		desc    string
		config  AddConfig
		want    []Quote
		wantErr error
	}{
		{
			desc:   "single quote 1",
			config: AddConfig{quote: "Phool hu gulab ka, chameli ka mat samajhna..Aashiq hu aapka apni saheli ka mat samajhna!", genre: "misc"},
			want: []Quote{
				{Text: "Phool hu gulab ka, chameli ka mat samajhna..Aashiq hu aapka apni saheli ka mat samajhna!", Genre: "misc"},
			},
		},
	}

	for _, tc := range tcs {
		/* Add quote to a buffer using runAddCmd */
		quoteStorage := ReadWriteSeekerUtil{ReadSeeker: bytes.NewReader([]byte{})}
		err := addQuoteToStorage(&quoteStorage, tc.config)

		/* Compare results */
		if tc.wantErr != nil {
			require.ErrorIs(t, err, tc.wantErr, tc.desc)
			continue
		}

		require.NoError(t, err, tc.desc)
		testQuoteStorage(t, tc.want, &quoteStorage, tc.desc)
	}

}

/* Tests addition of multiple quotes sequentially, and subsequent JSON file formation*/
/* TODO: test errors? */
func TestAddQuoteToStorageMulti(t *testing.T) {
	tcs := []struct {
		desc   string
		config AddConfig
		want   []Quote
	}{
		{
			desc:   "multi quote 1",
			config: AddConfig{genre: "misc", quote: "Phool hu gulab ka, chameli ka mat samajhna..Aashiq hu aapka apni saheli ka mat samajhna!"},
			want: []Quote{
				{Text: "Phool hu gulab ka, chameli ka mat samajhna..Aashiq hu aapka apni saheli ka mat samajhna!", Genre: "misc"},
			},
		},
		{
			desc:   "multi quote 2",
			config: AddConfig{genre: "romance", quote: "Humse door jaoge kaise? Humko tum bhulaoge kaise? Hum vo khushbu hai jo saanson me baste hai, apni saanson ko rok paoge kaise?"},
			want: []Quote{
				{Text: "Humse door jaoge kaise? Humko tum bhulaoge kaise? Hum vo khushbu hai jo saanson me baste hai, apni saanson ko rok paoge kaise?", Genre: "romance"},
			},
		},
	}

	/* We are testing the addition of multiple quotes, so first compile all the quotes into a single 'want' array */
	want := []Quote{}
	quoteStorage := ReadWriteSeekerUtil{ReadSeeker: bytes.NewReader([]byte{})}
	for _, tc := range tcs {
		err := addQuoteToStorage(&quoteStorage, tc.config)
		require.NoError(t, err, tc.desc)
		want = append(want, tc.want...)
	}
	testQuoteStorage(t, want, &quoteStorage, "test multiple adds")

}

func testQuoteStorage(t *testing.T, want []Quote, quoteStorage io.ReadSeeker, desc string) {
	t.Helper()
	var got []Quote

	/* Quote storage's seek pointer may have been moved around, but now we want to read its entire contents, so bring seek back to 0 */
	_, err := quoteStorage.Seek(0, 0)
	if err != nil {
		require.NoError(t, err, desc)
	}

	/* Read entire quote storage, unmarshal it and compare it to what we are expecting */
	data, err := io.ReadAll(quoteStorage)
	require.NoError(t, err, desc)
	err = json.Unmarshal(data, &got)
	require.NoError(t, err, desc)
	require.Equal(t, want, got, desc)
}

/*
TODO:
- Test the error message for something like -g with no flag specified
*/
// Test HandleAdd
// {
// 	desc: "no quote",
// 	args: []string{"-g", "romance"},
// 	err:  ErrNoPositionalArgs,
// },
// {
// 	desc: "genre flag but no genre specified",
// 	args: []string{"-g", "Abki baar bichde toh khwabo me mile, jaise sookhe hue phool kitabo me mile.."},
// 	err:  ErrNoGenreSpecified,
// },

func TestParseAddArgs(t *testing.T) {

	tcs := []struct {
		desc       string
		args       []string
		want       AddConfig
		wantErr    error
		wantStdout string
		wantStderr string
	}{
		{
			desc:       "help flag",
			args:       []string{"-h"},
			wantErr:    flag.ErrHelp,
			wantStderr: ADD_USAGE_STRING,
		},
		{
			desc:       "non existent flag",
			args:       []string{"-x"},
			wantErr:    errors.New("flag provided but not defined: -x"),
			wantStderr: ADD_USAGE_STRING,
		},
		{
			desc:       "genre flag but no genre specifed",
			args:       []string{"-g"},
			wantErr:    errors.New("flag needs an argument: -g"),
			wantStderr: ADD_USAGE_STRING,
		},
		{
			desc:       "no positional args",
			args:       []string{},
			wantErr:    ErrNoPositionalArgs,
			wantStderr: ADD_USAGE_STRING,
		},
		{
			desc: "no flags",
			args: []string{"randomquote"},
			want: AddConfig{genre: "misc", quote: "randomquote"},
		},
		{
			desc: "genre flag only",
			args: []string{"-g", "romance", "randomquote"},
			want: AddConfig{genre: "romance", quote: "randomquote"},
		},
	}

	for _, tc := range tcs {

		/* Execute parse */
		errBuf := bytes.Buffer{}
		got, err := parseAddArgs(&errBuf, tc.args)

		if tc.wantErr != nil {
			/* Assert if error strings are the same - error not compared directly because internal errors are also returned which will not match with error.Is */
			require.ErrorContains(t, err, tc.wantErr.Error(), tc.desc)

			/* Formulate the error output we are expecting to receive: for 'errHelp' we are expecting usage string; for any other error we are expecting err + usageString  */
			errWantBuf := bytes.Buffer{}
			if err != flag.ErrHelp {
				fmt.Fprint(&errWantBuf, err.Error())
				fmt.Fprintln(&errWantBuf)
			}
			fmt.Fprintln(&errWantBuf, completeAddUsageString)

			/* Compare the error string we are expecting with the one we want */
			require.Equal(t, errWantBuf.String(), errBuf.String(), tc.desc)
			continue
		}

		/* Non-error case */
		require.NoError(t, err, tc.desc)
		require.Equal(t, tc.want, got, tc.desc)
		require.Equal(t, tc.wantStderr, errBuf.String(), tc.desc) /* Standard output remains empty */

	}
}
