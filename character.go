package fate

import (
	"fmt"
	"github.com/go-xorm/xorm"
)

//Character 字符
type Character struct {
	PinYin                   []string `xorm:"default() notnull pin_yin"`                     //拼音
	Character                string   `xorm:"default() notnull character"`                   //字符
	Radical                  string   `xorm:"default() notnull radical"`                     //部首
	RadicalStroke            int      `xorm:"default(0) notnull radical_stroke"`             //部首笔画
	IsKangXi                 bool     `xorm:"default(0) notnull is_kang_xi"`                 //是否康熙字典
	KangXi                   string   `xorm:"default() notnull kang_xi"`                     //康熙
	KangXiStroke             int      `xorm:"default(0) notnull kang_xi_stroke"`             //康熙笔画
	SimpleRadical            string   `xorm:"default() notnull simple_radical"`              //简体部首
	SimpleRadicalStroke      int      `xorm:"default(0) notnull simple_radical_stroke"`      //简体部首笔画
	SimpleTotalStroke        int      `xorm:"default(0) notnull simple_total_stroke"`        //简体笔画
	TraditionalRadical       string   `xorm:"default() notnull traditional_radical"`         //繁体部首
	TraditionalRadicalStroke int      `xorm:"default(0) notnull traditional_radical_stroke"` //繁体部首笔画
	TraditionalTotalStroke   int      `xorm:"default(0) notnull traditional_total_stroke"`   //简体部首笔画
	NameScience              bool     `xorm:"default(0) notnull name_science"`               //姓名学
	WuXing                   string   `xorm:"default() notnull wu_xing"`                     //五行
	Lucky                    string   `xorm:"default() notnull lucky"`                       //吉凶寓意
	Regular                  bool     `xorm:"default(0) notnull regular"`                    //常用
	TraditionalCharacter     []string `xorm:"default() notnull traditional_character"`       //繁体字
	VariantCharacter         []string `xorm:"default() notnull variant_character"`           //异体字
	Comment                  []string `xorm:"default() notnull comment"`                     //解释
}

func getCharacter(f *fate, fn func(engine *xorm.Engine) *xorm.Session) (*Character, error) {
	s := fn(f.chardb)
	c := new(Character)
	b, e := s.Get(c)
	if e == nil && b {
		return c, nil
	}
	return nil, fmt.Errorf("character get error: %v", e)
}

func Char(name string) func(engine *xorm.Engine) *xorm.Session {
	return func(engine *xorm.Engine) *xorm.Session {
		return engine.Where("character = ?", name).
			Or("kang_xi = ?", name).
			Or("traditional_character = ?", name)
	}
}
