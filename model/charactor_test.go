package model_test

import (
	"log"
	"testing"

	"github.com/godcong/fate/model"
)

func TestCharacter_Sync(t *testing.T) {
	c := model.Character{}
	c.Sync()
}

func TestCharacter_List(t *testing.T) {
	var cl []model.Character

	model.CharacterList("火", 7, &cl)
	log.Println(cl)
}
