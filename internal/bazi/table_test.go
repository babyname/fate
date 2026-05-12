package bazi

import (
	"testing"
)

var expectedHiddenStems = map[string]map[string]float64{
	"子": {"癸": 1.0},
	"丑": {"己": 0.50, "癸": 0.30, "辛": 0.20},
	"寅": {"甲": 0.50, "丙": 0.30, "戊": 0.20},
	"卯": {"乙": 1.0},
	"辰": {"戊": 0.50, "乙": 0.30, "癸": 0.20},
	"巳": {"丙": 0.50, "庚": 0.30, "戊": 0.20},
	"午": {"丁": 0.70, "己": 0.30},
	"未": {"己": 0.50, "丁": 0.30, "乙": 0.20},
	"申": {"庚": 0.50, "壬": 0.30, "戊": 0.20},
	"酉": {"辛": 1.0},
	"戌": {"戊": 0.50, "辛": 0.30, "丁": 0.20},
	"亥": {"壬": 0.70, "甲": 0.30},
}

var diZhiNames = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}

func TestDizhi_HiddenStemCount(t *testing.T) {
	for i, name := range diZhiNames {
		expected := expectedHiddenStems[name]
		actual := len(dizhi[i])
		if actual != len(expected) {
			t.Errorf("地支 %s(索引%d): 期望 %d 个藏干, 实际 %d 个", name, i, len(expected), actual)
		}
	}
}

func TestDizhi_MissingStemsFixed(t *testing.T) {
	missingChecks := []struct {
		diZhi   string
		diIndex int
		stem    string
	}{
		{"寅", 2, "戊"},
		{"巳", 5, "戊"},
		{"午", 6, "己"},
		{"申", 8, "戊"},
	}
	for _, tc := range missingChecks {
		dz := dizhi[tc.diIndex]
		if _, ok := dz[tc.stem]; !ok {
			t.Errorf("地支 %s(索引%d): 缺失藏干 %s (Bug #139 未修复)", tc.diZhi, tc.diIndex, tc.stem)
		}
	}
}

func TestDizhi_AllStemsPresent(t *testing.T) {
	for i, name := range diZhiNames {
		dz := dizhi[i]
		expected := expectedHiddenStems[name]
		for stem := range expected {
			if _, ok := dz[stem]; !ok {
				t.Errorf("地支 %s(索引%d): 缺失藏干 %s", name, i, stem)
			}
		}
		for stem := range dz {
			if _, ok := expected[stem]; !ok {
				t.Errorf("地支 %s(索引%d): 多余藏干 %s", name, i, stem)
			}
		}
	}
}

func TestDizhi_StemValueFormula(t *testing.T) {
	for i, name := range diZhiNames {
		dz := dizhi[i]
		expected := expectedHiddenStems[name]
		for stem, proportion := range expected {
			values, ok := dz[stem]
			if !ok {
				t.Fatalf("地支 %s(索引%d): 藏干 %s 不存在", name, i, stem)
			}
			if len(values) != 12 {
				t.Errorf("地支 %s 藏干 %s: 期望12个月数据, 实际%d个月", name, stem, len(values))
				continue
			}
			tIdx, ok := tianIndex[stem]
			if !ok {
				t.Fatalf("天干 %s 不在 tianIndex 中", stem)
			}
			for month := 0; month < 12; month++ {
				expectedValue := int(float64(tiangan[month][tIdx]) * proportion)
				if values[month] != expectedValue {
					t.Errorf("地支 %s 藏干 %s 月份%d: 期望 %d (tiangan[%d][%d]=%d * %.2f), 实际 %d",
						name, stem, month, expectedValue, month, tIdx, tiangan[month][tIdx], proportion, values[month])
				}
			}
		}
	}
}

