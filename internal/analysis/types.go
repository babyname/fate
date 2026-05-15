package analysis

import (
	"github.com/babyname/fate/ent"
)

// GeItem represents a single ge (格) item in WuGe analysis
type GeItem struct {
	Name          string `json:"name"`
	Stroke        int    `json:"stroke"`
	Lucky         string `json:"lucky"`
	DaYan         string `json:"da_yan"`
	SkyNine       string `json:"sky_nine"`
	YinYangWuXing string `json:"yin_yang_wu_xing"`
	Analysis      string `json:"analysis"`
}

// WuGeResult represents the result of WuGe analysis
type WuGeResult struct {
	TianGe GeItem `json:"tian_ge"`
	RenGe  GeItem `json:"ren_ge"`
	DiGe   GeItem `json:"di_ge"`
	WaiGe  GeItem `json:"wai_ge"`
	ZongGe GeItem `json:"zong_ge"`
}

// CharInfo represents character information
type CharInfo struct {
	Char              string `json:"char"`
	TraditionalChar   string `json:"traditional_char"`
	Pinyin            string `json:"pinyin"`
	WuXing            string `json:"wu_xing"`
	SimplifiedStroke  int    `json:"simplified_stroke"`
	TraditionalStroke int    `json:"traditional_stroke"`
	ScienceStroke     int    `json:"science_stroke"`
	KangxiStroke      int    `json:"kangxi_stroke"`
	Radical           string `json:"radical"`
	Meaning           string `json:"meaning"`
	IsXiYong          bool   `json:"is_xi_yong"`
}

// ScoreDetail represents detailed score breakdown
type ScoreDetail struct {
	WenHuaYinXiang float64 `json:"wen_hua_yin_xiang"`
	WuXingBaZi     float64 `json:"wu_xing_ba_zi"`
	ShengXiao      float64 `json:"sheng_xiao"`
	WuGeShuLi      float64 `json:"wu_ge_shu_li"`
}

// NameResult represents a single name analysis result
type NameResult struct {
	Rank         int           `json:"rank"`
	FullName     string        `json:"full_name"`
	Surname      string        `json:"surname"`
	FirstName    string        `json:"first_name"`
	Strokes      string        `json:"strokes"`
	Char1        CharInfo      `json:"char1"`
	Char2        CharInfo      `json:"char2"`
	WuGe         *WuGeResult   `json:"wu_ge"`
	SanCai       string        `json:"san_cai"`
	SanCaiLuck   string        `json:"san_cai_luck"`
	SanCaiDetail string        `json:"san_cai_detail"`
	JiChuYun     string        `json:"ji_chu_yun"`
	ChengGongYun string        `json:"cheng_gong_yun"`
	RenJiGuanXi  string        `json:"ren_ji_guan_xi"`
	ZhouYi       *ZhouYiResult `json:"zhou_yi"`
	Score        float64       `json:"score"`
	Grade        string        `json:"grade"`
	ScoreDetail  ScoreDetail   `json:"score_detail"`
	Interpret    string        `json:"interpret"`
}

// BaziSection represents bazi information section
type BaziSection struct {
	Sizhu         [4]string `json:"sizhu"`
	Wuxing        [4]string `json:"wuxing"`
	Nayin         [4]string `json:"nayin"`
	Zodiac        string    `json:"zodiac"`
	Constellation string    `json:"constellation"`
}

// WuXingSection represents wuxing analysis section
type WuXingSection struct {
	DayGan     string   `json:"day_gan"`
	DayWuxing  string   `json:"day_wuxing"`
	QiangRuo   string   `json:"qiang_ruo"`
	XiWuxing   []string `json:"xi_wuxing"`
	YongWuxing string   `json:"yong_wuxing"`
	JiWuxing   []string `json:"ji_wuxing"`
	ChouWuxing []string `json:"chou_wuxing"`
	XianWuxing []string `json:"xian_wuxing"`
	Method     string   `json:"method"`
	MethodName string   `json:"method_name"`
	GeJuName   string   `json:"geju_name"`
	Analysis   string   `json:"analysis"`
}

// FateReport represents a complete fate analysis report
type FateReport struct {
	GeneratedAt string         `json:"generated_at"`
	Surname     string         `json:"surname"`
	Born        string         `json:"born"`
	Sex         string         `json:"sex"`
	Bazi        *BaziSection   `json:"bazi"`
	WuXing      *WuXingSection `json:"wu_xing"`
	TotalNames  int            `json:"total_names"`
	TopNames    []NameResult   `json:"top_names"`
}

// NameSource represents source characters for a name
type NameSource struct {
	C1 *ent.Character
	C2 *ent.Character
}
