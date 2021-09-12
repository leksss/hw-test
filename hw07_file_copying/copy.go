package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrInvalidFileNames      = errors.New("invalid file names")
)

func Copy(src, dst string, offset, limit int64) error {
	if src == "" || dst == "" {
		return ErrInvalidFileNames
	}

	srcStat, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("os.Stat: %w", err)
	}

	if !srcStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > srcStat.Size() {
		return ErrOffsetExceedsFileSize
	}

	bufSize := srcStat.Size()

	if limit > srcStat.Size()-offset {
		limit = srcStat.Size() - offset
	}

	if limit > 0 && limit < bufSize {
		bufSize = limit
	}

	reader, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer reader.Close()

	writer, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}
	defer writer.Close()

	n, err := reader.ReadAt(make([]byte, bufSize), offset)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return fmt.Errorf("reader.ReadAt: %w", err)
		}
	}

	if offset > 0 {
		_, err = reader.Seek(offset, 0)
		if err != nil {
			return fmt.Errorf("reader.Seek: %w", err)
		}
	}

	bar := pb.Full.Start64(bufSize)
	barReader := bar.NewProxyReader(reader)

	if _, err := io.CopyN(writer, barReader, int64(n)); err != nil {
		return fmt.Errorf("io.CopyN: %w", err)
	}

	bar.Finish()

	return nil
}
