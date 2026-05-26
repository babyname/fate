package analysis

import (
	"github.com/babyname/fate/ent"
)

type GeItem struct {
	Name          string `json:"name"`
	Stroke        int    `json:"stroke"`
	Lucky         string `json:"lucky"`
	DaYan         string `json:"da_yan"`
	SkyNine       string `json:"sky_nine"`
	YinYangWuXing string `json:"yin_yang_wu_xing"`
	Analysis      string `json:"analysis"`
}

type WuGeResult struct {
	TianGe GeItem `json:"tian_ge"`
	RenGe  GeItem `json:"ren_ge"`
	DiGe   GeItem `json:"di_ge"`
	WaiGe  GeItem `json:"wai_ge"`
	ZongGe GeItem `json:"zong_ge"`
}

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

type ScoreDetail struct {
	WenHuaYinXiang float64 `json:"wen_hua_yin_xiang"`
	WuXingBaZi     float64 `json:"wu_xing_ba_zi"`
	ShengXiao      float64 `json:"sheng_xiao"`
	WuGeShuLi      float64 `json:"wu_ge_shu_li"`
	YinYun         float64 `json:"yin_yun"`
}

type NameResult struct {
	Rank         int            `json:"rank"`
	FullName     string         `json:"full_name"`
	Surname      string         `json:"surname"`
	FirstName    string         `json:"first_name"`
	Strokes      string         `json:"strokes"`
	Char1        CharInfo       `json:"char1"`
	Char2        CharInfo       `json:"char2"`
	WuGe         *WuGeResult    `json:"wu_ge"`
	SanCai       string         `json:"san_cai"`
	SanCaiLuck   string         `json:"san_cai_luck"`
	SanCaiDetail string         `json:"san_cai_detail"`
	JiChuYun     string         `json:"ji_chu_yun"`
	ChengGongYun string         `json:"cheng_gong_yun"`
	RenJiGuanXi  string         `json:"ren_ji_guan_xi"`
	ZhouYi       *ZhouYiBasic   `json:"zhou_yi"`
	Bazi         *BaziBasic     `json:"bazi,omitempty"`
	WuXing       *WuXingSection `json:"wu_xing,omitempty"`
	Score        float64        `json:"score"`
	Grade        string         `json:"grade"`
	ScoreDetail  ScoreDetail    `json:"score_detail"`
	Interpret    string         `json:"interpret"`
	HasPoetry    bool           `json:"has_poetry"`
}

type ZhouYiBasic struct {
	BenGuaName     string `json:"ben_gua_name"`
	BenGuaDesc     string `json:"ben_gua_desc"`
	BenGuaJiXiong  string `json:"ben_gua_ji_xiong"`
	BianGuaName    string `json:"bian_gua_name"`
	DongYaoDesc    string `json:"dong_yao_desc"`
	DongYaoJiXiong string `json:"dong_yao_ji_xiong"`
	DaXiang        string `json:"da_xiang"`
	Score          int    `json:"score"`
}

type BaziBasic struct {
	FourPillars   [4]string `json:"four_pillars"`
	FiveElements  [4]string `json:"five_elements"`
	NaYin         [4]string `json:"na_yin"`
	Zodiac        string    `json:"zodiac"`
	Constellation string    `json:"constellation"`
}

type WuXingSection struct {
	DayGan            string   `json:"day_gan"`
	DayWuxing         string   `json:"day_wuxing"`
	Strength          string   `json:"strength"`
	FavorableElements []string `json:"favorable_elements"`
	UsefulElement     string   `json:"useful_element"`
	UnfavorableElements []string `json:"unfavorable_elements"`
	HostileElements   []string `json:"hostile_elements"`
	IdleElements      []string `json:"idle_elements"`
	Method            string   `json:"method"`
	MethodName        string   `json:"method_name"`
	GeJuName          string   `json:"geju_name"`
	Analysis          string   `json:"analysis"`
}

type NameSource struct {
	C1 *ent.Character
	C2 *ent.Character
}
