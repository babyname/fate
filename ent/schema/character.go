package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Character holds the schema definition for the Character entity.
type Character struct {
	ent.Schema
}

const WuXingGold = "金"
const WuXingWood = "木"
const WuXingWater = "水"
const WuXingFire = "火"
const WuXingSoil = "土"
const WuXingNone = "无"

//const "大凶"
//const "凶"
//const "凶多于吉"
//const "吉凶参半"
//const "吉多于凶"
//const "吉"
//const "大吉"

// Fields of the Character.
func (Character) Fields() []ent.Field {
	return []ent.Field{
		field.String("hash").Unique(),
		//Hash                     string   `xorm:"pk hash"`
		field.Strings("pin_yin"),
		//PinYin                   []string `xorm:"default() notnull pin_yin"`                               //拼音
		field.String("ch"),
		//Ch                       string   `xorm:"default() notnull ch"`                                    //字符
		field.Int8("science_stroke"),
		//ScienceStroke            int      `xorm:"default(0) notnull science_stroke" json:"science_stroke"` //科学笔画
		field.String("radical"),
		//Radical                  string   `xorm:"default() notnull radical"`                               //部首
		field.Int8("radical_stroke"),
		//RadicalStroke            int      `xorm:"default(0) notnull radical_stroke"`                       //部首笔画
		field.Int8("stroke"),
		//Stroke                   int      `xorm:"default() notnull stroke"`                                //总笔画数
		field.Bool("is_kang_xi"),
		//IsKangXi                 bool     `xorm:"default(0) notnull is_kang_xi"`                           //是否康熙字典
		field.String("kang_xi"),
		//KangXi                   string   `xorm:"default() notnull kang_xi"`                               //康熙
		field.String("kang_xi_stroke"),
		//KangXiStroke             int      `xorm:"default(0) notnull kang_xi_stroke"`                       //康熙笔画
		field.String("simple_radical"),
		//SimpleRadical            string   `xorm:"default() notnull simple_radical"`                        //简体部首
		field.String("simple_radical_stroke"),
		//SimpleRadicalStroke      int      `xorm:"default(0) notnull simple_radical_stroke"`                //简体部首笔画
		field.Int8("simple_total_stroke"),
		//SimpleTotalStroke        int      `xorm:"default(0) notnull simple_total_stroke"`                  //简体笔画
		field.String("traditional_radical"),
		//TraditionalRadical       string   `xorm:"default() notnull traditional_radical"`                   //繁体部首
		field.Int8("traditional_radical_stroke"),
		//TraditionalRadicalStroke int      `xorm:"default(0) notnull traditional_radical_stroke"`           //繁体部首笔画
		field.Int8("traditional_total_stroke"),
		//TraditionalTotalStroke   int      `xorm:"default(0) notnull traditional_total_stroke"`             //简体部首笔画
		field.Bool("is_name_science"),
		//NameScience              bool     `xorm:"default(0) notnull name_science"`                         //姓名学
		field.String("wu_xing"),
		//WuXing                   string   `xorm:"default() notnull wu_xing"`                               //五行
		field.String("lucky"),
		//Lucky                    string   `xorm:"default() notnull lucky"`                                 //吉凶寓意
		field.Bool("is_regular"),
		//Regular                  bool     `xorm:"default(0) notnull regular"`                              //常用
		field.Strings("traditional_character"),
		//TraditionalCharacter     []string `xorm:"default() notnull traditional_character"`                 //繁体字
		field.Strings("variant_character"),
		//VariantCharacter         []string `xorm:"default() notnull variant_character"`                     //异体字
		field.Text("comment"),
		//Comment                  []string `xorm:"default() notnull comment"`                               //解释
	}
}

// Edges of the Character.
func (Character) Edges() []ent.Edge {
	return nil
}
