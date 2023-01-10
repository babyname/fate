package logger

import (
	"context"
	"golang.org/x/exp/slog"
)

type Logger interface {
	LogDepth(calldepth int, level slog.Level, msg string, args ...any)
	LogAttrsDepth(calldepth int, level slog.Level, msg string, attrs ...slog.Attr)
	Handler() slog.Handler
	Context() context.Context
	With(args ...any) *slog.Logger
	WithGroup(name string) *slog.Logger
	WithContext(ctx context.Context) *slog.Logger
	Enabled(level slog.Level) bool
	Log(level slog.Level, msg string, args ...any)
	LogAttrs(level slog.Level, msg string, attrs ...slog.Attr)
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, err error, args ...any)
}

// func LogDepth(calldepth int, level slog.Level, msg string, args ...any) {
//
// }
// func LogAttrsDepth(calldepth int, level slog.Level, msg string, attrs ...slog.Attr) {
//
// }
//
//	func Handler() slog.Handler {
//		return logger.Handler()
//	}
//
//	func Context() context.Context {
//		return logger.Context()
//	}
//
//	func With(args ...any) *slog.Logger {
//		return logger.With(args...)
//	}

func WithGroup(name string) *slog.Logger {
	return logger.WithGroup(name)
}

//
//func WithContext(ctx context.Context) *slog.Logger {
//
//}
//
//func Enabled(level slog.Level) bool {
//
//}
//
//func Log(level slog.Level, msg string, args ...any) {
//
//}
//
//func LogAttrs(level slog.Level, msg string, attrs ...slog.Attr) {
//
//}
//
//func Debug(msg string, args ...any) {
//
//}
//
//func Info(msg string, args ...any) {
//
//}
//
//func Warn(msg string, args ...any) {
//
//}
//
//func Error(msg string, err error, args ...any) {
//
//}
