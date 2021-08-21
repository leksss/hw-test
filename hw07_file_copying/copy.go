package main

import (
	"errors"
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
		return err
	}

	if !srcStat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > srcStat.Size() {
		return ErrOffsetExceedsFileSize
	}

	iLimit := int(limit)
	bufSize := int(srcStat.Size())

	if limit+offset > srcStat.Size() {
		iLimit = int(srcStat.Size() - offset)
	}

	if iLimit > 0 && iLimit < bufSize {
		bufSize = iLimit
	}

	buf := make([]byte, bufSize)

	reader, err := os.Open(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	writer, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer writer.Close()

	n, err := reader.ReadAt(buf, offset)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}

	if offset > 0 {
		_, err = reader.Seek(offset, 0)
		if err != nil {
			return err
		}
	}

	bar := pb.Full.Start(bufSize)
	barReader := bar.NewProxyReader(reader)

	if _, err := io.CopyN(writer, barReader, int64(n)); err != nil {
		return err
	}

	bar.Finish()

	return nil
}