func TestDizhi_TotalMatchesFormula(t *testing.T) {
	for i, name := range diZhiNames {
		dz := dizhi[i]
		expected := expectedHiddenStems[name]
		for month := 0; month < 12; month++ {
			expectedTotal := 0
			for stem, proportion := range expected {
				expectedTotal += int(float64(tiangan[month][tianIndex[stem]]) * proportion)
			}
			actualTotal := 0
			for _, values := range dz {
				actualTotal += values[month]
			}
			if actualTotal != expectedTotal {
				t.Errorf("地支 %s 月份%d: 藏干总和 %d != 公式计算总和 %d",
					name, month, actualTotal, expectedTotal)
			}
		}
	}
}

func TestDizhi_ProportionSum(t *testing.T) {
	for _, name := range diZhiNames {
		propSum := 0.0
		for _, prop := range expectedHiddenStems[name] {
			propSum += prop
		}
		if propSum != 1.0 {
			t.Errorf("地支 %s: 藏干比例总和 = %.2f, 期望 1.0", name, propSum)
		}
	}
}

func TestDizhi_StemBelongsToCorrectWuXing(t *testing.T) {
	for i, name := range diZhiNames {
		dz := dizhi[i]
		diZhiWuXing := wuXingDiZhi[name]
		hasMainElement := false
		for stem := range dz {
			stemWuXing := wuXingTianGan[stem]
			if stemWuXing == diZhiWuXing {
				hasMainElement = true
				break
			}
		}
		if !hasMainElement {
			t.Errorf("地支 %s(五行%s): 没有本气藏干属于相同五行", name, diZhiWuXing)
		}
	}
}

func TestDizhi_NoZeroValues(t *testing.T) {
	for i, name := range diZhiNames {
		dz := dizhi[i]
		for stem, values := range dz {
			for month, v := range values {
				if v <= 0 {
					t.Errorf("地支 %s 藏干 %s 月份%d: 强度值 %d <= 0", name, stem, month, v)
				}
			}
		}
	}
}

func TestDizhi_MainStemIsStrongest(t *testing.T) {
	for i, name := range diZhiNames {
		dz := dizhi[i]
		expected := expectedHiddenStems[name]
		if len(expected) <= 1 {
			continue
		}
		mainStem := ""
		maxProp := 0.0
		for stem, prop := range expected {
			if prop > maxProp {
				maxProp = prop
				mainStem = stem
			}
		}
		for month := 0; month < 12; month++ {
			mainVal := dz[mainStem][month]
			for stem, values := range dz {
				if stem == mainStem {
					continue
				}
				if values[month] > mainVal {
					t.Errorf("地支 %s 月份%d: 本气 %s=%d 被非本气 %s=%d 超过",
						name, month, mainStem, mainVal, stem, values[month])
				}
			}
		}
	}
}

func TestDizhi_SpecificValues_Bug139(t *testing.T) {
	type check struct {
		diZhi  string
		idx    int
		stem   string
		month  int
		expect int
	}
	checks := []check{
		{"寅", 2, "甲", 0, 600},
		{"寅", 2, "甲", 2, 570},
		{"寅", 2, "丙", 0, 300},
		{"寅", 2, "丙", 2, 360},
		{"寅", 2, "戊", 0, 200},
		{"寅", 2, "戊", 6, 240},
		{"巳", 5, "丙", 0, 500},
		{"巳", 5, "丙", 2, 600},
		{"巳", 5, "庚", 0, 300},
		{"巳", 5, "戊", 0, 200},
		{"巳", 5, "戊", 6, 240},
		{"午", 6, "丁", 0, 700},
		{"午", 6, "丁", 2, 840},
		{"午", 6, "己", 0, 300},
		{"午", 6, "己", 6, 360},
		{"申", 8, "庚", 0, 500},
		{"申", 8, "庚", 9, 600},
		{"申", 8, "壬", 0, 360},
		{"申", 8, "戊", 0, 200},
		{"申", 8, "戊", 6, 240},
	}
	for _, c := range checks {
		dz := dizhi[c.idx]
		vals, ok := dz[c.stem]
		if !ok {
			t.Errorf("地支 %s(索引%d): 藏干 %s 不存在", c.diZhi, c.idx, c.stem)
			continue
		}
		if vals[c.month] != c.expect {
			t.Errorf("地支 %s 藏干 %s 月份%d: 期望 %d, 实际 %d",
				c.diZhi, c.stem, c.month, c.expect, vals[c.month])
		}
	}
}

