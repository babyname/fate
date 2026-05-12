# 五行喜忌实现

## 实现总览

五行喜忌实现包括：
- 五行统计函数
- 生扶力量计算
- 克泄耗力量计算
- 强弱判断
- 喜用五行推导
-忌神五行推导
- 数据表和辅助函数

---

## 核心实现代码

### wuxing_xiji.go 实现

```go
package chronos

import (
    "strings"
    "github.com/godcong/fate/errors"
)

// 计算五行喜忌
func CalculateWuxingXiji(bazi *BaziInfo) (*WuxingXijiInfo, error) {
    // 1. 获取日主
    dayGan := bazi.Sizhu[2][0:1]
    dayWuxing := getGanWuxing(dayGan)
    
    // 2. 统计五行次数
    wuxingCount := countWuxing(bazi)
    
    // 3. 计算生扶力量
    shengfuPower := calculateShengfuPower(dayWuxing, wuxingCount)
    
    // 4. 计算克泄耗力量
    kexiehaoPower := calculateKexiehaoPower(dayWuxing, wuxingCount)
    
    // 5. 判断强弱
    strength := judgeStrength(shengfuPower, kexiehaoPower)
    
    // 6. 推导喜用五行
    xiWuxing := deriveXiWuxing(dayWuxing, strength)
    
    // 7. 推导忌神五行
    jiWuxing := deriveJiWuxing(dayWuxing, strength)
    
    // 8. 生成分析说明
    analysis := generateAnalysis(dayGan, dayWuxing, strength, xiWuxing, jiWuxing)
    
    // 9. 组装返回
    return &WuxingXijiInfo{
        DayGan:       dayGan,
        DayWuxing:    dayWuxing,
        XiWuxing:     xiWuxing,
        JiWuxing:     jiWuxing,
        Analysis:     analysis,
        SuggestWuxing: strings.Join(xiWuxing, ""),
    }, nil
}
```

---

## 辅助函数实现

### countWuxing() 实现

```go
// 统计五行次数
func countWuxing(bazi *BaziInfo) map[string]float64 {
    wuxingCount := make(map[string]float64)
    
    // 统计天干五行
    for i, zhu := range bazi.Sizhu {
        gan := zhu[0:1]
        ganWuxing := getGanWuxing(gan)
        
        // 月柱天干权重加成
        if i == 1 {
            wuxingCount[ganWuxing] += 1.5
        } else {
            wuxingCount[ganWuxing] += 1.0
        }
    }
    
    // 统计地支五行
    for i, zhu := range bazi.Sizhu {
        zhi := zhu[1:2]
        zhiWuxing := getZhiWuxing(zhi)
        
        // 月柱地支权重加成
        if i == 1 {
            wuxingCount[zhiWuxing] += 1.5
        } else {
            wuxingCount[zhiWuxing] += 1.0
        }
    }
    
    // 统计藏干五行（权重0.5）
    for _, canggan := range bazi.Canggan {
        for _, gan := range canggan {
            ganWuxing := getGanWuxing(gan)
            wuxingCount[ganWuxing] += 0.5
        }
    }
    
    return wuxingCount
}
```

---

### calculateShengfuPower() 实现

```go
// 计算生扶力量
func calculateShengfuPower(dayWuxing string, wuxingCount map[string]float64) float64 {
    power := 0.0
    
    // 同类五行
    power += wuxingCount[dayWuxing] * 1.0
    
    // 生我五行
    shengwoWuxing := getShengwoWuxing(dayWuxing)
    power += wuxingCount[shengwoWuxing] * 0.8
    
    return power
}
```

---

### calculateKexiehaoPower() 实现

