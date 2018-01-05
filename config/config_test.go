package config_test

import (
	"testing"

	"github.com/godcong/fate/config"
)

func TestLoad(t *testing.T) {
	if config.DefaultConfig().GetString("database.name") == "" {
		t.Error("database.name")
	}

}
