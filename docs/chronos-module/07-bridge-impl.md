# 桥接实现

## 实现总览

桥接实现包括 lunar_adapter.go 和 bridge.go，用于适配 lunar-go API。

---

## bridge.go 实现

```go
package chronos

// 桥接 lunar-go
func BridgeLunar(date time.Time, isLunar bool) (*LunarAdapter, error) {
    if isLunar {
        return FromLunar(date)
    }
    return FromSolar(date)
}
```

---

## lunar_adapter.go 实现

```go
package chronos

import "github.com/6tail/lunar-go/calendar"

type LunarAdapter struct {
    lunar *calendar.Lunar
}

func FromSolar(date time.Time) (*LunarAdapter, error) {
    solar := calendar.NewSolarFromDate(date)
    return &LunarAdapter{lunar: solar.GetLunar()}, nil
}

func FromLunar(date time.Time) (*LunarAdapter, error) {
    lunar := calendar.NewLunarFromDate(date)
    return &LunarAdapter{lunar: lunar}, nil
}

func (a *LunarAdapter) GetEightChar() *EightCharAdapter {
    return &EightCharAdapter{eightChar: a.lunar.GetEightChar()}
}

func (a *LunarAdapter) GetZodiac() string {
    return a.lunar.GetShengXiao()
}
```

---

## EightCharAdapter 实现

```go
type EightCharAdapter struct {
    eightChar *calendar.EightChar
}

func (a *EightCharAdapter) GetSizhu() [4]string {
    return [4]string{
        a.eightChar.GetYear().GetString(),
        a.eightChar.GetMonth().GetString(),
        a.eightChar.GetDay().GetString(),
        a.eightChar.GetTime().GetString(),
    }
}

func (a *EightCharAdapter) GetCanggan() [4][]string {
    return [4][]string{
        getCangganList(a.eightChar.GetYear().GetZhi()),
        getCangganList(a.eightChar.GetMonth().GetZhi()),
        getCangganList(a.eightChar.GetDay().GetZhi()),
        getCangganList(a.eightChar.GetTime().GetZhi()),
    }
}
```

---

## 辅助函数

```go
// 获取藏干列表
func getCangganList(zhi string) []string {
    cangganMap := map[string][]string{
        "子": {"癸"}, "丑": {"己", "辛", "癸"},
        "寅": {"甲", "丙", "戊"}, "卯": {"乙"},
        "辰": {"戊", "乙", "癸"}, "巳": {"丙", "戊", "庚"},
        "午": {"丁", "己"}, "未": {"己", "丁", "乙"},
        "申": {"庚", "壬", "戊"}, "酉": {"辛"},
        "戌": {"戊", "辛", "丁"}, "亥": {"壬", "甲"},
    }
    return cangganMap[zhi]
}
```

---

## 实现要点

| 要点 | 说明 |
|-----|------|
| 适配器模式 | 隔离 lunar-go API |
| 数据转换 | 对象 → 数组格式 |
| 错误处理 | 适配器捕获 lunar-go 错误 |

---

## 总结

桥接实现使用适配器模式，通过 LunarAdapter 和 EightCharAdapter 隔离 lunar-go，转换数据格式为数组格式，降低 API 变化风险。

**核心实现**：BridgeLunar()、LunarAdapter、EightCharAdapter
**辅助函数**：getCangganList()
**桥接优势**：隔离 API、转换格式、降低风险