```go
// 计算克泄耗力量
func calculateKexiehaoPower(dayWuxing string, wuxingCount map[string]float64) float64 {
    power := 0.0
    
    // 克我五行
    kewoWuxing := getKewoWuxing(dayWuxing)
    power += wuxingCount[kewoWuxing] * 1.0
    
    // 泄我五行
    xiewoWuxing := getXiewoWuxing(dayWuxing)
    power += wuxingCount[xiewoWuxing] * 0.8
    
    // 耗我五行
    haowoWuxing := getHaowoWuxing(dayWuxing)
    power += wuxingCount[haowoWuxing] * 0.7
    
    return power
}
```

---

### judgeStrength() 实现

```go
// 判断强弱
func judgeStrength(shengfuPower, kexiehaoPower float64) string {
    if shengfuPower + kexiehaoPower == 0 {
        return "中和"
    }
    
    ratio := shengfuPower / (shengfuPower + kexiehaoPower)
    
    if ratio > 0.6 {
        return "身强"
    } else if ratio < 0.4 {
        return "身弱"
    } else {
        return "中和"
    }
}
```

---

### deriveXiWuxing() 实现

```go
// 推导喜用五行
func deriveXiWuxing(dayWuxing string, strength string) []string {
    xiWuxing := []string{}
    
    switch strength {
    case "身强":
        // 身强喜克泄耗
        xiWuxing = append(xiWuxing, 
            getKewoWuxing(dayWuxing),   // 克我
            getXiewoWuxing(dayWuxing),  // 泄我
            getHaowoWuxing(dayWuxing),  // 耗我
        )
    case "身弱":
        // 身弱喜生扶
        xiWuxing = append(xiWuxing,
            dayWuxing,                  // 同类
            getShengwoWuxing(dayWuxing), // 生我
        )
    case "中和":
        // 中和喜用平衡
        xiWuxing = append(xiWuxing,
            dayWuxing,
            getShengwoWuxing(dayWuxing),
        )
    }
    
    return xiWuxing
}
```

---

### deriveJiWuxing() 实现

```go
// 推导忌神五行
func deriveJiWuxing(dayWuxing string, strength string) []string {
    jiWuxing := []string{}
    
    switch strength {
    case "身强":
        // 身强忌生扶
        jiWuxing = append(jiWuxing,
            dayWuxing,
            getShengwoWuxing(dayWuxing),
        )
    case "身弱":
        // 身弱忌克泄耗
        jiWuxing = append(jiWuxing,
            getKewoWuxing(dayWuxing),
            getXiewoWuxing(dayWuxing),
            getHaowoWuxing(dayWuxing),
        )
    case "中和":
        // 中和忌神较少
        jiWuxing = append(jiWuxing,
            getKewoWuxing(dayWuxing),
        )
    }
    
    return jiWuxing
}
```

---

### generateAnalysis() 实现

```go
// 生成分析说明
func generateAnalysis(dayGan, dayWuxing, strength string, xiWuxing, jiWuxing []string) string {
    var analysis string
    
    // 日主信息
    analysis += fmt.Sprintf("日主为%s（%s五行）。", dayGan, dayWuxing)
    
    // 强弱判断
    analysis += fmt.Sprintf("日主%s。", strength)
    
    // 喜忌说明
    if strength == "身强" {
        analysis += "身强喜克泄耗，"
        analysis += fmt.Sprintf("喜用五行：%s。", strings.Join(xiWuxing, "、"))
        analysis += fmt.Sprintf("忌神五行：%s。", strings.Join(jiWuxing, "、"))
    } else if strength == "身弱" {
        analysis += "身弱喜生扶，"
        analysis += fmt.Sprintf("喜用五行：%s。", strings.Join(xiWuxing, "、"))
        analysis += fmt.Sprintf("忌神五行：%s。", strings.Join(jiWuxing, "、"))
    } else {
        analysis += "日主中和，五行较为平衡。"
        analysis += fmt.Sprintf("建议补益：%s。", strings.Join(xiWuxing, "、"))
    }
    
    return analysis
}
```

---

## 数据表设计

### constants.go 数据表

