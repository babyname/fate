package main

import (
	"fmt"
	"time"

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
			fmt.Println("end", time.Now().String())
		},
	}
	cmd.Flags().StringVarP(&last, "last", "l", "", "指定名占位: 名或*名")
	cmd.Flags().StringVarP(&born, "born", "b", time.Now().Format(chronos.DateFormat), "设置新生儿生日 2016/01/02 15:04")
	cmd.Flags().StringVarP(&sex, "sex", "s", "boy", "设置新生儿性别")
	cmd.Flags().StringVarP(&output, "output", "o", "", "设置输出路径")
	return cmd
}
