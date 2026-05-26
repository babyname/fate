package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	v2 "github.com/godcong/chronos/v2"
	"github.com/babyname/fate"
	"github.com/babyname/fate/internal/analysis"
	"github.com/babyname/fate/internal/database"
	"github.com/babyname/fate/internal/log"
	"github.com/babyname/fate/internal/repository"
	"github.com/babyname/fate/internal/seeddb"

	"github.com/spf13/cobra"
)

func cmdName() *cobra.Command {
	var born string
	var last string
	var output string
	var filter string
	var sex string
	var xiyong string
	var strictness string
	cmd := &cobra.Command{
		Use:   "name",
		Short: "生成姓名列表",
		Args:  cobra.MaximumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			fmt.Println("start", time.Now().String())
			repo := getRepository()
			if repo == nil {
				fmt.Println("数据库未初始化")
				return
			}
			if err := seeddb.EnsureSeeded(context.Background(), repo); err != nil {
				log.Error("ensure seeded", err)
			}
			f := fate.NewWithRepo(cfg, repo)
			s := f.NewSessionWithFilter(fate.NewFilter(fate.FilterOption{
				CharacterFilter:     true,
				CharacterFilterType: 0,
				MinStroke:           3,
				MaxStroke:           18,
				RegularFilter:       true,
				DaYanFilter:         true,
				WuXingFilter:        true,
				SexFilter:           false,
				XiYongMethod:        xiyong,
				FilterStrictness:    strictness,
			}))
			fmt.Println("last", last)
			l, ok := getLastChar(last)
			if !ok {
				fmt.Println("请输入姓氏")
				return
			}
			bornStr := born
			if len(args) > 0 {
				bornStr = born + " " + args[0]
			}
			fmt.Println("born", bornStr)
			b, err := parseBornTime(bornStr)
			if err != nil {
				fmt.Println("请输入正确的出生日期，格式: 2024/06/15 10:30 或 2024/06/15")
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
			if s.Err() != nil {
				fmt.Println("输出时发生了一些错误", s.Err().Error())
				return
			}
			<-s.Context().Done()
			fmt.Println("query end", time.Now().String())
			time.Sleep(3 * time.Second)
			if len(filter) > 0 {
				fmt.Println("注意: 过滤功能已迁移到API层，请使用 qiming 服务")
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
	cmd.Flags().StringVarP(&last, "surname", "s", "", "指定姓氏")
	cmd.Flags().StringVarP(&born, "born", "b", time.Now().Format(v2.DateFormatYMDHMS), "设置新生儿生日 2016/01/02 15:04")
	cmd.Flags().StringVarP(&sex, "gender", "g", "boy", "设置新生儿性别(boy/girl)")
	cmd.Flags().StringVarP(&filter, "filter", "f", "", "从结果中过滤掉指定汉字")
	cmd.Flags().StringVarP(&output, "output", "o", "", "设置输出路径")
	cmd.Flags().StringVar(&xiyong, "xiyong", "balance", "喜用神算法: balance(平衡用神法) 或 geju(格局用神法)")
	cmd.Flags().StringVar(&strictness, "strictness", "moderate", "五格筛选严格度: strict, moderate, relaxed")

	cmd.AddCommand(cmdNameDetail())
	return cmd
}

func cmdNameDetail() *cobra.Command {
	var born string
	var last string
	var sex string
	cmd := &cobra.Command{
		Use:   "detail",
		Short: "查看名字详细分析",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			char1 := args[0]
			char2 := args[1]

			fmt.Println("last", last)
			l, ok := getLastChar(last)
			if !ok {
				fmt.Println("请输入姓氏")
				return
			}
			b, err := parseBornTime(born)
			if err != nil {
				fmt.Println("请输入正确的出生日期，格式: 2024/06/15 10:30 或 2024/06/15")
				return
			}
			sx := 1
			if sex == "girl" {
				sx = 0
			}

			db := getRepository()
			if db == nil {
				fmt.Println("数据库未初始化")
				return
			}

			fateData, err := v2.GetFateData(&v2.FateInput{
				BirthDate: b,
				Gender:    sx,
				Surname:   l[0],
			})
			if err != nil {
				fmt.Println("获取八字信息失败", err)
				return
			}

			ctx := cmd.Context()
			c1, err := db.GetCharacter(ctx, repository.Char(char1))
			if err != nil {
				fmt.Println("查询字「"+char1+"」失败:", err)
				return
			}
			c2, err := db.GetCharacter(ctx, repository.Char(char2))
			if err != nil {
				fmt.Println("查询字「"+char2+"」失败:", err)
				return
			}

			l1 := 0
			l2 := 0
			ln, err := db.QueryLastName(ctx, l)
			if err == nil && ln[0] != nil {
				l1 = ln[0].ScienceStroke
				if ln[1] != nil {
					l2 = ln[1].ScienceStroke
				}
			}

			nr := analysis.BuildNameResult(0, l[0], c1, c2, l1, l2, fateData)
			printNameDetail(nr)
		},
	}
	cmd.Flags().StringVarP(&last, "surname", "s", "", "指定姓氏")
	cmd.Flags().StringVarP(&born, "born", "b", time.Now().Format(v2.DateFormatYMDHMS), "设置新生儿生日 2016/01/02 15:04")
	cmd.Flags().StringVarP(&sex, "gender", "g", "boy", "设置新生儿性别(boy/girl)")
	return cmd
}

func getRepository() *repository.Repository {
	b := database.New(cfg.Database)
	cli, err := b.Client()
	if err != nil {
		return nil
	}
	return repository.New(cli)
}

func printNameDetail(nr analysis.NameResult) {
	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("  %s  评分: %.1f  等级: %s\n", nr.FullName, nr.Score, nr.Grade)
	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("  姓氏: %s  名字: %s  笔画: %s\n", nr.Surname, nr.FirstName, nr.Strokes)
	fmt.Println()
	fmt.Println("  【字义分析】")
	fmt.Printf("  %s: 拼音%s 五行属%s 笔画(简/繁/科/康): %d/%d/%d/%d\n",
		nr.Char1.Char, nr.Char1.Pinyin, nr.Char1.WuXing,
		nr.Char1.SimplifiedStroke, nr.Char1.TraditionalStroke,
		nr.Char1.ScienceStroke, nr.Char1.KangxiStroke)
	if nr.Char1.Meaning != "" {
		fmt.Printf("    释义: %s\n", nr.Char1.Meaning)
	}
	if nr.Char1.IsXiYong {
		fmt.Printf("    ★ 喜用神\n")
	}
	fmt.Printf("  %s: 拼音%s 五行属%s 笔画(简/繁/科/康): %d/%d/%d/%d\n",
		nr.Char2.Char, nr.Char2.Pinyin, nr.Char2.WuXing,
		nr.Char2.SimplifiedStroke, nr.Char2.TraditionalStroke,
		nr.Char2.ScienceStroke, nr.Char2.KangxiStroke)
	if nr.Char2.Meaning != "" {
		fmt.Printf("    释义: %s\n", nr.Char2.Meaning)
	}
	if nr.Char2.IsXiYong {
		fmt.Printf("    ★ 喜用神\n")
	}
	fmt.Println()
	if nr.WuGe != nil {
		fmt.Println("  【五格分析】")
		printGeItem("天格", nr.WuGe.TianGe)
		printGeItem("人格", nr.WuGe.RenGe)
		printGeItem("地格", nr.WuGe.DiGe)
		printGeItem("外格", nr.WuGe.WaiGe)
		printGeItem("总格", nr.WuGe.ZongGe)
		fmt.Println()
	}
	fmt.Printf("  三才: %s  吉凶: %s\n", nr.SanCai, nr.SanCaiLuck)
	if nr.SanCaiDetail != "" {
		fmt.Printf("  三才详解: %s\n", nr.SanCaiDetail)
	}
	fmt.Println()
	fmt.Printf("  基础运: %s\n", nr.JiChuYun)
	fmt.Printf("  成功运: %s\n", nr.ChengGongYun)
	fmt.Printf("  人际关系: %s\n", nr.RenJiGuanXi)
	fmt.Println()
	if nr.ZhouYi != nil {
		fmt.Println("  【周易卦象】")
		fmt.Printf("  本卦: %s（%s）\n", nr.ZhouYi.BenGuaName, nr.ZhouYi.BenGuaJiXiong)
		fmt.Printf("  变卦: %s\n", nr.ZhouYi.BianGuaName)
		if nr.ZhouYi.DongYaoDesc != "" {
			fmt.Printf("  动爻: %s\n", nr.ZhouYi.DongYaoDesc)
		}
		if nr.ZhouYi.DaXiang != "" {
			fmt.Printf("  大象: %s\n", nr.ZhouYi.DaXiang)
		}
		fmt.Println()
	}
	if nr.Bazi != nil {
		fmt.Println("  【八字分析】")
		fmt.Printf("  四柱: %s\n", strings.Join(nr.Bazi.FourPillars[:], " "))
		fmt.Printf("  五行: %s\n", strings.Join(nr.Bazi.FiveElements[:], " "))
		fmt.Printf("  纳音: %s\n", strings.Join(nr.Bazi.NaYin[:], " "))
		fmt.Printf("  生肖: %s  星座: %s\n", nr.Bazi.Zodiac, nr.Bazi.Constellation)
		fmt.Println()
	}
	if nr.WuXing != nil {
		fmt.Println("  【五行喜用神】")
		fmt.Printf("  日主: %s(%s)  八字强弱: %s\n", nr.WuXing.DayGan, nr.WuXing.DayWuxing, nr.WuXing.Strength)
		fmt.Printf("  喜用神: %s\n", strings.Join(nr.WuXing.FavorableElements, "、"))
		fmt.Printf("  忌神: %s\n", strings.Join(nr.WuXing.UnfavorableElements, "、"))
		if nr.WuXing.UsefulElement != "" {
			fmt.Printf("  用神: %s\n", nr.WuXing.UsefulElement)
		}
		if len(nr.WuXing.HostileElements) > 0 {
			fmt.Printf("  仇神: %s\n", strings.Join(nr.WuXing.HostileElements, "、"))
		}
		if len(nr.WuXing.IdleElements) > 0 {
			fmt.Printf("  闲神: %s\n", strings.Join(nr.WuXing.IdleElements, "、"))
		}
		fmt.Printf("  算法: %s\n", nr.WuXing.MethodName)
		if nr.WuXing.Analysis != "" {
			fmt.Printf("  分析: %s\n", nr.WuXing.Analysis)
		}
		fmt.Println()
	}
	fmt.Println("  【评分明细】")
	sd := nr.ScoreDetail
	fmt.Printf("  文化印象: %.1f  五行八字: %.1f  生肖: %.1f  五格数理: %.1f  音韵: %.1f\n",
		sd.WenHuaYinXiang, sd.WuXingBaZi, sd.ShengXiao, sd.WuGeShuLi, sd.YinYun)
	fmt.Println()
	fmt.Println("  【综合解读】")
	fmt.Printf("  %s\n", nr.Interpret)
	fmt.Println("═══════════════════════════════════════")
}

func printGeItem(name string, ge analysis.GeItem) {
	fmt.Printf("  %s: %d画  吉凶:%s  阴阳五行:%s\n", name, ge.Stroke, ge.Lucky, ge.YinYangWuXing)
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

func parseBornTime(s string) (time.Time, error) {
	formats := []string{
		v2.DateFormatYMDHMS,
		"2006/01/02 15:04",
		"2006/01/02 15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006/01/02",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("无法解析日期: %s", s)
}

func PrintScreen(output *fate.Output) {
	topNames := output.TopNames()
	allNames := output.AllNames()

	fmt.Println()
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("           Top 10 推荐名字             ")
	fmt.Println("═══════════════════════════════════════")
	for _, nr := range topNames {
		fmt.Printf("  #%d  %s  评分:%.1f  等级:%s  五行:%s%s\n",
			nr.Rank, nr.FullName, nr.Score, nr.Grade,
			nr.Char1.WuXing, nr.Char2.WuXing)
		if nr.WuGe != nil {
			fmt.Printf("      天格:%s 人格:%s 地格:%s 总格:%s\n",
				nr.WuGe.TianGe.Lucky, nr.WuGe.RenGe.Lucky,
				nr.WuGe.DiGe.Lucky, nr.WuGe.ZongGe.Lucky)
		}
		if nr.SanCai != "" {
			fmt.Printf("      三才:%s(%s)\n", nr.SanCai, nr.SanCaiLuck)
		}
		if nr.Interpret != "" {
			fmt.Printf("      %s\n", nr.Interpret)
		}
		fmt.Println()
	}

	fmt.Println("───────────────────────────────────────")
	fmt.Printf("  全部候选名字共 %d 个\n", len(allNames))
	showCount := 20
	if showCount > len(allNames) {
		showCount = len(allNames)
	}
	if showCount > 0 {
		fmt.Println()
		for i := 0; i < showCount; i++ {
			sn := allNames[i]
			surname := ""
			basic := output.Basic()
			if basic.LastName[0] != nil {
				surname = basic.LastName[0].Char
			}
			fmt.Printf("  %s%s%s  %.1f(%s)  |  ",
				surname, sn.Name[0].Char, sn.Name[1].Char,
				sn.Score, sn.Grade)
			if (i+1)%5 == 0 {
				fmt.Println()
			}
		}
		if showCount%5 != 0 {
			fmt.Println()
		}
		if len(allNames) > 20 {
			fmt.Printf("  ... 还有 %d 个名字未显示\n", len(allNames)-20)
		}
	}
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("Total", output.Total())
}

func WriteToFile(output *fate.Output, path string) error {
	log.Info("open file", "path", path)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	w := bufio.NewWriter(f)
	defer func(w *bufio.Writer) {
		err := w.Flush()
		if err != nil {
			log.Error("flush", err)
		}
	}(w)

	topNames := output.TopNames()
	allNames := output.AllNames()

	_, err = w.WriteString("═══════════════════════════════════════\n")
	if err != nil {
		return err
	}
	_, err = w.WriteString("           Top 10 推荐名字             \n")
	if err != nil {
		return err
	}
	_, err = w.WriteString("═══════════════════════════════════════\n")
	if err != nil {
		return err
	}

	for _, nr := range topNames {
		_, err = fmt.Fprintf(w, "#%d %s 评分:%.1f 等级:%s\n",
			nr.Rank, nr.FullName, nr.Score, nr.Grade)
		if err != nil {
			return err
		}
		if nr.WuGe != nil {
			_, err = fmt.Fprintf(w, "  天格:%s 人格:%s 地格:%s 总格:%s\n",
				nr.WuGe.TianGe.Lucky, nr.WuGe.RenGe.Lucky,
				nr.WuGe.DiGe.Lucky, nr.WuGe.ZongGe.Lucky)
			if err != nil {
				return err
			}
		}
		if nr.Interpret != "" {
			_, err = fmt.Fprintf(w, "  %s\n", nr.Interpret)
			if err != nil {
				return err
			}
		}
		_, err = w.WriteString("\n")
		if err != nil {
			return err
		}
	}

	_, err = w.WriteString("───────────────────────────────────────\n")
	if err != nil {
		return err
	}

	surname := ""
	basic := output.Basic()
	if basic.LastName[0] != nil {
		surname = basic.LastName[0].Char
	}

	for i, sn := range allNames {
		_, err = fmt.Fprintf(w, "%s%s%s(%.1f/%s)",
			surname, sn.Name[0].Char, sn.Name[1].Char, sn.Score, sn.Grade)
		if err != nil {
			return err
		}
		if (i+1)%10 == 0 {
			_, err = w.WriteString("\n")
		} else {
			_, err = w.WriteString("  |  ")
		}
		if err != nil {
			return err
		}
	}
	if len(allNames)%10 != 0 {
		_, err = w.WriteString("\n")
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(w, "\nTotal %d\n", output.Total())
	if err != nil {
		return err
	}
	return nil
}
