package mongo_test

import (
	"log"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/godcong/fate/mongo"
)

func TestNaYin(t *testing.T) {
	mongo.Dial("localhost", &mgo.Credential{
		Username: "root",
		Password: "v2RgzSuIaBlx",
	})
	ny := mongo.NaYin{
		GanZhi:   []string{"甲", "子"},
		WuXing:   "金",
		ZhiLiang: "海中金",
		Comment:  `甲子为从革之金，金气散漫，若得到戊申大驿土和癸巳长流水的相助，就会兴发起来。因为戊申是金之临官禄旺，癸巳是金之长生秀气，内藏火土金水生成之道，纳音各有所归，为朝元禄。怕遇见丁卯炉中火，丁酉山下火，戊午天上火克破，多为贫贱或短寿。甲子纳音金，禄官为金，天官藏地元。金溺水下，金死仲中，气泄于子，水旺而金衰，金沉水下而无光辉。须火暖寒体，见火金之成器，水土同宫，印旺亥子，明暗取官，可光耀而名。`,
	}
	mongo.C("nanyin").Insert(&ny)
}

func TestWuXing(t *testing.T) {
	mongo.Dial("localhost", &mgo.Credential{
		Username: "root",
		Password: "v2RgzSuIaBlx",
	})
	wx := mongo.GetWuXing()
	log.Printf("%+v", wx)
}

func TestZhouYi(t *testing.T) {
	mongo.Dial("localhost", &mgo.Credential{
		Username: "root",
		Password: "v2RgzSuIaBlx",
	})

	//insert
	//zy := mongo.ZhouYi{
	//	ShangGua: "乾",
	//	XiaGua:   "乾",
	//	GuaXiang: "乾",
	//	GuaMing:  "乾为天",
	//	GuaYi:    "自强不息",
	//	GuaYun:   "上上卦",
	//	XiangYue: "困龙得水好运交，不由喜气上眉梢，一切谋望皆如意，向后时运渐渐高。这个卦是同卦（下乾上乾）相叠。象征天，喻龙（德才的君子），又象征纯粹的阳和健，表明兴盛强健。乾卦是根据万物变通的道理，以“元、亨、利、贞”为卦辞，示吉祥如意，教导人遵守天道的德行。",
	//}
	//mongo.C("zhouyi").Insert(&zy)

	//get
	//zy := mongo.GetZhouYi()
	//log.Printf("%+v", zy)

}
