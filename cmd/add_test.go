package cmd

import (
	"bytes"
	"errors"
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
)

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
		wantErrStr string
	}{
		{
			desc:       "help flag",
			args:       []string{"-h"},
			wantErr:    flag.ErrHelp,
			wantErrStr: completeAddUsageString + "\n",
		},
		{
			desc:       "non existent flag",
			args:       []string{"-x"},
			wantErr:    errors.New("flag provided but not defined: -x"),
			wantErrStr: "flag provided but not defined: -x" + "\n" + completeAddUsageString + "\n",
		},
		{
			desc:       "genre flag but no genre specifed",
			args:       []string{"-g"},
			wantErr:    errors.New("flag needs an argument: -g"),
			wantErrStr: "flag needs an argument: -g" + "\n" + completeAddUsageString + "\n",
		},
		{
			desc:       "no positional args",
			args:       []string{},
			wantErr:    ErrNoPositionalArgs,
			wantErrStr: ErrNoPositionalArgs.Error() + "\n" + completeAddUsageString + "\n",
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
		got, err, gotErrStr := parseAddArgs(tc.args)

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

func TestParseQuotes(t *testing.T) {
	jsonData := `[{"text":"if not us, then who? if not now, then when?","genre":"revolution"},{"text":"abki baar bichde toh khwabo me mile..jaise sukhe hue phool kitabo me mile","genre":"romance"},{"text":"if the lessons of history teach us anything, it is that no one learns the lessons that history teaches us","genre":"misc"},{"text":"I'd love to say you make me weak in the knees, but to be quite honest, and completely upfront - you make my body forget it has knees at all","genre":"love"}]`
	quoteStorage := bytes.NewReader([]byte(jsonData))
	want := []Quote{
		{Text: "if not us, then who? if not now, then when?", Genre: "revolution"},
		{Text: "abki baar bichde toh khwabo me mile..jaise sukhe hue phool kitabo me mile", Genre: "romance"},
		{Text: "if the lessons of history teach us anything, it is that no one learns the lessons that history teaches us", Genre: "misc"},
		{Text: "I'd love to say you make me weak in the knees, but to be quite honest, and completely upfront - you make my body forget it has knees at all", Genre: "love"},
	}
	got, err := parseQuotes(quoteStorage)
	require.NoError(t, err)
	require.Equal(t, want, got)
}
