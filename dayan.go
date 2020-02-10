package fate

import (
	_ "github.com/godcong/fate/statik"
)

var daYanList [81]DaYan

func init() {
	daYanList = [81]DaYan{
		{Number: 1, Lucky: "吉", SkyNine: "太极之数", Comment: "太极之数，万物开泰，生发无穷，利禄亨通。"},
		{Number: 2, Lucky: "凶", SkyNine: "两仪之数", Comment: "两仪之数，混沌未开，进退保守，志望难达。"},
		{Number: 3, Lucky: "吉", SkyNine: "三才之数", Comment: "三才之数，天地人和，大事大业，繁荣昌隆。"},
		{Number: 4, Lucky: "凶", SkyNine: "四象之数", Comment: "四象之数，待于生发，万事慎重，不具营谋。"},
		{Number: 5, Lucky: "吉", SkyNine: "五行之数", Comment: "五行俱权，循环相生，圆通畅达，福祉无穷。"},
		{Number: 6, Lucky: "吉", SkyNine: "六爻之数", Comment: "六爻之数，发展变化，天赋美德，吉祥安泰。"},
		{Number: 7, Lucky: "吉", SkyNine: "七政之数", Comment: "七政之数，精悍严谨，天赋之力，吉星照耀。"},
		{Number: 8, Lucky: "半吉", SkyNine: "八卦之数", Comment: "八卦之数，乾坎艮震，巽离坤兑，无穷无尽。"},
		{Number: 9, Lucky: "凶", SkyNine: "大成之数", Comment: "大成之数，蕴涵凶险，或成或败，难以把握。"},
		{Number: 10, Lucky: "凶", SkyNine: "终结之数", Comment: "终结之数，雪暗飘零，偶或有成，回顾茫然。"},
		{Number: 11, Lucky: "吉", SkyNine: "旱苗逢雨", Comment: "万物更新，调顺发达，恢弘泽世，繁荣富贵。"},
		{Number: 12, Lucky: "凶", SkyNine: "掘井无泉", Comment: "无理之数，发展薄弱，虽生不足，难酬志向。"},
		{Number: 13, Lucky: "吉", SkyNine: "春日牡丹", Comment: "才艺多能，智谋奇略，忍柔当事，鸣奏大功。"},
		{Number: 14, Lucky: "凶", SkyNine: "破兆", Comment: "家庭缘薄，孤独遭难，谋事不达，悲惨不测。"},
		{Number: 15, Lucky: "吉", SkyNine: "福寿", Comment: "福寿圆满，富贵荣誉，涵养雅量，德高望重。"},
		{Number: 16, Lucky: "吉", SkyNine: "厚重", Comment: "厚重载德，安富尊荣，财官双美，功成名就。"},
		{Number: 17, Lucky: "半吉", SkyNine: "刚强", Comment: "权威刚强，突破万难，如能容忍，必获成功。"},
		{Number: 18, Lucky: "半吉", SkyNine: "铁镜重磨", Comment: "权威显达，博得名利，且养柔德，功成名就。"},
		{Number: 19, Lucky: "凶", SkyNine: "多难", Comment: "风云蔽日，辛苦重来，虽有智谋，万事挫折。"},
		{Number: 20, Lucky: "凶", SkyNine: "屋下藏金", Comment: "非业破运，灾难重重，进退维谷，万事难成。"},
		{Number: 21, Lucky: "吉", Sex: true, SkyNine: "明月中天", Comment: "光风霁月，万物确立，官运亨通，大搏名利。女性不宜此数。"},
		{Number: 22, Lucky: "凶", SkyNine: "秋草逢霜", Comment: "秋草逢霜，困难疾弱，虽出豪杰，人生波折。"},
		{Number: 23, Lucky: "吉", Sex: true, SkyNine: "壮丽", Comment: "旭日东升，壮丽壮观，权威旺盛，功名荣达。女性不宜此数。"},
		{Number: 24, Lucky: "吉", SkyNine: "掘藏得金", Comment: "家门余庆，金钱丰盈，白手成家，财源广进。"},
		{Number: 25, Lucky: "半吉", SkyNine: "荣俊", Comment: "资性英敏，才能奇特，克服傲慢，尚可成功。"},
		{Number: 26, Lucky: "凶", SkyNine: "变怪", Comment: "变怪之谜，英雄豪杰，波澜重叠，而奏大功。"},
		{Number: 27, Lucky: "凶", SkyNine: "增长", Comment: "欲望无止，自我强烈，多受毁谤，尚可成功。"},
		{Number: 28, Lucky: "凶", Sex: true, SkyNine: "阔水浮萍", Comment: "遭难之数，豪杰气概，四海漂泊，终世浮躁。女性不宜此数。"},
		{Number: 29, Lucky: "吉", SkyNine: "智谋", Comment: "智谋优秀，财力归集，名闻海内，成就大业。"},
		{Number: 30, Lucky: "半吉", SkyNine: "非运", Comment: "沉浮不定，凶吉难变，若明若暗，大成大败。"},
		{Number: 31, Lucky: "吉", SkyNine: "春日花开", Comment: "智勇得志，博得名利，统领众人，繁荣富贵。"},
		{Number: 32, Lucky: "吉", SkyNine: "宝马金鞍", Comment: "侥幸多望，贵人得助，财帛如裕，繁荣至上。"},
		{Number: 33, Lucky: "吉", Sex: true, SkyNine: "旭日升天", Comment: "旭日升天，鸾凤相会，名闻天下，隆昌至极。女性不宜此数。"},
		{Number: 34, Lucky: "凶", SkyNine: "破家", Comment: "破家之身，见识短小，辛苦遭逢，灾祸至极。"},
		{Number: 35, Lucky: "吉", SkyNine: "高楼望月", Comment: "温和平静，智达通畅，文昌技艺，奏功洋洋。"},
		{Number: 36, Lucky: "半吉", SkyNine: "波澜重叠", Comment: "波澜重叠，沉浮万状，侠肝义胆，舍己成仁。"},
		{Number: 37, Lucky: "吉", SkyNine: "猛虎出林", Comment: "权威显达，热诚忠信，宜着雅量，终身荣富。"},
		{Number: 38, Lucky: "半吉", SkyNine: "磨铁成针", Comment: "意志薄弱，刻意经营，才识不凡，技艺有成。"},
		{Number: 39, Lucky: "半吉", SkyNine: "富贵荣华", Comment: "富贵荣华，财帛丰盈，暗藏险象，德泽四方。"},
		{Number: 40, Lucky: "凶", SkyNine: "退安", Comment: "智谋胆力，冒险投机，沉浮不定，退保平安。"},
		{Number: 41, Lucky: "吉", Max: true, SkyNine: "有德", Comment: "纯阳独秀，德高望重，和顺畅达，博得名利。此数为最大好运数。"},
		{Number: 42, Lucky: "凶", SkyNine: "寒蝉在柳", Comment: "博识多能，精通世情，如能专心，尚可成功。"},
		{Number: 43, Lucky: "凶", SkyNine: "散财破产", Comment: "散财破产，诸事不遂，虽有智谋，财来财去。"},
		{Number: 44, Lucky: "凶", SkyNine: "烦闷", Comment: "破家亡身，暗藏惨淡，事不如意，乱世怪杰。"},
		{Number: 45, Lucky: "吉", SkyNine: "顺风", Comment: "新生泰和，顺风扬帆，智谋经纬，富贵繁荣。"},
		{Number: 46, Lucky: "凶", SkyNine: "浪里淘金", Comment: "载宝沉舟，浪里淘金，大难尝尽，大功有成。"},
		{Number: 47, Lucky: "吉", SkyNine: "点石成金", Comment: "花开之象，万事如意，祯祥吉庆，天赋幸福。"},
		{Number: 48, Lucky: "吉", SkyNine: "古松立鹤", Comment: "智谋兼备，德量荣达，威望成师，洋洋大观。"},
		{Number: 49, Lucky: "半吉", SkyNine: "转变", Comment: "吉临则吉，凶来则凶，转凶为吉，配好三才。"},
		{Number: 50, Lucky: "半吉", SkyNine: "小舟入海", Comment: "一成一败，吉凶参半，先得庇荫，后遭凄惨。"},
		{Number: 51, Lucky: "半吉", SkyNine: "沉浮", Comment: "盛衰交加，波澜重叠，如能慎始，必获成功。"},
		{Number: 52, Lucky: "吉", SkyNine: "达眼", Comment: "卓识达眼，先见之明，智谋超群，名利双收。"},
		{Number: 53, Lucky: "凶", SkyNine: "曲卷难星", Comment: "外祥内患，外祸内安，先富后贫，先贫后富。"},
		{Number: 54, Lucky: "凶", SkyNine: "石上栽花", Comment: "石上栽花，难得有活，忧闷烦来，辛惨不绝。"},
		{Number: 55, Lucky: "半吉", SkyNine: "善恶", Comment: "善善得恶，恶恶得善，吉到极限，反生凶险。"},
		{Number: 56, Lucky: "凶", SkyNine: "浪里行舟", Comment: "历尽艰辛，四周障碍，万事龃龌，做事难成。"},
		{Number: 57, Lucky: "吉", SkyNine: "日照春松", Comment: "寒雪青松，夜莺吟春，必遭一过，繁荣白事。"},
		{Number: 58, Lucky: "半吉", SkyNine: "晚行遇月", Comment: "沉浮多端，先苦后甜，宽宏扬名，富贵繁荣。"},
		{Number: 59, Lucky: "凶", SkyNine: "寒蝉悲风", Comment: "寒蝉悲风，意志衰退，缺乏忍耐，苦难不休。"},
		{Number: 60, Lucky: "凶", SkyNine: "无谋", Comment: "无谋之人，漂泊不定，晦暝暗黑，动摇不安。"},
		{Number: 61, Lucky: "吉", SkyNine: "牡丹芙蓉", Comment: "牡丹芙蓉，花开富贵，名利双收，定享天赋。"},
		{Number: 62, Lucky: "凶", SkyNine: "衰败", Comment: "衰败之象，内外不和，志望难达，灾祸频来。"},
		{Number: 63, Lucky: "吉", SkyNine: "舟归平海", Comment: "富贵荣华，身心安泰，雨露惠泽，万事亨通。"},
		{Number: 64, Lucky: "凶", SkyNine: "非命", Comment: "骨肉分离，孤独悲愁，难得心安，做事不成。"},
		{Number: 65, Lucky: "吉", SkyNine: "巨流归海", Comment: "天长地久，家运隆昌，福寿绵长，事事成就。"},
		{Number: 66, Lucky: "凶", SkyNine: "岩头步马", Comment: "进退维谷，艰难不堪，等待时机，一跃而起。"},
		{Number: 67, Lucky: "吉", SkyNine: "顺风通达", Comment: "天赋幸运，四通八达，家道繁昌，富贵东来。"},
		{Number: 68, Lucky: "吉", SkyNine: "顺风吹帆", Comment: "智虑周密，集众信达，发明能智，拓展昂进。"},
		{Number: 69, Lucky: "凶", SkyNine: "非业", Comment: "非业非力，精神迫滞，灾害交至，遍偿痛苦。"},
		{Number: 70, Lucky: "凶", SkyNine: "残菊逢霜", Comment: "残菊逢霜，寂寞无碍，惨淡忧愁，晚景凄凉。"},
		{Number: 71, Lucky: "半吉", SkyNine: "石上金花", Comment: "石上金花，内心劳苦，贯彻始终，定可昌隆。"},
		{Number: 72, Lucky: "半吉", SkyNine: "劳苦", Comment: "荣苦相伴，阴云覆月，外表吉祥，内实凶祸。"},
		{Number: 73, Lucky: "半吉", SkyNine: "无勇", Comment: "盛衰交加，徒有高志，天王福祉，终世平安。"},
		{Number: 74, Lucky: "凶", SkyNine: "残菊经霜", Comment: "残菊经霜，秋叶寂寞，无能无智，辛苦繁多。"},
		{Number: 75, Lucky: "凶", SkyNine: "退守", Comment: "退守保吉，发迹甚迟，虽有吉象，无谋难成。"},
		{Number: 76, Lucky: "凶", SkyNine: "离散", Comment: "倾覆离散，骨肉分离，内外不和，虽劳无功。"},
		{Number: 77, Lucky: "半吉", SkyNine: "半吉", Comment: "家庭有悦，半吉半凶，能获援护，陷落不幸。"},
		{Number: 78, Lucky: "凶", SkyNine: "晚苦", Comment: "祸福参半，先天智能，中年发达，晚景困苦。"},
		{Number: 79, Lucky: "凶", SkyNine: "云头望月", Comment: "云头望月，身疲力尽，穷迫不伸，精神不定。"},
		{Number: 80, Lucky: "凶", SkyNine: "遁吉", Comment: "辛苦不绝，早入隐遁，安心立命，化凶转吉。"},
		{Number: 81, Lucky: "吉", SkyNine: "万物回春", Comment: "最吉之数，还本归元，吉祥重叠，富贵尊荣。"},
	}
}

// DaYan ...
type DaYan struct {
	Number  int
	Lucky   string
	Max     bool
	Sex     bool //male(false),female(true)
	SkyNine string
	Comment string
}

//IsSex 女性不宜此数
func (dy DaYan) IsSex() bool {
	return dy.Sex
}

//IsMax 是否最大好运数
func (dy DaYan) IsMax() bool {
	return dy.Max
}

//GetDaYan 获取大衍之数
func GetDaYan(idx int) DaYan {
	if idx <= 0 {
		panic("wrong idx")
	}
	i := (idx - 1) % 81
	return daYanList[i]
}