```go
package chronos

// 天干五行对照表
var ganWuxingMap = map[string]string{
    "甲": "木", "乙": "木",
    "丙": "火", "丁": "火",
    "戊": "土", "己": "土",
    "庚": "金", "辛": "金",
    "壬": "水", "癸": "水",
}

// 地支五行对照表
var zhiWuxingMap = map[string]string{
    "子": "水", "丑": "土",
    "寅": "木", "卯": "木",
    "辰": "土", "巳": "火",
    "午": "火", "未": "土",
    "申": "金", "酉": "金",
    "戌": "土", "亥": "水",
}

// 五行生克关系表
var wuxingShengkeMap = map[string]map[string]string{
    "木": {"shengwo": "水", "kewo": "金", "xiewo": "火", "haowo": "土"},
    "火": {"shengwo": "木", "kewo": "水", "xiewo": "土", "haowo": "金"},
    "土": {"shengwo": "火", "kewo": "木", "xiewo": "金", "haowo": "水"},
    "金": {"shengwo": "土", "kewo": "火", "xiewo": "水", "haowo": "木"},
    "水": {"shengwo": "金", "kewo": "土", "xiewo": "木", "haowo": "火"},
}

// 获取天干五行
func getGanWuxing(gan string) string {
    return ganWuxingMap[gan]
}

// 获取地支五行
func getZhiWuxing(zhi string) string {
    return zhiWuxingMap[zhi]
}

// 获取生我五行
func getShengwoWuxing(dayWuxing string) string {
    return wuxingShengkeMap[dayWuxing]["shengwo"]
}

// 获取克我五行
func getKewoWuxing(dayWuxing string) string {
    return wuxingShengkeMap[dayWuxing]["kewo"]
}

// 获取泄我五行
func getXiewoWuxing(dayWuxing string) string {
    return wuxingShengkeMap[dayWuxing]["xiewo"]
}

// 获取耗我五行
func getHaowoWuxing(dayWuxing string) string {
    return wuxingShengkeMap[dayWuxing]["haowo"]
}
```

---

## 错误处理

### 默认值处理

```go
// 默认五行喜忌（降级处理）
func defaultWuxingXiji(bazi *BaziInfo) *WuxingXijiInfo {
    dayGan := bazi.Sizhu[2][0:1]
    dayWuxing := getGanWuxing(dayGan)
    
    xiWuxing := []string{dayWuxing, getShengwoWuxing(dayWuxing)}
    jiWuxing := []string{getKewoWuxing(dayWuxing)}
    
    return &WuxingXijiInfo{
        DayGan:       dayGan,
        DayWuxing:    dayWuxing,
        XiWuxing:     xiWuxing,
        JiWuxing:     jiWuxing,
        Analysis:     "无法分析五行喜忌，使用默认值",
        SuggestWuxing: strings.Join(xiWuxing, ""),
    }
}
```

---

## 实现要点

### 要点总结

| 要点 | 说明 | 重要性 |
|-----|------|-------|
| **五行统计** | 统计八字中五行出现次数 | 高 |
| **月令加成** | 月柱五行权重加成 | 中 |
| **藏干权重** | 藏干五行权重调整为0.5 | 中 |
| **强弱判断** | 根据生扶比例判断强弱 | 高 |
| **喜忌推导** | 根据强弱推导喜用忌神 | 高 |
| **分析生成** | 生成详细的分析说明 | 中 |
| **数据表** | 预加载常量数据表 | 中 |
| **错误处理** | 降级处理、默认值 | 中 |

---

## 总结

五行喜忌实现包括五行统计、生扶力量计算、克泄耗力量计算、强弱判断、喜用五行推导、忌神五行推导。使用数据表预加载常量数据，提高性能。错误处理使用降级处理，返回默认值。

**核心实现**：CalculateWuxingXiji()
**辅助函数**：countWuxing()、calculateShengfuPower()、calculateKexiehaoPower()等
**数据表**：ganWuxingMap、zhiWuxingMap、wuxingShengkeMap
**错误处理**：defaultWuxingXiji()（降级处理）