package fate

import (
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/babyname/fate/config"
	"github.com/babyname/fate/internal/rating"
	"github.com/babyname/fate/internal/wuge"
	"github.com/babyname/fate/internal/wuxing"
	v2 "github.com/godcong/chronos/v2"
)

type nameAnalysis struct {
	name       Name
	rating     *rating.NameRating
	tianGe     int
	renGe      int
	diGe       int
	waiGe      int
	zongGe    int
	sanCaiStr  string
	sanCaiLuck string
}

var badChars = map[string]bool{
	"不": true, "丧": true, "衰": true, "蠢": true, "鼠": true,
	"嫌": true, "疮": true, "惨": true, "惭": true, "愁": true,
	"病": true, "死": true, "贫": true, "灾": true, "凶": true,
	"恶": true, "毒": true, "杀": true, "伤": true, "残": true,
	"废": true, "孤": true, "寡": true, "贱": true, "奴": true,
	"鬼": true, "魔": true, "妖": true, "邪": true, "贼": true,
	"偷": true, "骗": true, "叛": true, "伪": true, "奸": true,
	"劣": true, "笨": true, "傻": true, "呆": true, "痴": true,
	"疯": true, "癫": true, "狂": true, "躁": true, "怒": true,
	"恨": true, "怨": true, "冤": true, "苦": true, "痛": true,
	"哭": true, "泪": true, "忧": true, "惧": true, "怕": true,
	"恐": true, "慌": true, "踩": true, "皱": true, "嗽": true,
	"厂": true, "凳": true, "仆": true, "辜": true, "创": true,
	"雹": true, "绑": true, "粥": true, "帽": true, "牌": true,
	"伞": true, "亩": true, "币": true, "弊": true, "幅": true,
	"划": true, "射": true, "叹": true, "团": true,
	"境": true, "复": true, "台": true, "伙": true,
	"么": true, "无": true, "惑": true,
}

func isBadNameChar(char string) bool {
	return badChars[char]
}

