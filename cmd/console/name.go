package main

import (
	"bufio"
	"fmt"
	"os"
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
	var filter string
	var sex string
	cmd := &cobra.Command{
		Use:   "name",
		Short: "生成姓名列表",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("start", time.Now().String())
			f, err := fate.New(cfg)
			if err != nil {
				log.Error("new fate", err)
				return
			}
			s := f.NewSessionWithFilter(fate.NewFilter(fate.FilterOption{
				CharacterFilter:     true,
				CharacterFilterType: 0,
				MinStroke:           3,
				MaxStroke:           18,
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
			//for input.Output().NextName() {
			//	names = append(names, s.Name(fn))
			//}
			if s.Err() != nil {
				fmt.Println("输出时发生了一些错误", s.Err().Error())
				return
			}
			<-s.Context().Done()
			fmt.Println("query end", time.Now().String())
			time.Sleep(3 * time.Second)
			if len(filter) > 0 {
				fr := []rune(filter)
				for i := range fr {
					filters := input.Output().Filter(string(fr[i]))
					fmt.Println("文字", string(fr[i]), "过滤了", filters, "个")
				}
			}
			fmt.Println("end", time.Now().String())
			if output != "" {
				fmt.Println("结果将输出到", output)
				err := WriteToFile(input.Output(), output)
				if err == nil {
					return
				}
			}
			PrintScreen(input.Output())
			fmt.Println("Finished")
		},
	}
	cmd.Flags().StringVarP(&last, "last", "l", "", "指定姓氏")
	cmd.Flags().StringVarP(&born, "born", "b", time.Now().Format(chronos.DateFormat), "设置新生儿生日 2016/01/02 15:04")
	cmd.Flags().StringVarP(&sex, "sex", "s", "boy", "设置新生儿性别")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "从结果中过滤掉指定汉字")
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

func PrintScreen(output *fate.Output) {
OutputName:
	for name, ok := output.NextName(); ; {
		fmt.Print("Name:")
		for j := 0; j < 10; j++ {
			if !ok {
				break OutputName
			}
			fmt.Print(name, "  |  ")
			name, ok = output.NextName()
		}
		fmt.Printf("\r\n")
	}
	fmt.Printf("\r\n")
	fmt.Println("Total", output.Total())
}

func WriteToFile(output *fate.Output, path string) error {
	log.Info("open file", "path", path)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer func(w *bufio.Writer) {
		err := w.Flush()
		if err != nil {
			log.Error("flush", err)
		}
	}(w)

OutputName:
	for name, ok := output.NextName(); ; {
		for j := 0; j < 10; j++ {
			if !ok {
				break OutputName
			}
			_, err = w.WriteString(name.String())
			if err != nil {
				return err
			}
			_, err = w.WriteString("(" + name.Strokes() + ")")
			if err != nil {
				return err
			}
			_, err = w.WriteString("  |  ")
			if err != nil {
				return err
			}
			name, ok = output.NextName()
		}
		_, err = w.WriteString("\n")
		if err != nil {
			return err
		}

	}

	_, err = w.WriteString("\n")
	if err != nil {
		return err
	}
	_, err = w.WriteString(fmt.Sprintf("Total %d", output.Total()))
	if err != nil {
		return err
	}
	return nil
}
