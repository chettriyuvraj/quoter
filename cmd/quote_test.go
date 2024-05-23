package cmd

import (
	"bytes"
	"flag"
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
			continue
		}
		require.NoError(t, err, tc.desc)
		require.Equal(t, tc.want, got, tc.desc)
		require.Equal(t, tc.output, buf.String(), tc.desc)

	}
}
