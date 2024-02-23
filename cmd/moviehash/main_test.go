package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/opensubtitlescli/moviehash/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintsHelpMessage(t *testing.T) {
	out, code := run("-h")
	assert.Equal(t, 0, code)
	assert.Equal(t, help(), out)
}

func TestPrintsMoviehashVersion(t *testing.T) {
	out, code := run("-v")
	assert.Equal(t, 0, code)
	assert.Equal(t, "0.1.0\n", out)
}

func TestPrintsHelpMessageIfThereAreNoArguments(t *testing.T) {
	out, code := run()
	assert.Equal(t, 2, code)
	assert.Equal(t, help(), out)
}

func TestPrintsHelpMessageIfArgumentsAreInvalid(t *testing.T) {
	out, code := run("-invalid")
	assert.Equal(t, 2, code)
	assert.Equal(t, "flag provided but not defined: -invalid\n" + usage(), out)
}

func TestPrintsTheMoviehashOfTheFile(t *testing.T) {
	m, err := testdata.FilesMap()
	require.NoError(t, err)

	_, err = os.Stat(m[0][0])
	if err != nil && errors.Is(err, os.ErrNotExist) {
		fmt.Printf("warn: %s not found\n", m[0][0])
		return
	}

	out, code := run(m[0][0])
	assert.Equal(t, 0, code)
	assert.Equal(t, m[0][1] + "\n", out)
}

func TestPrintsTheMoviehashOfFiles(t *testing.T) {
	m, err := testdata.FilesMap()
	require.NoError(t, err)

	_, err = os.Stat(m[0][0])
	if err != nil && errors.Is(err, os.ErrNotExist) {
		fmt.Printf("warn: %s not found\n", m[0][0])
		return
	}

	_, err = os.Stat(m[1][0])
	if err != nil && errors.Is(err, os.ErrNotExist) {
		fmt.Printf("warn: %s not found\n", m[1][0])
		return
	}

	out, code := run(m[0][0], m[1][0])
	assert.Equal(t, 0, code)
	assert.Equal(t, m[0][1] + "\n" + m[1][1] + "\n", out)
}

func TestPrintsAnErrorIfTheFileDoesNotExist(t *testing.T) {
	out, code := run("./main.txt")
	assert.Equal(t, 2, code)
	assert.Equal(t, "open ./main.txt: no such file or directory\n" + usage(), out)
}

func TestPrintsAnErrorIfTheHashCannotBeCalculated(t *testing.T) {
	out, code := run("./main.go")
	assert.Equal(t, 2, code)
	m := fmt.Sprintf("moviehash: the %q file is too small\n", "./main.go")
	assert.Equal(t, m + usage(), out)
}

func run(args ...string) (string, int) {
	args = append([]string{"moviehash"}, args...)
	buf := &bytes.Buffer{}
	cmd := new(buf)
	code := cmd.run(args)
	return buf.String(), code
}

func help() string {
	return "moviehash: calculate the moviehash of the file.\n" + usage()
}

func usage() string {
	return "usage: moviehash [-hv] <path...>\n"
}
