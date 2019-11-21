package fate

import (
	"fmt"
	"github.com/xormsharp/builder"
	"github.com/xormsharp/xorm"
)

//Character 字符
type Character struct {
	Hash   string   `xorm:"pk hash"`
	PinYin []string `xorm:"default() notnull pin_yin"` //拼音
	Ch     string   `xorm:"default() notnull ch"`      //字符
	//ScienceStroke            string   `xorm:"science_stroke" json:"science_stroke"`
	Radical                  string   `xorm:"default() notnull radical"`                     //部首
	RadicalStroke            int      `xorm:"default(0) notnull radical_stroke"`             //部首笔画
	Stroke                   int      `xorm:"default() notnull stroke"`                      //总笔画数
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

func getCharacters(f *fateImpl, fn func(engine *xorm.Engine) *xorm.Session) ([]*Character, error) {
	s := fn(f.chardb)
	var c []*Character
	e := s.Find(&c)
	if e == nil {
		return c, nil
	}
	return nil, fmt.Errorf("character list get error:%w", e)
}

func getCharacter(f *fateImpl, fn func(engine *xorm.Engine) *xorm.Session) (*Character, error) {
	s := fn(f.chardb)
	var c Character
	b, e := s.Get(&c)
	if e == nil && b {
		return &c, nil
	}
	return nil, fmt.Errorf("character get error:%w", e)
}

type CharacterOptions func(session *xorm.Session) *xorm.Session

func Regular() CharacterOptions {
	return func(session *xorm.Session) *xorm.Session {
		return session.And("regular = ?", 1)
	}
}

func Stoker(s int, options ...CharacterOptions) func(engine *xorm.Engine) *xorm.Session {
	return func(engine *xorm.Engine) *xorm.Session {
		session := engine.Where("pin_yin IS NOT NULL").
			And(builder.Eq{"stroke": s}.
				Or(builder.Eq{"kang_xi_stroke": s}).
				Or(builder.Eq{"simple_total_stroke": s}).
				Or(builder.Eq{"traditional_total_stroke": s}))
		for _, option := range options {
			session = option(session)
		}
		return session
	}

}

func Char(name string) func(engine *xorm.Engine) *xorm.Session {
	return func(engine *xorm.Engine) *xorm.Session {
		return engine.Where(builder.Eq{"ch": name}.
			Or(builder.Eq{"kang_xi": name}).
			Or(builder.Eq{"traditional_character": name}))
	}
}
