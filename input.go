package fate

import (
	"time"
)

type Sex int //girl:0,boy:1

type Input struct {
	Last [2]string
	Born time.Time
	Sex  Sex
}
