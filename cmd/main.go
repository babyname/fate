package main

import (
	"fmt"
	"log"

	"github.com/godcong/fate"
)

func main() {
	fmt.Println("请输入姓: ")
	sur := ""
	fmt.Scanln(&sur)
	log.Println(sur)
	f := fate.NewFactory(sur)

	fmt.Println("是否使用最优五行配置：1. Yes 2. No 其他. Yes")
	five := 0
	fmt.Scanln(&five)
	log.Println(five)

	f.LoadThreeFive(five)

	f.SecondName()
	//fmt.Println("")

	//log.Println(f)
	f.ThirdName()
	log.Println(f.GetName())
}
