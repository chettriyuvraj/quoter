package cmd

import (
	"bytes"
	"errors"
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseQuoteArgs(t *testing.T) {

	tcs := []struct {
		desc       string
		args       []string
		want       QuoteConfig
		wantErr    error
		wantErrStr string
	}{
		{
			desc:       "help flag",
			args:       []string{"-h"},
			wantErr:    flag.ErrHelp,
			wantErrStr: completeQuoteUsageString + "\n",
		},
		{
			desc:       "non existent flag",
			args:       []string{"-x"},
			wantErr:    errors.New("flag provided but not defined: -x"),
			wantErrStr: "flag provided but not defined: -x" + "\n" + completeQuoteUsageString + "\n",
		},
		{
			desc:       "genre flag but no genre specifed",
			args:       []string{"-g"},
			wantErr:    errors.New("flag needs an argument: -g"),
			wantErrStr: "flag needs an argument: -g" + "\n" + completeQuoteUsageString + "\n",
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
		got, err, gotErrStr := parseQuoteArgs(tc.args)

		if tc.wantErr != nil {
			/* Assert if error strings are the same - error not compared directly because internal errors are also returned which will not match with error.Is */
			require.ErrorContains(t, err, tc.wantErr.Error(), tc.desc)

			/* Compare the error string we are expecting with the one we want */
			require.Equal(t, tc.wantErrStr, gotErrStr, tc.desc)
			continue
		}

		/* Non-error case */
		require.NoError(t, err, tc.desc)
		require.Equal(t, tc.want, got, tc.desc)
		require.Equal(t, tc.wantErrStr, gotErrStr, tc.desc) /* Standard output remains empty */

	}
}

func TestGetRandomQuoteCmd(t *testing.T) {
	jsonData := `[{"text":"if not us, then who? if not now, then when?","genre":"revolution"},{"text":"abki baar bichde toh khwabo me mile..jaise sukhe hue phool kitabo me mile","genre":"romance"},{"text":"if the lessons of history teach us anything, it is that no one learns the lessons that history teaches us","genre":"misc"},{"text":"I'd love to say you make me weak in the knees, but to be quite honest, and completely upfront - you make my body forget it has knees at all","genre":"love"}]`
	quoteStorage := bytes.NewReader([]byte(jsonData))

	tcs := []struct {
		desc    string
		config  QuoteConfig
		want    Quote
		wantErr error
	}{
		{
			desc:   "quote of romance genre",
			config: QuoteConfig{Genre: "romance"},
			want:   Quote{Genre: "romance", Text: "abki baar bichde toh khwabo me mile..jaise sukhe hue phool kitabo me mile"},
		},
		// { TODO: How to test random quote?
		// 	desc:   "quote with no genre specified",
		// 	config: QuoteConfig{Genre: "romance"},
		// 	want: "Humse door jaoge kaise? Humko tum bhulaoge kaise? Hum vo khushbu hai jo saanson me baste hai, apni saanson ko rok paoge kaise?"
		// },
	}

	for _, tc := range tcs {
		got, err := getRandomQuote(quoteStorage, tc.config)
		if tc.wantErr != nil {
			require.ErrorIs(t, err, tc.wantErr, tc.desc)
			continue
		}
		require.NoError(t, err, tc.desc)
		require.Equal(t, tc.want, got)
	}

}
