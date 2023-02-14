package log

import (
	"context"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

const (
	envLogFile = "ENV_LOG_FILE"
)

var (
	output         = os.Stderr
	opts           = slog.HandlerOptions{AddSource: true}
	WipeData       = false
	WipeDataLength = 1024
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
	err := openLogFile(env)
	if err != nil {
		return
	}
	defaultLogger = Default()
}

func openLogFile(path string) error {
	dir, _ := filepath.Split(path)
	if dir != "" {
		_ = os.MkdirAll(dir, 0755)
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	output = file
	return nil
}

func SetGlobalLogger(l Logger) {
	defaultLogger = l
	initialized.Store(true)
}

func Default() *slog.Logger {
	return slog.New(opts.NewJSONHandler(output))
}
