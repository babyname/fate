package log

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/babyname/fate/config"
)

const (
	envLogFile = "ENV_LOG_FILE"
)

var (
	opts   = slog.HandlerOptions{AddSource: true}
	output *os.File
)

type Logger interface {
	Handler() slog.Handler
	Context() context.Context
	With(args ...any) *slog.Logger
	WithGroup(name string) *slog.Logger
	WithContext(ctx context.Context) *slog.Logger
	Enabled(level slog.Level) bool
	Log(level slog.Level, msg string, args ...any)
	LogDepth(calldepth int, level slog.Level, msg string, args ...any)
	LogAttrs(level slog.Level, msg string, attrs ...slog.Attr)
	LogAttrsDepth(calldepth int, level slog.Level, msg string, attrs ...slog.Attr)
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, err error, args ...any)
}

func init() {
	env, exist := os.LookupEnv(envLogFile)
	if !exist {
		return
	}
	file, err := openLogFile(env)
	if err != nil {
		return
	}
	output = file
	l := slog.New(slog.NewTextHandler(file, &opts))
	slog.SetDefault(l)
}

func openLogFile(path string) (*os.File, error) {
	dir, _ := filepath.Split(path)
	if dir != "" {
		_ = os.MkdirAll(dir, 0755)
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func LoadGlobalConfig(cfg config.LogConfig) error {
	file, err := openLogFile(cfg.Path)
	if err != nil {
		return err
	}

	opts.AddSource = cfg.ShowSource
	var h slog.Handler
	switch cfg.LogType {
	case "text":
		h = slog.NewTextHandler(file, &opts)
	case "json":
		h = slog.NewJSONHandler(file, &opts)
	}

	l := slog.New(h)
	l.Enabled(context.Background(), stringToLevel(cfg.Level))
	slog.SetDefault(l)
	return nil

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
