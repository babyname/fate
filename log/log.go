package log

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/babyname/fate/config"
)

var (
	sourceAttr func() slog.Attr
)

func Logger(name string, attr ...any) *slog.Logger {
	if sourceAttr != nil {
		attr = append(attr, sourceAttr())
	}
	return slog.With(attr...).WithGroup(name)
}

func retSourceAttr() slog.Attr {
	_, file, line, _ := runtime.Caller(1)
	return slog.Group("source",
		slog.String("file", file),
		slog.Int("line", line),
	)
}

func SetGlobalLogger(cfg config.LogConfig) error {
	if cfg.ShowSource {
		sourceAttr = retSourceAttr
	}
	file := openLogFile(cfg.Path)

	opts := slog.HandlerOptions{
		Level: stringToLevel(cfg.Level),
	}
	var h slog.Handler
	switch cfg.LogType {
	default:
		h = slog.NewTextHandler(file, &opts)
	case "json":
		h = slog.NewJSONHandler(file, &opts)
	}
	slog.SetDefault(slog.New(h))
	return nil
}

func openLogFile(path string) *os.File {
	dir, _ := filepath.Split(path)
	if dir != "" {
		_ = os.MkdirAll(dir, 0755)
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		slog.Error("failed to open log file", "error", err)
		return os.Stderr
	}
	return file
}

func stringToLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return -4
	case "INFO":
		return 0
	case "WARN":
		return 4
	default:
		return 8
	}
}
