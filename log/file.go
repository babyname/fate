package log

import (
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

type FileLogger struct {
	*slog.Logger
	file *os.File
}

func NewFileLogger(path string) (*FileLogger, error) {
	dir, _ := filepath.Split(path)
	if dir != "" {
		_ = os.MkdirAll(dir, 0755)
	}
	var fl FileLogger
	var err error
	fl.file, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	fl.Logger = slog.New(opts.NewJSONHandler(fl.file))
	return &fl, nil
}

func (l *FileLogger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}
