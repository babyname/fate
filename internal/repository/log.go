package repository

import (
	logger "github.com/babyname/fate/log"
)

var log logger.Logger

func Logger(name string) {
	log = logger.WithGroup(name)
}
