package debug

import "log"

func Println(v... interface{})  {
	log.Println(v)
}

func Print(v... interface{})  {
	log.Print(v)
}