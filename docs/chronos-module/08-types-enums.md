# 类型定义和枚举

## types.go 定义

```go
package chronos

// FateInput 输入参数
type FateInput struct {
    BirthDate time.Time
    Gender    int
    IsLunar   bool
    Surname   string
}

// FateData 输出数据
type FateData struct {
    SolarDate  string
    LunarDate  string
    Gender     int
    Bazi       *BaziInfo
    WuxingXiji *WuxingXijiInfo
}

// BaziInfo 八字信息
type BaziInfo struct {
    Sizhu      [4]string
    Wuxing     [4]string
    Nayin      [4]string
    Shishen    [4]string
    Canggan    [4][]string
    Xunkong    [4]string
    Zodiac     string
    Constellation string
}

// WuxingXijiInfo 五行喜忌信息
type WuxingXijiInfo struct {
    DayGan       string
    DayWuxing    string
    XiWuxing     []string
    JiWuxing     []string
    Analysis     string
    SuggestWuxing string
}
```

---

## constants.go 常量

```go
// 天干五行
var ganWuxingMap = map[string]string{
    "甲": "木", "乙": "木",
    "丙": "火", "丁": "火",
    "戊": "土", "己": "土",
    "庚": "金", "辛": "金",
    "壬": "水", "癸": "水",
}

// 地支五行
var zhiWuxingMap = map[string]string{
    "子": "水", "丑": "土",
    "寅": "木", "卯": "木",
    "辰": "土", "巳": "火",
    "午": "火", "未": "土",
    "申": "金", "酉": "金",
    "戌": "土", "亥": "水",
}
```

---

## 类型要点

| 类型 | 说明 |
|-----|------|
| FateInput | 输入参数封装 |
| FateData | 输出数据封装 |
| BaziInfo | 八字信息（数组格式） |
| WuxingXijiInfo | 五行喜忌信息 |

---

## 总结

类型定义包括 FateInput、FateData、BaziInfo、WuxingXijiInfo，使用数组格式简化输出。常量定义包括天干五行、地支五行对照表。

**核心类型**：FateInput、FateData、BaziInfo、WuxingXijiInfo
**核心常量**：ganWuxingMap、zhiWuxingMap