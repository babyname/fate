package regular

import (
	"testing"

	"github.com/godcong/fate"
	"github.com/godcong/fate/config"
)

// TestNew ...
func TestNew(t *testing.T) {
	c := config.LoadConfig()
	db := fate.InitDatabaseWithConfig(*c)
	regular := New(db)
	regular.Run()
}
