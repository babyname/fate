package regular

import (
	"github.com/babyname/fate"
	"github.com/babyname/fate/config"
	"testing"
)

// TestNew ...
func TestNew(t *testing.T) {
	c := config.LoadConfig()
	db := fate.InitDatabaseWithConfig(*c)
	regular := New(db)
	regular.Run()
}
