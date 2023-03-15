package main

import (
	"fmt"
	"time"

	"github.com/babyname/fate"
	"github.com/babyname/fate/log"
	"github.com/godcong/chronos"

	"github.com/spf13/cobra"
)

func cmdName() *cobra.Command {
	var born string
	var last string
	var output string
	var sex string
	cmd := &cobra.Command{
		Use:   "name",
		Short: "生成姓名列表",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("start", time.Now().String())
			f, err := fate.New(cfg)
			if err != nil {
				log.Error("new fate", err)
			}
			s := f.NewSessionWithFilter(fate.NewFilter(fate.FilterOption{
				CharacterFilter:     true,
				CharacterFilterType: 0,
				MinCharacter:        3,
				MaxCharacter:        18,
				RegularFilter:       true,
				DaYanFilter:         true,
				WuXingFilter:        true,
				SexFilter:           false,
			}))
			fmt.Println("last", last)
			l, ok := getLastChar(last)
			if !ok {
				fmt.Println("请输入姓氏")
				return
			}
			fmt.Println("born", born)
			b, err := time.Parse(chronos.DateFormat, born)
			if err != nil {
				fmt.Println("请输入正确的出生日期")
				return
			}
			sx := 1
			if sex == "girl" {
				sx = 0
			}
			input := &fate.Input{
				Last: l,
				Born: b,
				Sex:  fate.Sex(sx),
			}
			err = s.Start(input)
			if err != nil {
				fmt.Println("发生了一些错误", err.Error())
				return
			}
			var names []fate.Name
			//for input.Output().NextName() {
			//	names = append(names, s.Name(fn))
			//}
			if s.Err() != nil {
				fmt.Println("输出时发生了一些错误", s.Err().Error())
				return
			}
			<-s.Context().Done()
			fmt.Println("end", time.Now().String())
			time.Sleep(3 * time.Second)
			for i := 0; i < len(names); i += 10 {
				fmt.Print("Name:")
				for j := i; j < len(names) && j < i+10; j++ {
					fmt.Print(names[j], "  |  ")
				}
				fmt.Printf("\r\n")
			}
			fmt.Println("Total", len(names), input.Output().Total())
		},
	}
	cmd.Flags().StringVarP(&last, "last", "l", "", "指定姓氏")
	cmd.Flags().StringVarP(&born, "born", "b", time.Now().Format(chronos.DateFormat), "设置新生儿生日 2016/01/02 15:04")
	cmd.Flags().StringVarP(&sex, "sex", "s", "boy", "设置新生儿性别")
	cmd.Flags().StringVarP(&output, "output", "o", "", "设置输出路径")
	return cmd
}

func getLastChar(s string) ([2]string, bool) {
	var l [2]string
	switch len([]rune(s)) {
	case 1:
		l[0] = string([]rune(s)[0])
		return l, true
	case 2:
		l[0] = string([]rune(s)[0])
		l[1] = string([]rune(s)[1])
	default:
		return l, false
	}
	return l, true
}
