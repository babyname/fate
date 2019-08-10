package fate_test

import (
	"github.com/godcong/fate"
	"testing"
	"time"
)

func TestFate_FirstRunInit(t *testing.T) {
	newFate := fate.NewFate("李", time.Now())
	newFate.FirstRunInit()
}
