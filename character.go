package fate

import (
	xt "github.com/free-utils-go/xorm_type_assist"
	"xorm.io/builder"
	"xorm.io/xorm"
)

//Character 字符
type Character struct {
	Ch            string      `xorm:"not null pk comment('汉字') VARCHAR(1)"`
	PinYin        []string    `xorm:"not null comment('拼音') TEXT"`
	IsRegular     xt.BoolType `xorm:"not null comment('常用,~bool') VARCHAR(1)"`
	IsDuoYin      xt.BoolType `xorm:"not null comment('多音字,~bool') VARCHAR(1)"`
	IsSurname     xt.BoolType `xorm:"comment('姓氏,~bool') VARCHAR(1)"` //姓名学（是否为姓氏）
	SurnameGender string      `xorm:"comment('性别') VARCHAR(1)"`
	WuXing        string      `xorm:"not null comment('五行') VARCHAR(1)"`
	Lucky         string      `xorm:"comment('吉凶寓意') TEXT"`
	Radical       string      `xorm:"not null comment('部首') VARCHAR(1)"`
	Stroke        int         `xorm:"default 0 comment('笔画数') INT(11)"`    //手写笔画数
	ScienceStroke int         `xorm:"default 0 comment('姓名学笔画数') INT(11)"` //科学笔画
}

// CharacterOptions ...
type CharacterOptions func(session *xorm.Session) *xorm.Session

// Regular 过滤常用字符
func Regular() CharacterOptions {
	return func(session *xorm.Session) *xorm.Session {
		return session.And("is_regular = ?", xt.TRUE)
	}
}

//get common stroke
//名需要采用这种笔画
func (ch *Character) getStroke() int {
	var result int
	if ch.Stroke != 0 {
		result = ch.Stroke
	}

	if result == 0 {
		panic("笔画数找不到")
	}

	return result
}

//get ancient stroke
func (ch *Character) getStrokeScience() int {
	var result int
	if ch.ScienceStroke != 0 {
		result = ch.ScienceStroke
	}

	if result == 0 {
		panic("ScienceStroke找不到")
	}

	return result
}

// Stoker 设置笔画数条件
func Stoker(s int, options ...CharacterOptions) func(engine *xorm.Engine) *xorm.Session {
	return func(engine *xorm.Engine) *xorm.Session {
		session := engine.Where("pin_yin IS NOT NULL").
			And(builder.Eq{"stroke": s})
		for _, option := range options {
			session = option(session)
		}
		return session
	}

}

// Stoker 设置康熙笔画数条件
func StokerX(s int, options ...CharacterOptions) func(engine *xorm.Engine) *xorm.Session {
	return func(engine *xorm.Engine) *xorm.Session {
		session := engine.Where("pin_yin IS NOT NULL").
			And(builder.Eq{"science_stroke": s})
		for _, option := range options {
			session = option(session)
		}
		return session
	}

}

// 从字符产生查询条件
func Char(name string) func(engine *xorm.Engine) *xorm.Session {
	return func(engine *xorm.Engine) *xorm.Session {
		return engine.Where(builder.Eq{"ch": name})
	}
}
