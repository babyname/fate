package analysis

// GuaDetail 表示卦象的详细解读信息
type GuaDetail struct {
	DaXiang   string
	YunShi    string
	ShiYe     string
	JingShang string
	QiuMing   string
	HunLian   string
	JueCe     string
}

var guaDetailMap = map[string]*GuaDetail{
	"乾为天":  {DaXiang: "天行健，君子以自强不息", YunShi: "大吉，万事如意，事业兴旺", ShiYe: "事业顺利，有贵人相助，可成大业", JingShang: "经商顺利，财运亨通，可大胆投资", QiuMing: "名声远扬，功名有望", HunLian: "婚姻美满，门当户对", JueCe: "刚健中正，宜积极进取"},
	"坤为地":  {DaXiang: "地势坤，君子以厚德载物", YunShi: "平稳，宜守不宜攻", ShiYe: "宜守成，不宜冒进，以柔顺为上", JingShang: "宜合作，不宜独断，以和为贵", QiuMing: "功名较迟，需耐心等待", HunLian: "婚姻平稳，宜柔顺相处", JueCe: "柔顺守正，宜退守不宜进取"},
	"震为雷":  {DaXiang: "洊雷震，君子以恐惧修省", YunShi: "变动不居，有惊无险", ShiYe: "事业有变，需谨慎应对", JingShang: "市场波动，宜保守经营", QiuMing: "名声有起伏，需修身养性", HunLian: "感情有波折，需坦诚相待", JueCe: "临危不惧，谨慎行事"},
	"巽为风":  {DaXiang: "随风巽，君子以申命行事", YunShi: "顺风吹火，事半功倍", ShiYe: "宜顺水推舟，借势而为", JingShang: "宜与人合作，不宜独断", QiuMing: "名声渐起，宜谦逊待人", HunLian: "感情顺利，宜以诚相待", JueCe: "谦逊行事，宜柔不宜刚"},
	"坎为水":  {DaXiang: "水洊至习坎，君子以常德行习教事", YunShi: "险中求通，需耐心等待", ShiYe: "事业有险，需坚韧不拔", JingShang: "财运不佳，宜守不宜攻", QiuMing: "功名受阻，需坚持不懈", HunLian: "感情有阻，需耐心等待", JueCe: "处险不惊，以诚待人"},
	"离为火":  {DaXiang: "明两作离，大人以继明照于四方", YunShi: "光明磊落，前途光明", ShiYe: "事业有成，宜发扬光大", JingShang: "财运亨通，宜光明正大经营", QiuMing: "功名有望，名声远播", HunLian: "感情热烈，宜理性对待", JueCe: "光明正大，宜进不宜退"},
	"艮为山":  {DaXiang: "兼山艮，君子以思不出其位", YunShi: "宜止则止，安守本分", ShiYe: "宜守成，不宜冒进", JingShang: "宜稳扎稳打，不宜冒险", QiuMing: "功名平稳，需安分守己", HunLian: "感情平稳，宜安守本分", JueCe: "知止而后安，宜静不宜动"},
	"兑为泽":  {DaXiang: "丽泽兑，君子以朋友讲习", YunShi: "喜悦和乐，人际关系佳", ShiYe: "事业顺遂，有朋友相助", JingShang: "财运不错，宜和气生财", QiuMing: "名声渐起，宜广结善缘", HunLian: "感情甜蜜，宜珍惜眼前", JueCe: "和悦待人，宜喜不宜忧"},
	"天水讼":  {DaXiang: "天与水违行讼，君子以作事谋始", YunShi: "有争讼之象，宜和解", ShiYe: "事业有争端，宜冷静处理", JingShang: "商业纠纷，宜和解不宜诉讼", QiuMing: "功名有阻，需忍让", HunLian: "感情有争执，需互相理解", JueCe: "避免争讼，以和为贵"},
	"水天需":  {DaXiang: "云上于天需，君子以饮食宴乐", YunShi: "等待时机，不宜急进", ShiYe: "需耐心等待，时机未到", JingShang: "投资需等待，不宜急进", QiuMing: "功名有待，需积蓄力量", HunLian: "感情需等待，不宜急躁", JueCe: "耐心等待，时机自到"},
	"天雷无妄": {DaXiang: "天下雷行无妄，君子以对时育万物", YunShi: "无妄之灾，需谨慎行事", ShiYe: "事业有意外，需谨慎应对", JingShang: "不宜投机，宜脚踏实地", QiuMing: "功名有阻，需守正道", HunLian: "感情需真诚，不宜虚妄", JueCe: "守正不妄，顺其自然"},
	"雷天大壮": {DaXiang: "雷在天上大壮，君子以非礼弗履", YunShi: "气势壮盛，但需守礼", ShiYe: "事业壮盛，但需谦虚", JingShang: "财运壮盛，但不宜贪婪", QiuMing: "功名壮盛，但需守礼", HunLian: "感情热烈，但需理性", JueCe: "壮而守礼，不可逞强"},
	"天风姤":  {DaXiang: "天下有风姤，君子以施命诰四方", YunShi: "偶遇之象，需谨慎选择", ShiYe: "有新的机遇，需明辨是非", JingShang: "有新的合作，需审慎选择", QiuMing: "有新的机会，需把握", HunLian: "有新的邂逅，需谨慎", JueCe: "遇事需明辨，不可轻信"},
	"风天小畜": {DaXiang: "风行天上小畜，君子以懿文德", YunShi: "力量不足，需蓄积", ShiYe: "事业需积累，不宜冒进", JingShang: "资金不足，宜小规模经营", QiuMing: "功名尚需时日", HunLian: "感情需培养，不宜急躁", JueCe: "蓄积力量，以柔克刚"},
	"天山遁":  {DaXiang: "天下有山遁，君子以远小人", YunShi: "退避之象，宜隐忍", ShiYe: "宜暂时退守，不宜冒进", JingShang: "宜收缩战线，保存实力", QiuMing: "功名有阻，宜隐忍", HunLian: "感情需冷静，不宜冲动", JueCe: "退而保全，待时而动"},
	"山天大畜": {DaXiang: "天在山中大畜，君子以多识前言往行", YunShi: "积蓄力量，大有可为", ShiYe: "积蓄力量，可成大业", JingShang: "资金充裕，可大展宏图", QiuMing: "功名有望，需广学多识", HunLian: "感情稳定，宜珍惜", JueCe: "蓄势待发，厚积薄发"},
	"天泽履":  {DaXiang: "上天下泽履，君子以辩上下定民志", YunShi: "如履薄冰，需谨慎", ShiYe: "需谨慎行事，步步为营", JingShang: "宜稳健经营，不宜冒险", QiuMing: "功名需循规蹈矩", HunLian: "感情需以礼相待", JueCe: "循礼而行，谨慎为上"},
	"泽天夬":  {DaXiang: "泽上于天夬，君子以施禄及下", YunShi: "决断之象，需果断", ShiYe: "需果断决策，不可犹豫", JingShang: "需果断投资，把握时机", QiuMing: "功名可成，需果断", HunLian: "感情需果断抉择", JueCe: "果断决策，不可优柔"},
	"地雷复":  {DaXiang: "雷在地中复，君子以至日闭关", YunShi: "一阳来复，否极泰来", ShiYe: "事业转机，宜把握时机", JingShang: "财运转好，宜把握机会", QiuMing: "功名有望，需努力", HunLian: "感情回春，宜珍惜", JueCe: "否极泰来，宜积极进取"},
	"雷地豫":  {DaXiang: "雷出地奋豫，君子以作乐崇德", YunShi: "安乐之象，但需警惕", ShiYe: "事业安逸，但需居安思危", JingShang: "财运尚可，但不宜懈怠", QiuMing: "功名平稳，需继续努力", HunLian: "感情和乐，需珍惜", JueCe: "居安思危，乐而不淫"},
	"地水师":  {DaXiang: "地中有水师，君子以容民畜众", YunShi: "统率之象，需以德服人", ShiYe: "宜带领团队，以德服人", JingShang: "宜合作经营，以和为贵", QiuMing: "功名在军政方面有发展", HunLian: "感情需以诚相待", JueCe: "以德服人，不可用强"},
	"水地比":  {DaXiang: "地上有水比，君子以建万国亲诸侯", YunShi: "亲比之象，宜团结", ShiYe: "宜团结合作，共同发展", JingShang: "宜合伙经营，互利互惠", QiuMing: "功名需靠人脉", HunLian: "感情亲密，宜珍惜", JueCe: "亲比团结，以和为贵"},
	"地风升":  {DaXiang: "地中生木升，君子以顺德积小以高大", YunShi: "上升之象，步步高升", ShiYe: "事业上升，宜循序渐进", JingShang: "财运上升，宜稳扎稳打", QiuMing: "功名上升，需继续努力", HunLian: "感情升温，宜珍惜", JueCe: "循序渐进，步步高升"},
	"风地观":  {DaXiang: "风行地上观，君子以省方观民设教", YunShi: "观察之象，宜审时度势", ShiYe: "宜观察形势，不宜轻举妄动", JingShang: "宜观望市场，不宜盲目投资", QiuMing: "功名需观察时机", HunLian: "感情需观察，不宜冲动", JueCe: "审时度势，谋定后动"},
	"地山谦":  {DaXiang: "地中有山谦，君子以裒多益寡", YunShi: "谦虚之象，以退为进", ShiYe: "宜谦虚谨慎，以退为进", JingShang: "宜薄利多销，以德经商", QiuMing: "功名需谦虚，不可骄傲", HunLian: "感情需谦让，以和为贵", JueCe: "谦虚谨慎，以退为进"},
	"山地剥":  {DaXiang: "山附于地剥，君子以顺天止上", YunShi: "剥落之象，需守成", ShiYe: "事业有损，宜守不宜攻", JingShang: "财运不佳，宜保守经营", QiuMing: "功名有阻，需忍耐", HunLian: "感情有变，需冷静", JueCe: "守成待变，不宜冒进"},
	"地泽临":  {DaXiang: "泽上有地临，君子以教思无穷容保民无疆", YunShi: "亲临之象，宜关怀", ShiYe: "宜亲力亲为，关怀下属", JingShang: "宜亲临一线，了解市场", QiuMing: "功名有望，需亲力亲为", HunLian: "感情亲近，宜珍惜", JueCe: "亲临关怀，以德服人"},
	"泽地萃":  {DaXiang: "泽上于地萃，君子以除戎器戒不虞", YunShi: "聚集之象，宜团结", ShiYe: "宜聚集人才，共同发展", JingShang: "宜集中资源，合力经营", QiuMing: "功名需众人相助", HunLian: "感情融洽，宜珍惜", JueCe: "聚集力量，团结一致"},
	"雷水解":  {DaXiang: "雷雨作解，君子以赦过宥罪", YunShi: "解除困难，柳暗花明", ShiYe: "困难解除，宜积极进取", JingShang: "困境解除，宜把握机会", QiuMing: "功名有望，需努力", HunLian: "感情转好，宜珍惜", JueCe: "困难解除，宜积极进取"},
	"水雷屯":  {DaXiang: "云雷屯，君子以经纶", YunShi: "创业维艰，需坚持", ShiYe: "事业初创，困难重重", JingShang: "创业艰难，需坚持不懈", QiuMing: "功名尚早，需努力", HunLian: "感情初创，需耐心培养", JueCe: "创业维艰，坚持不懈"},
	"雷风恒":  {DaXiang: "雷风恒，君子以立不易方", YunShi: "恒久不变，宜守常", ShiYe: "宜持之以恒，不可半途而废", JingShang: "宜长期经营，不可朝三暮四", QiuMing: "功名需持之以恒", HunLian: "感情需长久经营", JueCe: "持之以恒，不可变节"},
	"风雷益":  {DaXiang: "风雷益，君子以见善则迁有过则改", YunShi: "增益之象，宜进取", ShiYe: "事业增益，宜积极进取", JingShang: "财运增益，宜扩大经营", QiuMing: "功名增益，需努力", HunLian: "感情增益，宜珍惜", JueCe: "见善则迁，积极进取"},
	"雷山小过": {DaXiang: "山上有雷小过，君子以行过乎恭丧过乎哀用过乎俭", YunShi: "小有过越，需谨慎", ShiYe: "小有失误，需谨慎行事", JingShang: "小有损失，需谨慎经营", QiuMing: "功名小有阻碍", HunLian: "感情小有波折", JueCe: "小过宜改，谨慎行事"},
	"山雷颐":  {DaXiang: "山下有雷颐，君子以慎言语节饮食", YunShi: "颐养之象，宜修养", ShiYe: "宜修养身心，积蓄力量", JingShang: "宜谨慎理财，量入为出", QiuMing: "功名需修养", HunLian: "感情需细心呵护", JueCe: "慎言节欲，修养身心"},
	"雷泽归妹": {DaXiang: "泽上有雷归妹，君子以永终知敝", YunShi: "归妹之象，需慎重", ShiYe: "需慎重选择，不可轻率", JingShang: "投资需谨慎，不宜轻率", QiuMing: "功名需正当途径", HunLian: "感情需慎重，不可轻率", JueCe: "慎重选择，不可轻率"},
	"泽雷随":  {DaXiang: "泽中有雷随，君子以向晦入宴息", YunShi: "随从之象，宜顺势", ShiYe: "宜顺应时势，不可逆行", JingShang: "宜跟随大势，不可逆势而为", QiuMing: "功名需顺势而为", HunLian: "感情需顺其自然", JueCe: "顺势而为，不可强求"},
	"水风井":  {DaXiang: "木上有水井，君子以劳民劝相", YunShi: "井养之象，宜修身", ShiYe: "宜修身养性，服务他人", JingShang: "宜诚信经营，服务客户", QiuMing: "功名需修身", HunLian: "感情需用心经营", JueCe: "修身养性，服务他人"},
	"风水涣":  {DaXiang: "风行水上涣，君子以享帝立庙", YunShi: "涣散之象，需凝聚", ShiYe: "人心涣散，需凝聚力量", JingShang: "团队涣散，需加强管理", QiuMing: "功名有散，需努力", HunLian: "感情有散，需珍惜", JueCe: "凝聚人心，团结一致"},
	"水山蹇":  {DaXiang: "山上有水蹇，君子以反身修德", YunShi: "艰难之象，需修身", ShiYe: "事业艰难，需修身养德", JingShang: "经营困难，需节俭持家", QiuMing: "功名艰难，需坚持", HunLian: "感情艰难，需忍耐", JueCe: "修身养德，以德服人"},
	"山水蒙":  {DaXiang: "山下出泉蒙，君子以果行育德", YunShi: "启蒙之象，宜学习", ShiYe: "宜学习进修，提升能力", JingShang: "宜学习经营之道", QiuMing: "功名需学习", HunLian: "感情需学习相处之道", JueCe: "启蒙学习，提升自我"},
	"水泽节":  {DaXiang: "泽上有水节，君子以制数度议德行", YunShi: "节制之象，宜守度", ShiYe: "宜节制有度，不可过度", JingShang: "宜量入为出，不可奢侈", QiuMing: "功名需守度", HunLian: "感情需节制", JueCe: "节制有度，不可过度"},
	"泽水困":  {DaXiang: "泽无水困，君子以致命遂志", YunShi: "困顿之象，需坚守", ShiYe: "事业困顿，需坚守信念", JingShang: "经营困难，需坚持", QiuMing: "功名受阻，需坚持", HunLian: "感情受困，需忍耐", JueCe: "坚守信念，不可放弃"},
	"火雷噬嗑": {DaXiang: "雷电噬嗑，君子以明罚敕法", YunShi: "咬合之象，宜决断", ShiYe: "需果断处理问题", JingShang: "需果断决策，不可犹豫", QiuMing: "功名需决断", HunLian: "感情需果断处理问题", JueCe: "果断决断，不可犹豫"},
	"雷火丰":  {DaXiang: "雷电皆至丰，君子以折狱致刑", YunShi: "丰盛之象，宜把握", ShiYe: "事业丰盛，宜把握时机", JingShang: "财运丰盛，宜把握机会", QiuMing: "功名丰盛，需谦虚", HunLian: "感情丰盛，宜珍惜", JueCe: "丰盛之时，需谦虚谨慎"},
	"火风鼎":  {DaXiang: "木上有火鼎，君子以正位凝命", YunShi: "鼎新之象，宜革新", ShiYe: "宜改革创新，开创新局面", JingShang: "宜更新经营模式", QiuMing: "功名鼎新，需创新", HunLian: "感情需创新相处方式", JueCe: "革故鼎新，开创未来"},
	"风火家人": {DaXiang: "风自火出家人，君子以言有物而行有恒", YunShi: "家庭之象，宜和睦", ShiYe: "宜家庭和睦，事业顺遂", JingShang: "宜家族经营，和气生财", QiuMing: "功名需家庭支持", HunLian: "感情和睦，宜珍惜", JueCe: "家庭和睦，事业顺遂"},
	"火山旅":  {DaXiang: "山上有火旅，君子以明慎用刑而不留狱", YunShi: "旅行之象，宜谨慎", ShiYe: "宜外出发展，但需谨慎", JingShang: "宜开拓外地市场", QiuMing: "功名在外地有发展", HunLian: "感情有分离之象", JueCe: "旅居在外，需谨慎行事"},
	"山火贲":  {DaXiang: "山下有火贲，君子以明庶政无敢折狱", YunShi: "文饰之象，宜修饰", ShiYe: "宜注重形象，提升品质", JingShang: "宜注重品牌，提升形象", QiuMing: "功名需文饰", HunLian: "感情需用心经营", JueCe: "注重修饰，提升品质"},
	"火泽睽":  {DaXiang: "上火下泽睽，君子以同而异", YunShi: "背离之象，需沟通", ShiYe: "人际关系有矛盾，需沟通", JingShang: "合作有分歧，需协商", QiuMing: "功名有阻碍", HunLian: "感情有分歧，需沟通", JueCe: "求同存异，沟通协商"},
	"泽火革":  {DaXiang: "泽中有火革，君子以治历明时", YunShi: "变革之象，宜革新", ShiYe: "宜改革创新，顺应时势", JingShang: "宜变革经营模式", QiuMing: "功名需变革", HunLian: "感情需改变相处方式", JueCe: "顺应时势，勇于变革"},
	"火天大有": {DaXiang: "火在天上大有，君子以遏恶扬善顺天休命", YunShi: "大有之象，大吉大利", ShiYe: "事业大成，前途光明", JingShang: "财运大好，可大展宏图", QiuMing: "功名大就，名利双收", HunLian: "感情美满，天作之合", JueCe: "大有收获，宜谦虚谨慎"},
	"天火同人": {DaXiang: "天与火同人，君子以类族辨物", YunShi: "团结之象，宜合作", ShiYe: "宜团结合作，共同发展", JingShang: "宜合伙经营，互利互惠", QiuMing: "功名需众人相助", HunLian: "感情和谐，志同道合", JueCe: "团结合作，共同发展"},
	"火地晋":  {DaXiang: "明出地上晋，君子以自照明德", YunShi: "晋升之象，宜进取", ShiYe: "事业晋升，宜积极进取", JingShang: "财运上升，宜扩大经营", QiuMing: "功名晋升，前途光明", HunLian: "感情升温，宜珍惜", JueCe: "积极进取，前途光明"},
	"地火明夷": {DaXiang: "明入地中明夷，君子以莅众用晦而明", YunShi: "光明受伤，宜韬晦", ShiYe: "宜韬光养晦，不宜张扬", JingShang: "宜低调经营，不宜张扬", QiuMing: "功名有阻，需忍耐", HunLian: "感情需低调处理", JueCe: "韬光养晦，等待时机"},
	"火水未济": {DaXiang: "火在水上未济，君子以慎辨物居方", YunShi: "未完成之象，需努力", ShiYe: "事业未成，需继续努力", JingShang: "经营未成，需坚持不懈", QiuMing: "功名未就，需努力", HunLian: "感情未定，需耐心", JueCe: "坚持不懈，终有所成"},
	"水火既济": {DaXiang: "水在火上既济，君子以思患而豫防之", YunShi: "已完成之象，需防患", ShiYe: "事业已成，需居安思危", JingShang: "经营有成，需防患未然", QiuMing: "功名已就，需谦虚", HunLian: "感情稳定，需珍惜", JueCe: "居安思危，防患未然"},
	"山风蛊":  {DaXiang: "山下有风蛊，君子以振民育德", YunShi: "腐败之象，需革新", ShiYe: "事业有弊，需改革", JingShang: "经营有弊，需改革", QiuMing: "功名需革新", HunLian: "感情有问题，需沟通", JueCe: "革除弊端，勇于改革"},
	"风山渐":  {DaXiang: "山上有木渐，君子以居贤德善俗", YunShi: "渐进之象，宜循序渐进", ShiYe: "宜循序渐进，不可急躁", JingShang: "宜稳步发展，不可冒进", QiuMing: "功名渐进，需耐心", HunLian: "感情渐入佳境", JueCe: "循序渐进，不可急躁"},
	"山泽损":  {DaXiang: "山下有泽损，君子以惩忿窒欲", YunShi: "减损之象，宜节俭", ShiYe: "宜减少开支，节俭经营", JingShang: "宜减少投资，保守经营", QiuMing: "功名有损，需忍耐", HunLian: "感情有损，需珍惜", JueCe: "减损欲望，节俭持家"},
	"泽山咸":  {DaXiang: "山上有泽咸，君子以虚受人", YunShi: "感应之象，宜沟通", ShiYe: "宜与人沟通，互相理解", JingShang: "宜与客户沟通，了解需求", QiuMing: "功名需人脉", HunLian: "感情和谐，心心相印", JueCe: "虚心接受，互相理解"},
	"泽风大过": {DaXiang: "泽灭木大过，君子以独立不惧遁世无闷", YunShi: "过度之象，需节制", ShiYe: "做事过度，需节制", JingShang: "投资过度，需节制", QiuMing: "功名有过，需谦虚", HunLian: "感情过度，需理性", JueCe: "适度而行，不可过度"},
	"风泽中孚": {DaXiang: "泽上有风中孚，君子以议狱缓死", YunShi: "诚信之象，宜守信", ShiYe: "宜诚信经营，以信立业", JingShang: "宜诚信经商，以信为本", QiuMing: "功名需诚信", HunLian: "感情需真诚", JueCe: "诚信为本，以信立业"},
	"地天泰":  {DaXiang: "天地交泰，君子以裁成天地之道辅相天地之宜", YunShi: "通泰之象，大吉大利", ShiYe: "事业通达，万事亨通", JingShang: "财运亨通，宜大胆经营", QiuMing: "功名通达，前途无量", HunLian: "感情和谐，天作之合", JueCe: "通泰和顺，宜积极进取"},
	"天地否":  {DaXiang: "天地不交否，君子以俭德辟难不可荣以禄", YunShi: "闭塞之象，需忍耐", ShiYe: "事业受阻，宜守不宜攻", JingShang: "经营困难，宜保守", QiuMing: "功名有阻，需忍耐", HunLian: "感情有隔，需沟通", JueCe: "守正待时，不可冒进"},
}

func getGuaDetail(guaMing string) *GuaDetail {
	if v, ok := guaDetailMap[guaMing]; ok {
		return v
	}
	return nil
}
