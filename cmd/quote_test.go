package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseQuoteArgs(t *testing.T) {

	tcs := []struct {
		desc   string
		args   []string
		err    error
		want   QuoteConfig
		output string
	}{
		{
			desc:   "help flag",
			args:   []string{"-h"},
			err:    flag.ErrHelp,
			output: QUOTE_USAGE_STRING,
		},
		{
			desc:   "non existent flag",
			args:   []string{"-x"},
			err:    flag.ErrHelp,
			output: QUOTE_USAGE_STRING,
		},
		{
			desc: "genre flag only",
			args: []string{"-g", "romance"},
			want: QuoteConfig{Genre: "romance"},
		},
		{
			desc:   "genre flag but no genre specifed",
			args:   []string{"-g"},
			err:    flag.ErrHelp,
			output: QUOTE_USAGE_STRING,
		},
		{
			desc: "no flags",
			args: []string{""},
			want: QuoteConfig{Genre: ""},
		},
	}

	for _, tc := range tcs {
		buf := bytes.Buffer{}
		got, err := parseQuoteArgs(&buf, tc.args)
		if tc.err != nil {
			require.Error(t, tc.err, err, tc.desc)
			errWantBuf := bytes.Buffer{}
			if err != flag.ErrHelp {
				fmt.Fprint(&errWantBuf, err.Error())
				fmt.Fprintln(&errWantBuf)
			}
			fmt.Fprintln(&errWantBuf, completeQuoteUsageString)
			require.Equal(t, errWantBuf.String(), buf.String(), tc.desc)
			continue
		}
		require.NoError(t, err, tc.desc)
		require.Equal(t, tc.want, got, tc.desc)
		require.Equal(t, tc.output, buf.String(), tc.desc)

	}
}

func addQuotes(t *testing.T, quoteStorage io.ReadWriteSeeker) {
	t.Helper()

	quoteConfigs := []AddConfig{
		{genre: "misc", quote: "Phool hu gulab ka, chameli ka mat samajhna..Aashiq hu aapka apni saheli ka mat samajhna!"},
		{genre: "romance", quote: "Humse door jaoge kaise? Humko tum bhulaoge kaise? Hum vo khushbu hai jo saanson me baste hai, apni saanson ko rok paoge kaise?"},
	}
	for _, config := range quoteConfigs {
		err := runAddCmd(quoteStorage, config)
		require.NoError(t, err)
	}

}

func TestRunQuoteCmd(t *testing.T) {
	/* Add quotes to a storage first */
	quoteStorage := ReadWriteSeekerUtil{ReadSeeker: bytes.NewReader([]byte{})}
	addQuotes(t, &quoteStorage)

	tcs := []struct {
		desc   string
		config QuoteConfig
		err    error
		want   string
	}{
		{
			desc:   "quote of romance genre",
			config: QuoteConfig{Genre: "romance"},
			want:   "Humse door jaoge kaise? Humko tum bhulaoge kaise? Hum vo khushbu hai jo saanson me baste hai, apni saanson ko rok paoge kaise?\n",
		},
		// { TODO: How to test random quote?
		// 	desc:   "quote with no genre specified",
		// 	config: QuoteConfig{Genre: "romance"},
		// 	want: "Humse door jaoge kaise? Humko tum bhulaoge kaise? Hum vo khushbu hai jo saanson me baste hai, apni saanson ko rok paoge kaise?"
		// },
	}

	for _, tc := range tcs {
		buf := bytes.Buffer{}
		err := runQuoteCmd(&buf, &quoteStorage, tc.config)
		if tc.err != nil {
			require.Error(t, err, tc.err, tc.desc)
			continue
		}
		require.NoError(t, err, tc.desc)
		require.Equal(t, tc.want, buf.String())
	}

}
