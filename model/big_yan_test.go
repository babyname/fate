package model

import (
	"log"
	"testing"
)

func TestGetBigYan(t *testing.T) {
	for i := 0; i <= 82; i++ {
		v := GetBigYan(i)
		log.Println(v)
	}
}
