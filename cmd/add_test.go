package cmd

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunAddCmdSingleQuote(t *testing.T) {
	tcs := []struct {
		desc   string
		config AddConfig
		err    error
		want   []Quote
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
		buf := bytes.Buffer{}
		quoteStorage := ReadWriteSeekerUtil{ReadSeeker: bytes.NewReader([]byte{})}
		err := runAddCmd(&buf, &quoteStorage, tc.config)
		if tc.err != nil {
			require.Error(t, err, tc.err, tc.desc)
			continue
		}
		require.NoError(t, err, tc.desc)
		testQuoteStorage(t, tc.want, &quoteStorage, tc.desc)
	}

}

/* Tests addition of multiple quotes sequentially, and subsequent JSON file formation*/
func TestRunAddCmdMultiQuote(t *testing.T) {
	defer os.Remove(PERSIST_FILENAME)

	tcs := []struct {
		desc   string
		config AddConfig
		err    error
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

	want := []Quote{}
	quoteStorage := ReadWriteSeekerUtil{ReadSeeker: bytes.NewReader([]byte{})}
	for _, tc := range tcs {
		buf := bytes.Buffer{}
		err := runAddCmd(&buf, &quoteStorage, tc.config)
		require.NoError(t, err, tc.desc)
		want = append(want, tc.want...)
	}
	testQuoteStorage(t, want, &quoteStorage, "test multiple adds")

}

func testQuoteStorage(t *testing.T, want []Quote, quoteStorage io.ReadSeeker, desc string) {
	t.Helper()
	var got []Quote
	/* Quote storages seek pointer may have been moved around, but now we want its entire contents */
	_, err := quoteStorage.Seek(0, 0)
	if err != nil {
		require.NoError(t, err, desc)
	}
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
			output: ADD_USAGE_STRING,
		},
		{
			desc:   "non existent flag",
			args:   []string{"-x"},
			err:    flag.ErrHelp,
			output: ADD_USAGE_STRING,
		},
		{
			desc: "genre flag only",
			args: []string{"-g", "romance", "randomquote"},
			want: AddConfig{genre: "romance", quote: "randomquote"},
		},
		{
			desc:   "genre flag but no genre specifed",
			args:   []string{"-g"},
			err:    flag.ErrHelp,
			output: ADD_USAGE_STRING,
		},
		{
			desc: "no flags",
			args: []string{"randomquote"},
			want: AddConfig{genre: "misc", quote: "randomquote"},
		},
	}

	for _, tc := range tcs {
		buf := bytes.Buffer{}
		got, err := parseAddArgs(&buf, tc.args)
		if tc.err != nil {
			require.Error(t, tc.err, err, tc.desc)
			continue
		}
		require.NoError(t, err, tc.desc)
		require.Equal(t, tc.want, got, tc.desc)
		require.Equal(t, tc.output, buf.String(), tc.desc)

	}
}
