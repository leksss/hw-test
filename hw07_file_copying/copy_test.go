package main

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const outFile = "copy_test_out.txt"

func TestCopy(t *testing.T) {
	t.Run("blank file names", func(t *testing.T) {
		err := Copy("", "", 0, 0)
		assert.Equal(t, ErrInvalidFileNames, err)
	})

	t.Run("random dev", func(t *testing.T) {
		f := createOutFile()
		t.Cleanup(func() { removeOutFile(f) })

		err := Copy("/dev/urandom", f.Name(), 0, 0)
		assert.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("limit is bigger than file size", func(t *testing.T) {
		f := createOutFile()
		t.Cleanup(func() { removeOutFile(f) })

		srcFile := "testdata/input.txt"
		srcStat, err := os.Stat(srcFile)
		require.NoError(t, err)

		err = Copy(srcFile, f.Name(), srcStat.Size()+1, 0)
		assert.Equal(t, ErrOffsetExceedsFileSize, err)
	})
}

func removeOutFile(file *os.File) {
	_ = os.Remove(file.Name())
}

func createOutFile() *os.File {
	f, err := os.CreateTemp("", outFile)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
