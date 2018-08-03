package mongo

import (
	"github.com/globalsign/mgo/bson"
	"strconv"
)

//Character 字符
type Character struct {
	ID                   bson.ObjectId `bson:"_id,omitempty"`          //id
	Character            string        `bson:"character"`              //字符
	Pinyin               string        `bson:"pinyin"`                 //拼音
	Radical              string        `bson:"radical"`                //部首
	RadicalStrokes       string        `bson:"radical_strokes"`        //部首笔画
	TotalStrokes         string        `bson:"total_strokes"`          //总笔画
	KangxiCharacter      string        `bson:"kangxi_character"`       //康熙字符
	KangxiStrokes        string        `bson:"kangxi_strokes"`         //康熙笔画数
	Phonetic             string        `bson:"phonetic"`               //注音
	CommonlyCharacters   string        `bson:"commonly_characters"`    //是否为常用字
	NameScience          string        `bson:"name_science"`           //姓名学
	FiveElementCharacter string        `bson:"five_element_character"` //汉字五行
	GodBadMoral          string        `bson:"god_bad_moral"`          //吉凶寓意
	DecompositionSearch  string        `bson:"decomposition_search"`   //首尾分解查字
	StrokeNumber         string        `bson:"stroke_number"`          //笔顺编号
	StrokeReadWrite      string        `bson:"stroke_read_write"`      //笔顺读写
}

const (
	KangXi = iota
	Simplified
	Traditional
	Total
	Radical
)

func (c *Character) GetStrokeByType(tp int) int {
	switch tp {
	case KangXi:
		i, _ := strconv.Atoi(c.KangxiStrokes)
		return i
	case Simplified, Total:
		i, _ := strconv.Atoi(c.TotalStrokes)
		return i
	case Traditional: //do nothing
	case Radical:
		i, _ := strconv.Atoi(c.RadicalStrokes)
		return i
	}
	return 0
}
