package model_test

import (
	"testing"

	"github.com/godcong/fate/model"
)

func TestORM(t *testing.T) {
	model.RunMigrate()
}