func TestNameAnalysis(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Database = config.DBConfig{
		Driver: "sqlite3",
		Name:   testDBPath,
	}

	f, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	born, _ := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")

	s := f.NewSessionWithFilter(NewFilter(FilterOption{
		CharacterFilter:     true,
		CharacterFilterType: CharacterFilterTypeDefault,
		MinStroke:           3,
		MaxStroke:           18,
		RegularFilter:       true,
		DaYanFilter:         true,
		WuXingFilter:        true,
	}))

	input := &Input{
		Last: [2]string{"张", ""},
		Born: born,
		Sex:  Sex(1),
	}

	start := time.Now()
	err = s.Start(input)
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	s.Wait()
	elapsed := time.Since(start)

	output := input.Output()
	total := output.Total()
	t.Logf("共生成 %d 个名字, 耗时 %v", total, elapsed)

	fateData, _ := v2.GetFateData(&v2.FateInput{
		BirthDate: born,
		Gender:    1,
		Surname:   "张",
	})

	lastName := output.Basic().LastName
	l1 := lastName[0].ScienceStroke
	var l2 int
	if lastName[1] != nil {
		l2 = lastName[1].ScienceStroke
	}

	rater := rating.NewRater(fateData)

	var analyses []nameAnalysis
	for name, ok := output.NextName(); ok; name, ok = output.NextName() {
		c1 := name.FirstName[0]
		c2 := name.FirstName[1]

		if isBadNameChar(c1.Char) || isBadNameChar(c2.Char) {
			continue
		}

		ge := wuge.CalcWuGe(l1, l2, c1.ScienceStroke, c2.ScienceStroke)
		tianGe := ge.TianGe()
		renGe := ge.RenGe()
		diGe := ge.DiGe()
		waiGe := ge.WaiGe()
		zongGe := ge.ZongGe()

		sanCai := wuxing.NewSanCai(tianGe, renGe, diGe)
		sanCaiStr := sanCai.String()
		sanCaiLuck, _ := wuxing.GetWuXing(sanCaiStr)

		nr := rater.RateName("张", c1, c2)

		analyses = append(analyses, nameAnalysis{
			name:       name,
			rating:     nr,
			tianGe:     tianGe,
			renGe:      renGe,
			diGe:       diGe,
			waiGe:      waiGe,
			zongGe:    zongGe,
			sanCaiStr:  sanCaiStr,
			sanCaiLuck: sanCaiLuck,
		})
	}

	sort.Slice(analyses, func(i, j int) bool {
		return analyses[i].rating.TotalScore > analyses[j].rating.TotalScore
	})

	fmt.Println()
	fmt.Println(strings.Repeat("═", 64))
	fmt.Println("                    姓名分析报告（美名腾风格）")
	fmt.Println(strings.Repeat("═", 64))

	fmt.Println()
	fmt.Println("【八字信息】")
	fmt.Printf("  出生时间: %s\n", born.Format("2006年01月02日 15:04"))
	if fateData != nil && fateData.Bazi != nil {
		bz := fateData.Bazi
		fmt.Printf("  四柱: %s  %s  %s  %s\n", bz.Sizhu[0], bz.Sizhu[1], bz.Sizhu[2], bz.Sizhu[3])
		fmt.Printf("  五行: %s  %s  %s  %s\n", bz.Wuxing[0], bz.Wuxing[1], bz.Wuxing[2], bz.Wuxing[3])
		fmt.Printf("  纳音: %s  %s  %s  %s\n", bz.Nayin[0], bz.Nayin[1], bz.Nayin[2], bz.Nayin[3])
		fmt.Printf("  生肖: %s    星座: %s\n", bz.Zodiac, bz.Constellation)
	}

	if fateData != nil && fateData.WuxingXiji != nil {
		wx := fateData.WuxingXiji
		fmt.Println()
		fmt.Println("【五行喜忌分析】")
		fmt.Printf("  日主五行: %s    强弱: %s\n", wx.DayWuxing, wx.QiangRuo)
		fmt.Printf("  喜用神: %s\n", strings.Join(wx.XiWuxing, "、"))
		fmt.Printf("  忌  神: %s\n", strings.Join(wx.JiWuxing, "、"))
		fmt.Printf("  分析: %s\n", wx.Analysis)
	}

	topN := 20
	filtered := len(analyses)
	if filtered < topN {
		topN = filtered
	}

	fmt.Println()
	fmt.Printf("【推荐名字 TOP %d（按评分排序）】\n", topN)
	fmt.Println(strings.Repeat("─", 64))

	for i := 0; i < topN; i++ {
		a := analyses[i]
		c1 := a.name.FirstName[0]
		c2 := a.name.FirstName[1]

		tianDaYan := wuge.Find(a.tianGe)
		renDaYan := wuge.Find(a.renGe)
		diDaYan := wuge.Find(a.diGe)
		waiDaYan := wuge.Find(a.waiGe)
		zongDaYan := wuge.Find(a.zongGe)

		fmt.Printf("\n  %2d. %s\n", i+1, a.name.String())
		fmt.Printf("      笔画: %s\n", a.name.Strokes())
		fmt.Printf("      ┌─────────────────────────────────────────────┐\n")
		fmt.Printf("      │ 天格 %2d  %s  %s\n", a.tianGe, tianDaYan.Lucky, tianDaYan.SkyNine)
		fmt.Printf("      │ 人格 %2d  %s  %s\n", a.renGe, renDaYan.Lucky, renDaYan.SkyNine)
		fmt.Printf("      │ 地格 %2d  %s  %s\n", a.diGe, diDaYan.Lucky, diDaYan.SkyNine)
		fmt.Printf("      │ 外格 %2d  %s  %s\n", a.waiGe, waiDaYan.Lucky, waiDaYan.SkyNine)
		fmt.Printf("      │ 总格 %2d  %s  %s\n", a.zongGe, zongDaYan.Lucky, zongDaYan.SkyNine)
		fmt.Printf("      └─────────────────────────────────────────────┘\n")
		fmt.Printf("      三才: %s（%s）\n", a.sanCaiStr, a.sanCaiLuck)
		fmt.Printf("      五行: %s%s", c1.WuXing, c2.WuXing)
		if fateData != nil && fateData.WuxingXiji != nil {
			xiList := fateData.WuxingXiji.XiWuxing
			if containsStr(xiList, c1.WuXing) {
				fmt.Printf("  「%s」为喜用", c1.Char)
			}
			if containsStr(xiList, c2.WuXing) {
				fmt.Printf("  「%s」为喜用", c2.Char)
			}
		}
		fmt.Println()

		p1 := ""
		p2 := ""
		if len(c1.Pinyin) > 0 {
			p1 = c1.Pinyin[0]
		}
		if len(c2.Pinyin) > 0 {
			p2 = c2.Pinyin[0]
		}
		fmt.Printf("      读音: %s %s\n", p1, p2)
		if c1.Meaning != "" || c2.Meaning != "" {
			m1 := c1.Meaning
			m2 := c2.Meaning
			if len(m1) > 30 {
				m1 = m1[:30] + "…"
			}
			if len(m2) > 30 {
				m2 = m2[:30] + "…"
			}
			fmt.Printf("      字义: %s—%s；%s—%s\n", c1.Char, m1, c2.Char, m2)
		}

		fmt.Printf("      评分: %.1f分（%s）\n", a.rating.TotalScore, a.rating.Grade)
		fmt.Printf("      五行分: %.0f  笔画分: %.0f  音韵分: %.0f\n",
			a.rating.WuXingScore, a.rating.BiHuaScore, a.rating.YinYunScore)
		fmt.Printf("      解读: %s\n", a.rating.Interpret)
	}

	fmt.Println()
	fmt.Println(strings.Repeat("═", 64))
	fmt.Printf("  共 %d 个吉名可选（过滤后 %d 个），以上展示评分最高的 %d 个\n", total, filtered, topN)
	fmt.Printf("  起名耗时: %v\n", elapsed)
	fmt.Println(strings.Repeat("═", 64))
}

func containsStr(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
