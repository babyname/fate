package chronosfate

import (
	"github.com/godcong/chronos/v2"
)

type XiYongMethod int

const (
	XiYongMethodBalance XiYongMethod = iota
	XiYongMethodGeJu
)

type FateInput struct {
	Calendar     chronos.Calendar
	Gender       int
	XiYongMethod XiYongMethod
}

type FateData struct {
	BaziInfo       BaziInfo       `json:"bazi_info"`
	WuxingXijiInfo WuxingXijiInfo `json:"wuxing_xiji_info"`
	WuxingStrength WuxingStrength `json:"wuxing_strength"`
	XiYongJiChou   XiYongJiChou   `json:"xi_yong_ji_chou"`
	GeJuInfo       *GeJuInfo      `json:"ge_ju_info,omitempty"`
}

type BaziInfo struct {
	SiZhu         [4]string   `json:"si_zhu"`
	WuXing        [4]string   `json:"wu_xing"`
	NaYin         [4]string   `json:"na_yin"`
	ShiShenGan    [4]string   `json:"shi_shen_gan"`
	ShiShenZhi    [4][]string `json:"shi_shen_zhi"`
	CangGan       [4][]string `json:"cang_gan"`
	DaYun         []int       `json:"da_yun"`
	Zodiac        string      `json:"zodiac"`
	Constellation string      `json:"constellation"`
}

type WuxingXijiInfo struct {
	Xi            string `json:"xi"`
	Ji            string `json:"ji"`
	RiZhuQiangRuo string `json:"ri_zhu_qiang_ruo"`
	TiaoHouShen   string `json:"tiao_hou_shen"`
}

type WuxingStrength struct {
	WuxingFen map[string]float64 `json:"wuxing_fen"`
	Total     float64            `json:"total"`
}

type XiYongJiChou struct {
	Xi   string `json:"xi"`
	Yong string `json:"yong"`
	Ji   string `json:"ji"`
	Chou string `json:"chou"`
}

type GeJuType int

const (
	GeJuZhengGuan GeJuType = iota
	GeJuQiSha
	GeJuZhengCai
	GeJuPianCai
	GeJuZhengYin
	GeJuPianYin
	GeJuShiShen
	GeJuShangGuan
	GeJuUnknown
)

type GeJuInfo struct {
	Type     GeJuType `json:"type"`
	Name     string   `json:"name"`
	YongShen string   `json:"yong_shen"`
	XiShen   string   `json:"xi_shen"`
	JiShen   string   `json:"ji_shen"`
	ChouShen string   `json:"chou_shen"`
	Analysis string   `json:"analysis"`
}

type FateError struct {
	Code    int
	Message string
}

func (e FateError) Error() string {
	return e.Message
}

const (
	ErrCodeInputInvalid = iota + 1
	ErrCodeDateRange
	ErrCodeCalculateBazi
	ErrCodeCalculateWuxing
)
