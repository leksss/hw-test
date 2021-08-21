package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const OutFile = "out.txt"

func TestCopy(t *testing.T) {
	t.Run("blank file names", func(t *testing.T) {
		err := Copy("", "", 0, 0)
		assert.Equal(t, ErrInvalidFileNames, err)

		t.Cleanup(func() { removeOutFile() })
	})

	t.Run("random dev", func(t *testing.T) {
		err := Copy("/dev/urandom", OutFile, 0, 0)
		assert.Equal(t, ErrUnsupportedFile, err)

		t.Cleanup(func() { removeOutFile() })
	})

	t.Run("limit is bigger than file size", func(t *testing.T) {
		srcFile := "testdata/input.txt"
		srcStat, err := os.Stat(srcFile)
		if err != nil {
			t.Fatal("getting file stat failed")
		}

		err = Copy(srcFile, OutFile, srcStat.Size()+1, 0)
		assert.Equal(t, ErrOffsetExceedsFileSize, err)

		t.Cleanup(func() { removeOutFile() })
	})
}

func removeOutFile() {
	_ = os.Remove(OutFile)
}
