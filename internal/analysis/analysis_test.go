package analysis

import (
	"testing"
)

func TestYinYangWuXingAttr(t *testing.T) {
	tests := []struct {
		stroke   int
		expected string
	}{
		{1, "阳木"},
		{2, "阴木"},
		{3, "阳火"},
		{4, "阴火"},
		{5, "阳土"},
		{6, "阴土"},
		{7, "阳金"},
		{8, "阴金"},
		{9, "阳水"},
		{10, "阴水"},
		{11, "阳木"},
		{12, "阴木"},
		{13, "阳火"},
	}
	for _, tt := range tests {
		result := yinYangWuXingAttr(tt.stroke)
		if result != tt.expected {
			t.Errorf("yinYangWuXingAttr(%d) = %s, want %s", tt.stroke, result, tt.expected)
		}
	}
}

func TestSanCaiDetail(t *testing.T) {
	detail := getSanCaiDetail("木木木")
	if detail == "" || detail == "暂无解析" {
		t.Error("Should have detail for 木木木")
	}
	detail = getSanCaiDetail("火火火")
	if detail == "" || detail == "暂无解析" {
		t.Error("Should have detail for 火火火")
	}
}

func TestJiChuYun(t *testing.T) {
	result := getJiChuYun("木", "木")
	if result == "" || result == "暂无解析" {
		t.Error("Should have jichuyun for 木木")
	}
}

func TestChengGongYun(t *testing.T) {
	result := getChengGongYun("木", "火")
	if result == "" || result == "暂无解析" {
		t.Error("Should have chenggongyun for 木火")
	}
}

func TestRenJiGuanXi(t *testing.T) {
	result := getRenJiGuanXi("金", "水")
	if result == "" || result == "暂无解析" {
		t.Error("Should have renjiguanxi for 金水")
	}
}

func TestZodiacWuXing(t *testing.T) {
	tests := []struct {
		zodiac   string
		expected string
	}{
		{"鼠", "水"},
		{"虎", "木"},
		{"龙", "土"},
		{"蛇", "火"},
		{"猴", "金"},
		{"", ""},
	}
	for _, tt := range tests {
		result := getZodiacWuXing(tt.zodiac)
		if result != tt.expected {
			t.Errorf("getZodiacWuXing(%s) = %s, want %s", tt.zodiac, result, tt.expected)
		}
	}
}

func TestCalcZhouYi(t *testing.T) {
	result := CalcZhouYi(11, 0, 13, 12)
	if result == nil {
		t.Fatal("CalcZhouYi should not return nil")
	}
	if result.BenGuaName == "" {
		t.Error("BenGuaName should not be empty")
	}
	if result.BianGuaName == "" {
		t.Error("BianGuaName should not be empty")
	}
	t.Logf("本卦: %s（%s）, 变卦: %s, 动爻: %s, 大象: %s",
		result.BenGuaName, result.BenGuaJiXiong, result.BianGuaName, result.DongYaoDesc, result.DaXiang)
}