func TestDizhi_DataEntryErrorsFixed(t *testing.T) {
	type check struct {
		diZhi  string
		idx    int
		stem   string
		month  int
		expect int
		desc   string
	}
	checks := []check{
		{"丑", 1, "辛", 4, 220, "丑辛月4: 230→220"},
		{"丑", 1, "辛", 9, 240, "丑辛月9: 248→240"},
		{"辰", 4, "癸", 5, 212, "辰癸月5: 200→212"},
		{"辰", 4, "戊", 5, 570, "辰戊月5: 600→570"},
		{"未", 7, "丁", 11, 300, "未丁月11: 318→300"},
		{"戌", 10, "丁", 11, 200, "戌丁月11: 212→200"},
		{"亥", 11, "壬", 10, 742, "亥壬月10: 724→742"},
	}
	for _, c := range checks {
		dz := dizhi[c.idx]
		vals, ok := dz[c.stem]
		if !ok {
			t.Errorf("%s: 藏干 %s 不存在", c.desc, c.stem)
			continue
		}
		if vals[c.month] != c.expect {
			t.Errorf("%s: 期望 %d, 实际 %d", c.desc, c.expect, vals[c.month])
		}
	}
}

func TestTiangan_Completeness(t *testing.T) {
	if len(tiangan) != 12 {
		t.Errorf("天干强度表: 期望12个月, 实际%d个月", len(tiangan))
	}
	for month, row := range tiangan {
		if len(row) != 10 {
			t.Errorf("天干强度表月份%d: 期望10个天干, 实际%d个", month, len(row))
		}
	}
}

func TestDiIndex_Completeness(t *testing.T) {
	expectedDi := map[string]int{
		"子": 0, "丑": 1, "寅": 2, "卯": 3, "辰": 4, "巳": 5,
		"午": 6, "未": 7, "申": 8, "酉": 9, "戌": 10, "亥": 11,
	}
	for k, v := range expectedDi {
		if diIndex[k] != v {
			t.Errorf("diIndex[%s]: 期望 %d, 实际 %d", k, v, diIndex[k])
		}
	}
	if len(diIndex) != 12 {
		t.Errorf("diIndex: 期望12个地支, 实际%d个", len(diIndex))
	}
}

func TestTianIndex_Completeness(t *testing.T) {
	expectedTian := map[string]int{
		"甲": 0, "乙": 1, "丙": 2, "丁": 3, "戊": 4,
		"己": 5, "庚": 6, "辛": 7, "壬": 8, "癸": 9,
	}
	for k, v := range expectedTian {
		if tianIndex[k] != v {
			t.Errorf("tianIndex[%s]: 期望 %d, 实际 %d", k, v, tianIndex[k])
		}
	}
	if len(tianIndex) != 10 {
		t.Errorf("tianIndex: 期望10个天干, 实际%d个", len(tianIndex))
	}
}

func TestWuXingTianGan_Completeness(t *testing.T) {
	stems := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	for _, s := range stems {
		if _, ok := wuXingTianGan[s]; !ok {
			t.Errorf("wuXingTianGan: 缺失天干 %s", s)
		}
	}
}

func TestWuXingDiZhi_Completeness(t *testing.T) {
	for _, name := range diZhiNames {
		if _, ok := wuXingDiZhi[name]; !ok {
			t.Errorf("wuXingDiZhi: 缺失地支 %s", name)
		}
	}
}
