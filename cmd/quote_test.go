package cmd

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseQuoteArgs(t *testing.T) {

	tcs := []struct {
		desc       string
		args       []string
		want       QuoteConfig
		wantErr    error
		wantStdout string
		wantStderr string
	}{
		{
			desc:       "help flag",
			args:       []string{"-h"},
			wantErr:    flag.ErrHelp,
			wantStderr: completeQuoteUsageString + "\n",
		},
		{
			desc:       "non existent flag",
			args:       []string{"-x"},
			wantErr:    errors.New("flag provided but not defined: -x"),
			wantStderr: "flag provided but not defined: -x" + "\n" + completeQuoteUsageString + "\n",
		},
		{
			desc:       "genre flag but no genre specifed",
			args:       []string{"-g"},
			wantErr:    errors.New("flag needs an argument: -g"),
			wantStderr: "flag needs an argument: -g" + "\n" + completeQuoteUsageString + "\n",
		},
		{
			desc: "genre flag only",
			args: []string{"-g", "romance"},
			want: QuoteConfig{Genre: "romance"},
		},
		{
			desc: "no flags",
			args: []string{""},
			want: QuoteConfig{Genre: ""},
		},
	}

	for _, tc := range tcs {

		/* Execute parse */
		errBuf := bytes.Buffer{}
		got, err := parseQuoteArgs(&errBuf, tc.args)

		if tc.wantErr != nil {
			/* Assert if error strings are the same - error not compared directly because internal errors are also returned which will not match with error.Is */
			require.ErrorContains(t, err, tc.wantErr.Error(), tc.desc)

			/* Compare the error string we are expecting with the one we want */
			require.Equal(t, tc.wantStderr, errBuf.String(), tc.desc)
			continue
		}

		/* Non-error case */
		require.NoError(t, err, tc.desc)
		require.Equal(t, tc.want, got, tc.desc)
		require.Equal(t, tc.wantStderr, errBuf.String(), tc.desc) /* Standard output remains empty */

	}
}

func addRandomQuotes(t *testing.T, quoteStorage io.ReadWriteSeeker) {
	t.Helper()

	quoteConfigs := []AddConfig{
		{genre: "misc", quote: "Phool hu gulab ka, chameli ka mat samajhna..Aashiq hu aapka apni saheli ka mat samajhna!"},
		{genre: "romance", quote: "Humse door jaoge kaise? Humko tum bhulaoge kaise? Hum vo khushbu hai jo saanson me baste hai, apni saanson ko rok paoge kaise?"},
	}
	for _, config := range quoteConfigs {
		err := addQuoteToStorage(quoteStorage, config)
		require.NoError(t, err)
	}

}

func TestGetRandomQuoteCmd(t *testing.T) {
	/* Add quotes to a storage first */
	quoteStorage := ReadWriteSeekerUtil{ReadSeeker: bytes.NewReader([]byte{})}
	addRandomQuotes(t, &quoteStorage)

	tcs := []struct {
		desc    string
		config  QuoteConfig
		want    Quote
		wantErr error
	}{
		{
			desc:   "quote of romance genre",
			config: QuoteConfig{Genre: "romance"},
			want:   Quote{Genre: "romance", Text: "Humse door jaoge kaise? Humko tum bhulaoge kaise? Hum vo khushbu hai jo saanson me baste hai, apni saanson ko rok paoge kaise?"},
		},
		// { TODO: How to test random quote?
		// 	desc:   "quote with no genre specified",
		// 	config: QuoteConfig{Genre: "romance"},
		// 	want: "Humse door jaoge kaise? Humko tum bhulaoge kaise? Hum vo khushbu hai jo saanson me baste hai, apni saanson ko rok paoge kaise?"
		// },
	}

	for _, tc := range tcs {
		got, err := getRandomQuote(&quoteStorage, tc.config)
		if tc.wantErr != nil {
			require.ErrorIs(t, err, tc.wantErr, tc.desc)
			continue
		}
		require.NoError(t, err, tc.desc)
		require.Equal(t, tc.want, got)
	}

}
