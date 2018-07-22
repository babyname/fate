package fate

type Martial struct {
	BiHua     string `bson:"bi_hua"json:"bi_hua"`         //笔画
	SanCai    string `bson:"san_cai"json:"san_cai"`       //三才
	BaZi      string `bson:"ba_zi"json:"ba_zi"`           //八字
	GuaXiang  string `bson:"gua_xiang"json:"gua_xiang"`   //卦象
	TianYun   string `bson:"tian_yun"json:"tian_yun"`     //天运
	ShengXiao string `bson:"sheng_xiao"json:"sheng_xiao"` //生肖
}
