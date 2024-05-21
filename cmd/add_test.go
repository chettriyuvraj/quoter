package cmd

import (
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
// 			desc: "add a non-persistent quote",
// 			args: []string{"Humse door jaoge kaise? Humko tum bhulaoge kaise? Hum vo khushbu hai jo saanson me baste hai, apni saanson ko rok paoge kaise?"},
// 			err:    nil,
// 			output: "",
// 		},
// 		{},
// 	}
// }

/* Not testing unknown flags error as it is equivalent to testing flag standard library */
func TestParseAddArgs(t *testing.T) {
	tcs := []struct {
		desc string
		args []string
		err  error
		want AddConfig
	}{
		{
			desc: "persist flag only",
			args: []string{"-p", "randomquote"},
			want: AddConfig{persist: true, genre: "misc"},
		},
		{
			desc: "genre flag only",
			args: []string{"-g", "romance", "randomquote"},
			want: AddConfig{genre: "romance"},
		},
		{
			desc: "genre and persist flag",
			args: []string{"-g", "romance", "-p", "randomquote"},
			want: AddConfig{persist: true, genre: "romance"},
		},
		{
			desc: "no flags",
			args: []string{"randomquote"},
			want: AddConfig{genre: "misc"},
		},
	}

	for _, tc := range tcs {
		fs := flag.NewFlagSet("add", flag.ContinueOnError)
		got, err := parseAddArgs(fs, tc.args)
		if tc.err != nil {
			require.Error(t, tc.err, err)
		} else {
			require.NoError(t, err)
		}
		require.Equal(t, tc.want, got)
	}
}
