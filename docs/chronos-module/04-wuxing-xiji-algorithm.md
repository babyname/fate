# 五行喜忌分析算法

## 算法总览

五行喜忌分析是 chronos 模块的核心算法，用于推导适合补益的五行（喜用五行）和不适合的五行（忌神五行）。

**算法流程**：

```
确定日主 → 判断日主强弱 → 推导喜用五行 → 推导忌神五行 → 生成分析说明
```

---

## 五行喜忌原理

### 五行基本概念

**五行**：金、木、水、火、土

**五行相生**：
木生火、火生土、土生金、金生水、水生木

**五行相克**：
木克土、土克水、水克火、火克金、金克木

---

### 日主概念

**日主**：日柱天干，代表命主本人

**日主五行**：日柱天干的五行属性

| 天干 | 五行 | 阴阳 |
|-----|------|------|
| 甲 | 木 | 阳 |
| 乙 | 木 | 阴 |
|丙 | 火 | 阳 |
| 丁 | 火 | 阴 |
|戊 | 土 | 阳 |
| 己 | 土 | 阴 |
| 庚 | 金 | 阳 |
| 辛 | 金 | 阴 |
| 壬 | 水 | 阳 |
| 癸 | 水 | 阴 |

---

### 日主强弱概念

**身强（日主强）**：
- 日主五行在八字中得到生扶（同类、生我者多）
- 特点：喜克泄耗（喜异类五行）

**身弱（日主弱）**：
- 日主五行在八字中受到克制（克我、泄我、耗我者多）
- 特点：喜生扶（喜同类五行、生我五行）

---

## 日主强弱判断算法

### 算法原理

**判断依据**：
- 统计八字中各五行出现次数
- 统计生扶日主的五行次数（同类、生我）
- 统计克泄耗日主的五行次数（克我、泄我、耗我）
- 对比判断强弱

---

### 算法步骤

**步骤1：统计五行次数**

```go
// 统计八字中五行出现次数
func countWuxing(bazi *BaziInfo) map[string]int {
    wuxingCount := make(map[string]int)
    
    // 统计四柱天干五行
    for _, zhu := range bazi.Sizhu {
        gan := zhu[0:1]  // 天干
        ganWuxing := getGanWuxing(gan)
        wuxingCount[ganWuxing]++
    }
    
    // 统计四柱地支五行
    for _, zhu := range bazi.Sizhu {
        zhi := zhu[1:2]  // 地支
        zhiWuxing := getZhiWuxing(zhi)
        wuxingCount[zhiWuxing]++
    }
    
    // 统计藏干五行
    for _, canggan := range bazi.Canggan {
        for _, gan := range canggan {
            ganWuxing := getGanWuxing(gan)
            wuxingCount[ganWuxing] += 0.5  // 藏干权重较低
        }
    }
    
    return wuxingCount
}
```

---

**步骤2：计算生扶力量**

```go
// 计算生扶力量（同类 + 生我）
func calculateShengfuPower(dayWuxing string, wuxingCount map[string]int) float64 {
    power := 0.0
    
    // 同类五行（与日主五行相同）
    power += float64(wuxingCount[dayWuxing]) * 1.0
    
    // 生我五行（生日主的五行）
    shengwoWuxing := getShengwoWuxing(dayWuxing)  // 如日主火，生我木
    power += float64(wuxingCount[shengwoWuxing]) * 0.8
    
    return power
}
```

---

**步骤3：计算克泄耗力量**

```go
// 计算克泄耗力量（克我 + 泄我 + 耗我）
func calculateKexiehaoPower(dayWuxing string, wuxingCount map[string]int) float64 {
    power := 0.0
    
    // 克我五行（克日主的五行）
    kewoWuxing := getKewoWuxing(dayWuxing)  // 如日主火，克我水
    power += float64(wuxingCount[kewoWuxing]) * 1.0
    
    // 泄我五行（日主生出的五行）
    xiewoWuxing := getXiewoWuxing(dayWuxing)  // 如日主火，泄我土
    power += float64(wuxingCount[xiewoWuxing]) * 0.8
    
    // 耗我五行（日主克出的五行）
    haowoWuxing := getHaowoWuxing(dayWuxing)  // 如日主火，耗我金
    power += float64(wuxingCount[haowoWuxing]) * 0.7
    
    return power
}
```

