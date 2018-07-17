package mgo

//Character 字符
type Character struct {
	Character            string `json:"character"`              //字符
	Pinyin               string `json:"pinyin"`                 //拼音
	Radical              string `json:"radical"`                //部首
	RadicalStrokes       string `json:"radical_strokes"`        //部首笔画
	TotalStrokes         string `json:"total_strokes"`          //总笔画
	KangxiCharacter      string `json:"kangxi_character"`       //康熙字符
	KangxiStrokes        string `json:"kangxi_strokes"`         //康熙笔画数
	Phonetic             string `json:"phonetic"`               //注音
	CommonlyCharacters   string `json:"commonly_characters"`    //是否为常用字
	NameScience          string `json:"name_science"`           //姓名学
	FiveElementCharacter string `json:"five_element_character"` //汉字五行
	GodBadMoral          string `json:"god_bad_moral"`          //吉凶寓意
	DecompositionSearch  string `json:"decomposition_search"`   //首尾分解查字
	StrokeNumber         string `json:"stroke_number"`          //笔顺编号
	StrokeReadWrite      string `json:"stroke_read_write"`      //笔顺读写
}
