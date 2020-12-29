package fate

//Martial six martials
type Martial struct {
	BiHua     bool `bson:"bi_hua" json:"bi_hua"`         //笔画
	SanCai    bool `bson:"san_cai" json:"san_cai"`       //三才
	BaZi      bool `bson:"ba_zi" json:"ba_zi"`           //八字
	GuaXiang  bool `bson:"gua_xiang" json:"gua_xiang"`   //卦象
	TianYun   bool `bson:"tian_yun" json:"tian_yun"`     //天运
	ShengXiao bool `bson:"sheng_xiao" json:"sheng_xiao"` //生肖
}
