package fate

import (
	"github.com/godcong/chronos"
	"strings"
)

// ZodiacPig ...
const (
	ZodiacPig  = "猪"
	ZodiacMice = "鼠"
)

//Zodiac 生肖
type Zodiac struct {
	Name      string
	Xi        string //喜
	XiRadical string
	Ji        string //忌
	JiRadical string
}

var zodiacList = map[string]Zodiac{
	ZodiacPig: {
		Name:      ZodiacPig,
		Xi:        "",
		XiRadical: "",
		Ji:        `一乙几刀力匕三上凡也土大川己已巳干弓之仁什今天太巴引戈支斤日比片牙王爿主乏代刊加功包卉占古召叱央它市布平弘弗旦玉申皮矛矢石示氾光列刑危夷妃州帆式戎早旨旭朵此虫血衣圮托汎伸佔但伯伶别判君呀坎壮希庇廷弟彤形杉系邑礽玖些依例刻券刷到制协卓卷呻坤奇奈宗宛延弦或戕旺易昌昆昂明昀昏昕昊昇枝欣版状直知社祀祁纠初采长忻玕佌旻炘炅罕泓泡泛玫表亮侯係前剋则勇垣奐宣巷帝帅建弭彦春昭映昧是星昱架柏矜祉祈祇禹穿竿紂红纪纫紇约紆羿虹衫风玡玠柰祅紈迅巡芝芽指珊玲珍珀玳倡候修刚奚孙家差席师庭时晋晏晃晁书核栩矩砷祕祐祠祟祖神祝祚纺纱纹素索纯纽级紜纳纸纷翁蚩袁袂衽讯託起躬珅倧扆祜紘紓衿衾邪邦那迎返近范茅苓班琉珮珠乾健凰副务勘曼唱唬婉婚将崇常张强彬彩彫旋晚晤晨曹勗桿烯祥票祭絃统扎绍紼细绅组终羚翌翎聆彪蛉被袒袖袍袋责玼珩埏紵紾袗邵述迦迪能情掛採捺涎浅琅球理现剴创劳喉场堤堠媛嵐巽帧弼彭斯普晰晴晶景智曾棕牌番短结绒绝紫絮丝络给绚絳翕蛟裁裂视诊费贵须邰琁珺牋矞絪絜軫陀郎郁送迷迺莫庄莉提扬琪琳琥琴琦琨勤势园戡暗暉暖暄会杨枫照煜牒祺禄禁经绢绥群蜀蛾蜕蜂裟裙补裘装裕试钾雉零琬琰琝媴絻郡通连速造透逢途準猿瑚瑟瑞瑙瑛瑜署僎划寥廖彰截畅碧禎福祸粽绽綰综绰綾绿紧缀网纲綺绸绵綵纶维绪緇綬臧蜜蜻裳裴裹裸製褚豪宾赶闽溒瑄瑋榬綪緁緆緋綖綦蜒裱魠郭都週逸进慢漫瑶琐玛瑰剧刘剑增审影暮槽桨毅奖莹稼缔练纬緻缄缅编缘缎缓緲緹翩蝴蝶褐复褓褊诞赏质辉驻鲁齿摎漻褌褋褙陈乡运游道达违遁撰潜澎璋璃瑾璀剂勋战历晓曇炽御縑縈县縝縉萤融褪裤褫讽醒骇璇璉縕螈錼阳邹远逊遣遥递蒋璟璞励弥戏戴曙墙矫禧禪绩繆缕总纵縴縵翼襄褸辕鍚鄔蟉襁适迁环璦璨断曜璧织缮绕繚绣繒蝉顏题蕥鎱际郑邓选迟薑璿嚥曝牘璽疆祷繫绎绳绘缴襠襟识赞譔还迈邀藏琼劝曦繽继耀腾龄繾饌邈瓏樱缠续边瓔弯禳衬鷚晒缨纤襴鷸蛮纘逻湾缆驪別壯協狀糾長則帥彥紅紀紉約風剛孫師時晉書紡紗紋純紐級納紙紛訊務將張強統紮紹細紳組終責淺現創勞場幀結絨絕絲絡給絢視診費貴須莊揚勢園會楊楓祿經絹綏蛻補裝試鉀連劃暢禍綻綜綽綠緊綴網綱綢綿綸維緒賓趕閩進瑤瑣瑪劇劉劍審槳獎瑩締練緯緘緬編緣緞緩複誕賞質輝駐魯齒陳鄉運遊達違潛劑勳戰曆曉熾禦縣螢褲諷駭陽鄒遠遜遙遞蔣勵彌戲牆矯績縷總縱轅適遷環斷織繕繞繡蟬題際鄭鄧選遲禱繹繩繪繳識贊還邁瓊勸繼騰齡櫻纏續邊彎襯曬纓纖蠻邏灣纜`,
		JiRadical: "",
	},
	ZodiacMice: {
		Name:      ZodiacMice,
		Xi:        "",
		XiRadical: "",
		Ji:        `丁二人士仁仃仇仍今介化午夭孔巴日火以付仕代仙仟卯央平旦未汀伙伊伍休仲件任仰仳份企光全印合地在圭妃寺早旨旬旭朴次竹羊臣行仵价伂位住伫伴佛何估佐佑伺伸似但作伯低伶余佈免坊址坍均坎壮孝宋巫志攸杜赤辛佟抑亚佯依侍佳使供例来佰佩仑佾侑兔味命坪坡坦妹岱幸旺易昌昆昂明昀昏昊昇东杵炎祀卧佼佶侄坵旻炅耵肖拒沫亭亮信侯侠保促俟俊俗俐係俞勉南哇型垠垣垢城姜姿宣封庠律徉春昭映昧是星昱柯柏柳炫为炳炬炯红美订酊香俋昶炷肯洋倍俸倩倖俩值倚倨俱倡候修倪俾伦冤冻卿埂埔埃姬娩宰差徐恙时晋晏晃晁书桓柴氧烊烈留羔耿袁袂酒钉马珋倢埕埒晟邦那迎胡苎苣茉英茆挽珠停假偃偌做伟健偶侦倏冕曼唱域坚堆埠基堵执婚张得从悠教晚晤曹勗欲焉烯瓷祥羚聊袋许顶偈偅偩偮娅欷欸羕邱迪悻傢傅备杰喜尧堪场堤堰报堡复普晰晴晶景智曾棚款焰然善翔贷距辜集冯傌媄焯羢郎郁脩莘莫荷提扬佣傲传仅僇吗塑塘涂塚塔填塌块坞妈微意暗晖暖暄会楠杨歇照煜焕义羡群详路钾铆驰驯渼羟郝连造逢菟署僧像侨境垫墅寿彻榴歌熙熊监绽台赫輓驳墉杩都逸进慢漾漫亿仪僵价侬俭刘墟增坠墩嬉德徵暮样槽楼欧羯蝴卖辉锌养驻驷驶驾驹鲁儆僾儌墡禡羬陈运道达潜骂儒傧壁垦坛壅曆晓昙炽燕熹笃糕羲谘醒骆嶪烨羱阳邹远优偿储壕壑壎曙营灿镁鍚骋骏鲜鲑邬蓨澲燨膳丛垒曜柠欤缮缯题骑镏骐郑邓选薑嚥垄坜曝羶羹辞鹏鄯薘繨邺迈壤曦耀议骞腾骝牺驱蓦欢鑑鞑骄骅惊驿验坝骤驞骥骊`,
		JiRadical: "",
	},
}

// GetZodiac ...
func GetZodiac(c chronos.Calendar) *Zodiac {
	z := chronos.GetZodiac(c.Lunar())
	if v, b := zodiacList[z]; b {
		return &v
	}
	return nil
}

func (z *Zodiac) zodiacJi(character *Character) int {
	if strings.IndexRune(z.Ji, []rune(character.Ch)[0]) != -1 {
		return -3
	}
	return 0
}

func filterZodiac(c chronos.Calendar, chars ...*Character) bool {
	return GetZodiac(c).PointCheck(3, chars...)
}

//PointCheck 检查point
func (z *Zodiac) PointCheck(limit int, chars ...*Character) bool {
	for _, c := range chars {
		if z.Point(c) < limit {
			return false
		}
	}
	return true
}

//Point 喜忌对冲，理论上喜忌都有的话，最好不要选给1，忌给0，喜给5，都没有给3
func (z *Zodiac) Point(character *Character) int {
	dp := 3
	dp += z.zodiacJi(character)
	dp += z.zodiacXi(character)
	return dp
}

func (z *Zodiac) zodiacXi(character *Character) int {
	if strings.IndexRune(z.Xi, []rune(character.Ch)[0]) != -1 {
		return 2
	}
	return 0
}
