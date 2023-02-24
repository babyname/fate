package log

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/babyname/fate/config"
	"golang.org/x/exp/slog"
)

const (
	envLogFile = "ENV_LOG_FILE"
)

var (
	opts = slog.HandlerOptions{AddSource: true}
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
	output.File = file
	//unsafeL := (*unsafe.Pointer)(unsafe.Pointer(output.File))
	//atomic.SwapPointer(unsafeL, unsafe.Pointer(file))
	return nil
}

func SetGlobalOutput(f *os.File) {
	output.File = f
	//unsafeL := (*unsafe.Pointer)(unsafe.Pointer(output.File))
	//atomic.SwapPointer(unsafeL, unsafe.Pointer(f))
}

func LoadGlobalConfig(cfg config.LogConfig) error {
	err := openLogFile(cfg.Path)
	if err != nil {
		return err
	}

	opts.AddSource = cfg.ShowSource
	var h slog.Handler = opts.NewJSONHandler(output.File)
	if cfg.LogType != "json" {
		h = opts.NewTextHandler(output.File)
	}

	l := slog.New(h)
	l.Enabled(stringToLevel(cfg.Level))
	output.Logger = l
	//unsafeL := (*unsafe.Pointer)(unsafe.Pointer(output.Logger))
	//atomic.SwapPointer(unsafeL, unsafe.Pointer(l))
	return nil
}

func Default() Logger {
	return output
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
