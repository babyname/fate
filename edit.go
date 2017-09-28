package fate

import (
	"strings"

	"log"

	"github.com/godcong/fate/model"
)

func InsertChar(typ, val string, stk int) {
	vs := strings.Split(val, "")
	for _, v := range vs {
		if strings.TrimSpace(v) != "" {
			c := model.Character{
				NameChar: v,
				NameType: typ,
				Strokes:  stk,
			}
			model.ORM().Create(&c)
			log.Println(v)
		}
	}
}
