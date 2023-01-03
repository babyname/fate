package logger

import (
	"github.com/babyname/fate/config"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"path/filepath"
)

func New(cfg config.Logger) *slog.Logger {
	output := io.Writer(os.Stderr)
	if cfg.Path != "" {
		file, err := openPathFile(cfg.Path)
		if err == nil {
			output = file
		}
	}
	opts := slog.HandlerOptions{AddSource: cfg.ShowSource}

	logger := slog.New(opts.NewJSONHandler(output))
	if cfg.Level != "" {
		logger.Enabled(enableLevel(cfg.Level))
	}
	return logger
}

func openPathFile(path string) (*os.File, error) {
	dir, _ := filepath.Split(path)
	if dir != "" {
		_ = os.MkdirAll(dir, 0755)
	}
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
}

func enableLevel(l string) slog.Level {
	switch l {
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelDebug
	default:
		return slog.LevelDebug
	}
}
