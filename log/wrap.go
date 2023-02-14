package log

import (
	"context"
	"sync/atomic"

	"golang.org/x/exp/slog"
)

var defaultLogger Logger = Default()
var initialized = &atomic.Bool{}

func ini() {
	//defaultLogger =
}

func IsInitialized() bool {
	return initialized.Load()
}

// Handler ...
// @return slog.Handler
func Handler() slog.Handler {
	return defaultLogger.Handler()
}

// Context ...
// @return context.Context
func Context() context.Context {
	return defaultLogger.Context()
}

// With ...
// @param ...any
// @return *slog.Logger
func With(args ...any) *slog.Logger {
	return defaultLogger.With(args...)
}

// WithGroup ...
// @param string
// @return *slog.Logger
func WithGroup(name string) *slog.Logger {
	return defaultLogger.WithGroup(name)
}

// WithContext ...
// @param context.Context
// @return *slog.Logger
func WithContext(ctx context.Context) *slog.Logger {
	return defaultLogger.WithContext(ctx)
}

// Enabled ...
// @param slog.Level
// @return bool
func Enabled(level slog.Level) bool {
	return defaultLogger.Enabled(level)
}

// LogDepth ...
// @param int
// @param slog.Level
// @param string
// @param ...any
func LogDepth(calldepth int, level slog.Level, msg string, args ...any) {
	defaultLogger.LogDepth(calldepth, level, msg, args...)
}

// LogAttrsDepth ...
// @param int
// @param slog.Level
// @param string
// @param ...slog.Attr
func LogAttrsDepth(calldepth int, level slog.Level, msg string, attrs ...slog.Attr) {
	defaultLogger.LogAttrsDepth(calldepth, level, msg, attrs...)
}

// Debug ...
// @param string
// @param ...any
func Debug(msg string, args ...any) {
	defaultLogger.LogDepth(1, slog.LevelDebug, msg, args...)
}

// Info ...
// @param string
// @param ...any
func Info(msg string, args ...any) {
	defaultLogger.LogDepth(1, slog.LevelInfo, msg, args...)
}

// Warn ...
// @param string
// @param ...any
func Warn(msg string, args ...any) {
	defaultLogger.LogDepth(1, slog.LevelWarn, msg, args...)
}

// Error calls Logger.Error on the default logger.
func Error(msg string, err error, args ...any) {
	defaultLogger.Error(msg, err, args...)
}

// Log ...
// @param slog.Level
// @param string
// @param ...any
func Log(level slog.Level, msg string, args ...any) {
	defaultLogger.LogDepth(1, level, msg, args...)
}

// LogAttrs ...
// @param Level
// @param string
// @param ...Attr
func LogAttrs(level slog.Level, msg string, attrs ...slog.Attr) {
	defaultLogger.LogAttrsDepth(1, level, msg, attrs...)
}
