// Package zhouyi 提供周易起卦相关功能。
package zhouyi

// GuaYao 表示卦爻及其吉凶信息。
type GuaYao struct {
	Yao     string `bson:"er_yao"`
	JiXiong string `bson:"er_yao_ji_xiong"`
}
