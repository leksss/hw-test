package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("empty input dir", func(t *testing.T) {
		_, err := ReadDir("")
		assert.Error(t, err)
	})

	t.Run("validate env name", func(t *testing.T) {
		assert.False(t, isValidEnvName("TTT=YYY"))
		assert.False(t, isValidEnvName(""))
		assert.True(t, isValidEnvName("TTTYYYaa"))
	})

	tests := []struct {
		input    string
		expected string
	}{
		{input: "dd\t", expected: "dd"},
		{input: "dd" + "\x00" + "dd", expected: "dd\ndd"},
		{input: `"hello"`, expected: `"hello"`},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result := filterText(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}

	t.Run("open file", func(t *testing.T) {
		_, err := readFileFirstLine("", "")
		assert.Error(t, err)

		_, err = readFileFirstLine("./testdata", "env")
		assert.Error(t, err)

		_, err = readFileFirstLine("./testdata/env", "")
		assert.Error(t, err)

		_, err = readFileFirstLine("./testdata/env", "BAR")
		assert.NoError(t, err)
	})

	t.Run("read first line", func(t *testing.T) {
		line, err := readFileFirstLine("./testdata/env", "BAR")
		assert.NoError(t, err)
		assert.Equal(t, "bar", line)

		line, err = readFileFirstLine("./testdata/env", "EMPTY")
		assert.NoError(t, err)
		assert.Equal(t, "", line)
	})
}