---

**步骤4：判断强弱**

```go
// 判断日主强弱
func judgeStrength(shengfuPower, kexiehaoPower float64) string {
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

### 五行关系函数

**五行生克关系**：

```go
// 生我五行（生日主的五行）
func getShengwoWuxing(dayWuxing string) string {
    switch dayWuxing {
    case "木": return "水"  // 水生木
    case "火": return "木"  // 木生火
    case "土": return "火"  // 火生土
    case "金": return "土"  // 土生金
    case "水": return "金"  // 金生水
    }
    return ""
}

// 克我五行（克日主的五行）
func getKewoWuxing(dayWuxing string) string {
    switch dayWuxing {
    case "木": return "金"  // 金克木
    case "火": return "水"  // 水克火
    case "土": return "木"  // 木克土
    case "金": return "火"  // 火克金
    case "水": return "土"  // 土克水
    }
    return ""
}

// 泄我五行（日主生出的五行）
func getXiewoWuxing(dayWuxing string) string {
    switch dayWuxing {
    case "木": return "火"  // 木生火
    case "火": return "土"  // 火生土
    case "土": return "金"  // 土生金
    case "金": return "水"  // 金生水
    case "水": return "木"  // 水生木
    }
    return ""
}

// 耗我五行（日主克出的五行）
func getHaowoWuxing(dayWuxing string) string {
    switch dayWuxing {
    case "木": return "土"  // 木克土
    case "火": return "金"  // 火克金
    case "土": return "水"  // 土克水
    case "金": return "木"  // 金克木
    case "水": return "火"  // 水克火
    }
    return ""
}
```

---

## 喜用五行推导算法

### 算法原理

**身强喜克泄耗**：
- 身强时，喜异类五行（克我、泄我、耗我）
- 目的：削弱日主力量，达到平衡

**身弱喜生扶**：
- 身弱时，喜同类五行和生我五行
- 目的：增强日主力量，达到平衡

---

### 算法实现

```go
// 推导喜用五行
func deriveXiWuxing(dayWuxing string, strength string) []string {
    xiWuxing := []string{}
    
    if strength == "身强" {
        // 身强喜克泄耗
        kewo := getKewoWuxing(dayWuxing)    // 克我
        xiewo := getXiewoWuxing(dayWuxing)  // 泄我
        haowo := getHaowoWuxing(dayWuxing)  // 耗我
        
        xiWuxing = append(xiWuxing, kewo, xiewo, haowo)
        
        // 按重要性排序：克我 > 泄我 > 耗我
        xiWuxing = sortXiWuxing(xiWuxing, dayWuxing)
        
    } else if strength == "身弱" {
        // 身弱喜生扶
        tonglei := dayWuxing               // 同类
        shengwo := getShengwoWuxing(dayWuxing)  // 生我
        
        xiWuxing = append(xiWuxing, tonglei, shengwo)
        
        // 按重要性排序：同类 > 生我
        xiWuxing = sortXiWuxing(xiWuxing, dayWuxing)
        
    } else {
        // 中和，喜用五行平衡
        xiWuxing = append(xiWuxing, dayWuxing, getShengwoWuxing(dayWuxing))
    }
    
    return xiWuxing
}
```

---

##忌神五行推导算法

### 算法原理

**忌神五行**：与喜用五行相反的五行

**身强忌生扶**：
- 身强时，忌同类五行和生我五行
- 原因：会进一步增强日主力量，导致失衡

**身弱忌克泄耗**：
- 身弱时，忌克我、泄我、耗我五行
- 原因：会进一步削弱日主力量，导致失衡

---

### 算法实现

```go
// 推导忌神五行
func deriveJiWuxing(dayWuxing string, strength string) []string {
    jiWuxing := []string{}
    
    if strength == "身强" {
        // 身强忌生扶
        tonglei := dayWuxing               // 同类
        shengwo := getShengwoWuxing(dayWuxing)  // 生我
        
        jiWuxing = append(jiWuxing, tonglei, shengwo)
        
    } else if strength == "身弱" {
        // 身弱忌克泄耗
        kewo := getKewoWuxing(dayWuxing)    // 克我
        xiewo := getXiewoWuxing(dayWuxing)  // 泄我
        haowo := getHaowoWuxing(dayWuxing)  // 耗我
        
        jiWuxing = append(jiWuxing, kewo, xiewo, haowo)
        
    } else {
        // 中和，忌神五行较少
        jiWuxing = append(jiWuxing, getKewoWuxing(dayWuxing))
    }
    
    return jiWuxing
}
```

---

## 算法示例

### 示例1：日主丙火，身强

**八字**：甲子年、乙丑月、丙寅日、丁卯时

**日主**：丙（火）

**五行统计**：
- 木：甲、乙、寅、卯 → 4次
- 火：丙、丁 → 2次
- 土：丑 → 1次
- 金：子（藏癸）、丑（藏辛） → 约1次
- 水：子（藏癸）、丑（藏癸） → 约1次

**生扶力量**：
- 同类（火）：2次 × 1.0 = 2.0
- 生我（木）：4次 × 0.8 = 3.2
- 总计：5.2

**克泄耗力量**：
- 克我（水）：1次 × 1.0 = 1.0
- 泄我（土）：1次 × 0.8 = 0.8
- 耗我（金）：1次 × 0.7 = 0.7
- 总计：2.5

**强弱判断**：
- 生扶比例：5.2 / (5.2 + 2.5) = 0.67 > 0.6 → 身强

**喜用五行**：
- 身强喜克泄耗：水（克我）、土（泄我）、金（耗我）
- 排序：水、土、金

**忌神五行**：
- 身强忌生扶：火（同类）、木（生我）

---

### 示例2：日主甲木，身弱

**八字**：庚申年、辛酉月、甲子日、乙丑时

**日主**：甲（木）

**五行统计**：
- 木：甲、乙 → 2次
- 金：庚、辛、申、酉 → 4次
- 水：子、丑（藏癸） → 约2次
- 土：丑 → 1次
- 火：无 → 0次

**生扶力量**：
- 同类（木）：2次 × 1.0 = 2.0
- 生我（水）：2次 × 0.8 = 1.6
- 总计：3.6

**克泄耗力量**：
- 克我（金）：4次 × 1.0 = 4.0
- 泄我（火）：0次 × 0.8 = 0.0
- 耗我（土）：1次 × 0.7 = 0.7
- 总计：4.7

**强弱判断**：
- 生扶比例：3.6 / (3.6 + 4.7) = 0.43 < 0.4 → 身弱（接近）

**喜用五行**：
- 身弱喜生扶：木（同类）、水（生我）

**忌神五行**：
- 身弱忌克泄耗：金（克我）、火（泄我）、土（耗我）

---

## 算法优化要点

### 1. 藏干权重调整

**优化**：
- 藏干五行权重调整为0.5
- 避免藏干影响过大

---

### 2. 月令加成

**优化**：
- 月柱地支五行权重加成（×1.5）
- 月令对日主强弱影响较大

---

### 3. 调候考虑

**优化**：
- 考虑季节对五行的影响
- 如夏季火旺、冬季水旺

---

## 算法局限性

### 局限性1：简化判断

**说明**：
- 本算法使用简化判断（统计五行次数）
- 未考虑八字格局（如正格、从格）

**应对**：
- 后续版本可考虑更复杂的判断算法
- 或支持多流派算法

---

### 局限性2：权重固定

**说明**：
- 权重固定（同类1.0、生我0.8、克我1.0等）
- 未考虑不同八字的具体情况

**应对**：
- 后续版本可考虑动态权重调整

---

## 总结

五行喜忌分析算法包括日主强弱判断、喜用五行推导、忌神五行推导。算法原理基于五行生克关系和日主强弱判断，通过统计五行次数、计算生扶力量和克泄耗力量，判断日主强弱，然后推导喜用五行和忌神五行。

**算法流程**：确定日主 → 判断强弱 → 推导喜用 → 推导忌神
**核心算法**：五行统计、生扶力量计算、克泄耗力量计算、强弱判断
**算法优化**：藏干权重、月令加成、调候考虑