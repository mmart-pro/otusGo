package main

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

var HideProgress = false // скрывать progressBar (для тестов)

const (
	// какой порцией копируется файл (занижено специально)
	ChunkSize int64 = 133
	// задержка после копирования блока (для процентиков)
	SleepTime = 10
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// источник
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}
	srcSize := srcFileInfo.Size()
	if offset > srcSize {
		return ErrOffsetExceedsFileSize
	}
	_, err = srcFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	// назначение
	dstFile, err := os.OpenFile(toPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// сколько байт нужно скопировать
	total := limit
	if total == 0 || total > srcSize-offset {
		total = srcSize - offset
	}

	var bar *pb.ProgressBar

	if !HideProgress {
		bar = pb.Start64(total)
		bar.Set(pb.Bytes, true)
		defer bar.Finish()
	}

	copied := int64(0)
	for {
		toCopy := total - copied
		if toCopy > ChunkSize {
			toCopy = ChunkSize
		}
		wr, werr := io.CopyN(dstFile, srcFile, toCopy)
		copied += wr

		if !HideProgress {
			bar.SetCurrent(copied)
			time.Sleep(time.Millisecond * SleepTime)
		}

		if werr != nil {
			return werr
		} else if copied == total {
			break
		}
	}

	return nil
}
