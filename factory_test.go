package fate_test

import (
	"testing"

	"log"

	"github.com/godcong/fate"
)

func TestNewFactory(t *testing.T) {
	f := fate.NewFactory("蒋")
	log.Println(f)
}
