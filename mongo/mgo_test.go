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
