package cmd

import (
	"fmt"
	"log"

	"github.com/godcong/fate"
	"github.com/godcong/fate/model"
)

func main() {
	fmt.Println("请输入姓: ")
	sur := "蒋"
	//fmt.Scanln(&sur)
	log.Println(sur)

	f := fate.NewFactory(sur)

	log.Println("是否使用最优五行配置：1. Yes 2. No 默认:Yes")
	five := 1
	//fmt.Scanln(&five)
	log.Println(five)

	f.LoadThreeFive(five)
	z := model.FindZodiac(model.ZODIAC_JI)

	f.SecondName(*z)
	//fmt.Println("")

	//log.Println(f)
	f.ThirdName(*z)
	log.Println(f.GetName())
}
