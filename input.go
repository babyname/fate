package fate

import (
	"time"
)

type Sex int //girl:0,boy:1

type Input struct {
	Name [2]rune
	Born time.Time
	Sex  Sex
}
