# lunar-go 桥接设计

## 桥接设计总览

chronos 使用**适配器模式**桥接 lunar-go，隔离 API 变化，转换数据格式。

**桥接目的**：
1. 避免直接依赖 lunar-go API（降低风险）
2. 转换数据格式为 fate 专用格式（数组格式）
3. 适配 lunar-go API 变化（易于升级）

---

## 桥接架构

```
chronos → BridgeLayer → lunar-go
         ↑
    lunar_adapter.go
```

**BridgeLayer职责**：
- 公历转农历（FromSolar）
- 农历转公历（FromLunar）
- 获取八字对象（GetEightChar）
- 数据格式转换

---

## 适配器设计

### lunar_adapter.go 结构

```go
package chronos

import (
    "time"
    "github.com/6tail/lunar-go/calendar"
)

// lunar 适配器
type LunarAdapter struct {
    lunar *calendar.Lunar
}

// 公历转农历
func FromSolar(date time.Time) (*LunarAdapter, error) {
    solar := calendar.NewSolarFromDate(date)
    lunar := solar.GetLunar()
    return &LunarAdapter{lunar: lunar}, nil
}

// 农历转公历
func FromLunar(date time.Time) (*LunarAdapter, error) {
    lunar := calendar.NewLunarFromDate(date)
    return &LunarAdapter{lunar: lunar}, nil
}

// 获取八字
func (a *LunarAdapter) GetEightChar() *EightCharAdapter {
    eightChar := a.lunar.GetEightChar()
    return &EightCharAdapter{eightChar: eightChar}
}

// 获取农历日期字符串
func (a *LunarAdapter) GetDateString() string {
    return a.lunar.GetString()
}

// 获取生肖
func (a *LunarAdapter) GetZodiac() string {
    return a.lunar.GetShengXiao()
}
```

---

## 数据转换设计

### EightCharAdapter 设计

```go
// EightChar 适配器
type EightCharAdapter struct {
    eightChar *calendar.EightChar
}

// 获取四柱干支（数组格式）
func (a *EightCharAdapter) GetSizhu() [4]string {
    return [4]string{
        a.eightChar.GetYear().GetString(),
        a.eightChar.GetMonth().GetString(),
        a.eightChar.GetDay().GetString(),
        a.eightChar.GetTime().GetString(),
    }
}

// 获取五行（数组格式）
func (a *EightCharAdapter) GetWuxing() [4]string {
    return [4]string{
        getGanWuxing(a.eightChar.GetYear().GetGan()),
        getGanWuxing(a.eightChar.GetMonth().GetGan()),
        getGanWuxing(a.eightChar.GetDay().GetGan()),
        getGanWuxing(a.eightChar.GetTime().GetGan()),
    }
}

// 获取藏干（二维数组格式）
func (a *EightCharAdapter) GetCanggan() [4][]string {
    return [4][]string{
        getZhiCanggan(a.eightChar.GetYear().GetZhi()),
        getZhiCanggan(a.eightChar.GetMonth().GetZhi()),
        getZhiCanggan(a.eightChar.GetDay().GetZhi()),
        getZhiCanggan(a.eightChar.GetTime().GetZhi()),
    }
}
```

---

## API 适配设计

### lunar-go API 适配表

| lunar-go API | fate API | 说明 |
|-------------|---------|------|
| `calendar.NewSolarFromDate()` | `FromSolar()` | 公历转农历 |
| `calendar.NewLunarFromDate()` | `FromLunar()` | 农历转公历 |
| `lunar.GetEightChar()` | `adapter.GetEightChar()` | 获取八字 |
| `eightChar.GetYear()` | `adapter.GetSizhu()[0]` | 年柱（数组格式） |
| `eightChar.GetMonth()` | `adapter.GetSizhu()[1]` | 月柱（数组格式） |
| `eightChar.GetDay()` | `adapter.GetSizhu()[2]` | 日柱（数组格式） |
| `eightChar.GetTime()` | `adapter.GetSizhu()[3]` | 时柱（数组格式） |

---

## 桥接优势

### 1. 避免直接依赖

**优势**：
- chronos 不直接调用 lunar-go API
- 通过适配器隔离，降低风险

---

### 2. 转换数据格式

**优势**：
- lunar-go 返回对象格式
- fate 使用数组格式（更简洁）
- 适配器自动转换格式

---

### 3. 易于升级

**优势**：
- lunar-go API 变化时，只需修改适配器
- chronos 代码无需修改
- 降低升级风险

---

## 桥接要点

### 要点总结

| 要点 | 说明 | 重要性 |
|-----|------|-------|
| **适配器模式** | 使用适配器隔离 lunar-go | 高 |
| **数据转换** | 转换对象格式为数组格式 | 高 |
| **API 适配** | 适配 lunar-go API | 高 |
| **错误隔离** | lunar-go 错误通过适配器转换 | 中 |
| **性能优化** | 缓存 lunar-go 数据 | 中 |

---

## 总结

chronos 使用适配器模式桥接 lunar-go，通过 LunarAdapter 和 EightCharAdapter 隔离 lunar-go API，转换数据格式为 fate 专用格式（数组格式）。桥接设计避免直接依赖 lunar-go，降低 API 变化风险，易于升级和维护。

**桥接架构**：chronos → BridgeLayer → lunar-go
**适配器设计**：LunarAdapter、EightCharAdapter
**数据转换**：对象格式 → 数组格式
**桥接优势**：避免直接依赖、转换数据格式、易于升级