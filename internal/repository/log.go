package repository

import (
	logger "github.com/babyname/fate/internal/log"
)

var _ logger.Logger

// Logger 设置包级别的日志记录器。
func Logger(name string) {
	_ = logger.WithGroup(name)
}
