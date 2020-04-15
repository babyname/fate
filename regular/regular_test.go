package regular

import (
	"github.com/godcong/fate"
	"github.com/godcong/fate/config"
	"testing"
)

// TestNew ...
func TestNew(t *testing.T) {
	c := config.LoadConfig()
	db := fate.InitDatabaseWithConfig(*c)
	regular := New(db)
	regular.Run()
}
