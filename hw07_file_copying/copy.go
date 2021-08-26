package main

import (
	"errors"
	"io"
	"os"
	"path"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile         = errors.New("unsupported file")
	ErrOffsetExceedsFileSize   = errors.New("offset exceeds file size")
	ErrFromPathEmpty           = errors.New("fromPath is empty")
	ErrToPathEmpty             = errors.New("toPath is empty")
	ErrFromPathToPathIdentical = errors.New("fromPath and toPath is identical")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" {
		return ErrFromPathEmpty
	}
	if toPath == "" {
		return ErrToPathEmpty
	}
	if fromPath == toPath {
		return ErrFromPathToPathIdentical
	}

	fstat, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if fstat.IsDir() {
		return ErrUnsupportedFile
	}
	if offset > fstat.Size() {
		return ErrOffsetExceedsFileSize
	}

	from, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer from.Close()

	if _, err := from.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	if limit == 0 {
		limit = fstat.Size()
	}

	targetDir := path.Dir(toPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	to, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer to.Close()

	bar := pb.Full.Start64(limit)
	bar.Set(pb.SIBytesPrefix, true)
	defer bar.Finish()

	reader := bar.NewProxyReader(io.LimitReader(from, limit))
	if _, err := io.Copy(to, reader); err != nil {
		return err
	}

	return nil
}
