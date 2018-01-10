package model_test

import (
	"testing"

	"github.com/godcong/fate/model"
)

func TestCharacter_Sync(t *testing.T) {
	c := model.Character{}
	c.Sync()
}
