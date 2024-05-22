package cmd

import (
	"bytes"
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
)

// func TestHandleAdd(t *testing.T) {
// 	tcs := []struct {
// 		desc   string
// 		args   []string
// 		err    error
// 		output string
// 	}{
// 		{
// 			desc:   "non-persistent quote",
// 			args:   []string{"Humse door jaoge kaise? Humko tum bhulaoge kaise? Hum vo khushbu hai jo saanson me baste hai, apni saanson ko rok paoge kaise?"},
// 			err:    nil,
// 			output: "",
// 		},
// 		{
// 			desc:   "persistent quote with no genre",
// 			args:   []string{"-p", "Humse door jaoge kaise? Humko tum bhulaoge kaise? Hum vo khushbu hai jo saanson me baste hai, apni saanson ko rok paoge kaise?"},
// 			err:    nil,
// 			output: "",
// 		},
// 		{
// 			desc:   "persistent quote with genre",
// 			args:   []string{"-p", "-g", "romance", "Humse door jaoge kaise? Humko tum bhulaoge kaise? Hum vo khushbu hai jo saanson me baste hai, apni saanson ko rok paoge kaise?"},
// 			err:    nil,
// 			output: "",
// 		},
// 		{
// 			desc:   "persistent quote with genre flag but no genre specified",
// 			args:   []string{"-p", "-g", "Humse door jaoge kaise? Humko tum bhulaoge kaise? Hum vo khushbu hai jo saanson me baste hai, apni saanson ko rok paoge kaise?"},
// 			err:    ErrNoGenreSpecified,
// 			output: ``,
// 		},
// 	}
// }

func TestParseAddArgs(t *testing.T) {

	tcs := []struct {
		desc   string
		args   []string
		err    error
		want   AddConfig
		output string
	}{
		{
			desc:   "help flag",
			args:   []string{"-h"},
			err:    flag.ErrHelp,
			output: USAGE_STRING,
		},
		{
			desc:   "non existent flag",
			args:   []string{"-x"},
			err:    flag.ErrHelp,
			output: USAGE_STRING,
		},
		{
			desc: "genre flag only",
			args: []string{"-g", "romance", "randomquote"},
			want: AddConfig{genre: "romance"},
		},
		{
			desc: "no flags",
			args: []string{"randomquote"},
			want: AddConfig{genre: "misc"},
		},
	}

	for _, tc := range tcs {
		buf := bytes.Buffer{}
		got, err := parseAddArgs(&buf, tc.args)
		if tc.err != nil {
			require.Error(t, tc.err, err, tc.desc)
		} else {
			require.NoError(t, err, tc.desc)
			require.Equal(t, tc.want, got, tc.desc)
			require.Equal(t, tc.output, buf.String(), tc.desc)
		}

	}
}
