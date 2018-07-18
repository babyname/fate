package model_test

import (
	"log"
	"testing"

	"github.com/godcong/fate/config"
	"github.com/godcong/fate/model"
)

func TestConnectDB(t *testing.T) {
	config := config.Default()
	log.Println(config.GetSub("database"))
	//log.Println(model.ConnectDB(config))
}

func TestCreateTables(t *testing.T) {
	model.CreateTables()
}

func TestSync(t *testing.T) {
	list := []model.SyncAble{
		new(model.FivePhase),
	}

	for _, l := range list {
		err := l.Sync()
		log.Println(err)
	}
}
