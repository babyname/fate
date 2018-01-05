package model_test

import (
	"log"
	"testing"

	"github.com/godcong/fate/config"
	"github.com/godcong/fate/model"
)

func TestConnectDB(t *testing.T) {
	config := config.DefaultConfig()
	log.Println(config.GetSub("database"))
	//log.Println(model.ConnectDB(config))
}

func TestCreateTables(t *testing.T) {
	model.CreateTables()
}
