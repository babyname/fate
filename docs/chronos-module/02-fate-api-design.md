# FateAPI 接口设计

## 接口总览

**FateAPI** 是 chronos 模块的核心接口，一次性返回所有计算数据。

**接口定义**：

```go
package chronos

func GetFateData(input *FateInput) (*FateData, error)
```

---

## 设计原则

### 1. 一次性计算

**原则**：一次性计算所有数据，避免多次调用

**原因**：
- 八字计算和五行喜忌分析依赖相同的数据（lunar-go）
- 多次调用会增加性能开销
- 一次性计算简化 naming 和 analysis 模块使用

---

### 2. 数据封装

**原则**：使用 FateData 封装所有计算结果

**原因**：
- 数据结构清晰，易于理解
- 避免返回多个分散的数据
- 易于扩展（添加新数据）

---

### 3. 输入简化

**原则**：使用 FateInput 简化输入参数

**原因**：
- 输入参数清晰，易于理解
- 避免过多参数
- 易于扩展（添加新参数）

---

## FateInput 设计

### 数据结构

```go
type FateInput struct {
    BirthDate time.Time  // 出生日期（公历）
    Gender    int        // 性别（1男，0女）
    IsLunar   bool       // 是否农历日期
    Surname   string     // 姓氏（可选）
}
```

---

### 字段说明

| 字段 | 类型 | 必选 | 说明 |
|-----|------|------|------|
| BirthDate | time.Time | 必选 | 出生日期（公历日期） |
| Gender | int | 必选 | 性别（1男，0女） |
| IsLunar | bool | 必选 | 是否农历日期（true表示农历） |
| Surname | string | 可选 | 姓氏（用于起名时筛选名字） |

---

### 输入验证

**验证规则**：

| 验证项 | 规则 | 说明 |
|-----|------|------|
| BirthDate | 1900-2100年 | lunar-go 数据范围限制 |
| Gender | 0或1 | 性别值限制 |
| IsLunar | true或false | 农历标志限制 |

**验证示例**：

```go
func ValidateInput(input *FateInput) error {
    if input.BirthDate.Year() < 1900 || input.BirthDate.Year() > 2100 {
        return errors.New("日期范围错误：1900-2100年")
    }
    if input.Gender != 0 && input.Gender != 1 {
        return errors.New("性别值错误：0或1")
    }
    return nil
}
```

---

## FateData 设计

### 数据结构

```go
type FateData struct {
    SolarDate  string            // 公历日期 "2020年1月23日 11:31"
    LunarDate  string            // 农历日期 "2019年腊月二十九 午时"
    Gender     int               // 性别 1
    Bazi       *BaziInfo         // 八字信息
    WuxingXiji *WuxingXijiInfo   // 五行喜忌信息
    Dayun      *DayunInfo        // 大运信息（可选）
}
```

---

### 字段说明

| 字段 | 类型 | 必选 | 说明 |
|-----|------|------|------|
| SolarDate | string | 必选 | 公历日期字符串（格式化输出） |
| LunarDate | string | 必选 | 农历日期字符串（格式化输出） |
| Gender | int | 必选 | 性别（1男，0女） |
| Bazi | *BaziInfo | 必选 | 八字完整信息 |
| WuxingXiji | *WuxingXijiInfo | 必选 | 五行喜忌完整信息 |
| Dayun | *DayunInfo | 可选 | 大运信息（可选功能） |

---

## BaziInfo 设计

### 数据结构

```go
type BaziInfo struct {
    Sizhu      [4]string      // 四柱干支 ["甲子", "乙丑", "丙寅", "丁卯"]
    Wuxing     [4]string      // 五行 ["木水", "木土", "火木", "火木"]
    Nayin      [4]string      // 纳音 ["海中金", "海中金", "炉中火", "炉中火"]
    Shishen    [4]string      // 十神 ["比肩", "劫财", "食神", "伤官"]
    Canggan    [4][]string    // 藏干 [["癸"], ["己","辛","癸"], ...]
    Xunkong    [4]string      // 旬空 ["戌亥", "申酉", "午未", "辰巳"]
    Zodiac     string         // 生肖 "鼠"
    Constellation string      // 星座 "水瓶座"
}
```

---

### 字段说明

| 字段 | 类型 | 说明 |
|-----|------|------|
| Sizhu | [4]string | 四柱干支（年柱、月柱、日柱、时柱） |
| Wuxing | [4]string | 四柱五行（每个柱的五行） |
| Nayin | [4]string | 四柱纳音（每个柱的纳音） |
| Shishen | [4]string | 天干十神（相对于日主） |
| Canggan | [4][]string | 地支藏干（每个地支的藏干） |
| Xunkong | [4]string | 四柱旬空（每个柱的旬空） |
| Zodiac | string | 生肖（根据年柱） |
| Constellation | string | 星座（根据公历日期） |

---

### 设计特点

**使用数组格式**：
- Sizhu、Wuxing、Nayin、Shishen 使用 `[4]string` 数组格式
- Canggan 使用 `[4][]string` 二维数组格式
- 简化输出，易于理解

**数组索引**：
- 索引0：年柱
- 索引1：月柱
- 索引2：日柱
- 索引3：时柱

---

## WuxingXijiInfo 设计

### 数据结构

```go
type WuxingXijiInfo struct {
    DayGan       string        // 日主（天干） "丙"
    DayWuxing    string        // 日主五行 "火"
    XiWuxing     []string      // 喜用五行 ["木", "火"]
    JiWuxing     []string      //忌神五行 ["金", "水"]
    Analysis     string        // 分析说明 "日主身强，喜克泄耗..."
    SuggestWuxing string       // 建议补益五行 "木火"
}
```

---

### 字段说明

| 字段 | 类型 | 说明 |
|-----|------|------|
| DayGan | string | 日主天干（日柱天干） |
| DayWuxing | string | 日主五行（日干五行） |
| XiWuxing | []string | 喜用五行列表（适合补益的五行） |
| JiWuxing | []string |忌神五行列表（不适合的五行） |
| Analysis | string | 五行喜忌分析说明（详细分析） |
| SuggestWuxing | string | 建议补益五行（简化的建议） |

---

### 设计特点

**提供列表形式**：
- XiWuxing、JiWuxing 使用 `[]string` 列表形式
- 易于 naming 模块筛选名字

**提供分析说明**：
- Analysis 提供详细的分析说明
- 便于用户理解五行喜忌

**提供简化建议**：
- SuggestWuxing 提简化的建议（如"木火"）
- 易于理解和使用

---

## 接口使用示例

### 基本使用

```go
// 创建输入
input := &chronos.FateInput{
    BirthDate: time.Date(2020, 1, 23, 11, 31, 0, 0, time.UTC),
    Gender:    1,
    IsLunar:   false,
    Surname:   "张",
}

// 获取数据
fateData, err := chronos.GetFateData(input)
if err != nil {
    log.Fatal(err)
}

// 使用数据
fmt.Println("八字：", fateData.Bazi.Sizhu)
fmt.Println("五行喜用：", fateData.WuxingXiji.XiWuxing)
```

---

### 农历日期使用

```go
// 使用农历日期
input := &chronos.FateInput{
    BirthDate: time.Date(2020, 1, 23, 11, 31, 0, 0, time.UTC),
    Gender:    1,
    IsLunar:   true,  // 标记为农历日期
}

// GetFateData 会自动将农历转为公历
fateData, err := chronos.GetFateData(input)
```

---

### naming 模块使用

```go
// naming 使用 FateData 筛选名字
fateData, err := chronos.GetFateData(input)

// 根据五行喜忌筛选名字
xiWuxing := fateData.WuxingXiji.XiWuxing
names := naming.FilterNames(xiWuxing)
```

---

### analysis 模块使用

```go
// analysis 使用 FateData 格式化输出
fateData, err := chronos.GetFateData(input)

// 格式化八字输出
baziOutput := analysis.FormatBaziOutput(fateData.Bazi)
fmt.Println(baziOutput)
```

---

## 接口错误处理

### 错误类型

```go
type FateError struct {
    Code    int    // 错误代码
    Message string // 错误消息
    Module  string // 所属模块
}

// 错误代码定义
const (
    ErrCodeInputInvalid    = 1001  // 输入无效
    ErrCodeDateRange       = 1002  // 日期范围错误
    ErrCodeCalculateBazi   = 2001  // 八字计算错误
    ErrCodeCalculateWuxing = 2002  // 五行喜忌计算错误
)
```

---

### 错误返回示例

```go
// 输入验证失败
fateData, err := chronos.GetFateData(input)
if err != nil {
    if fateErr, ok := err.(*chronos.FateError); ok {
        fmt.Printf("错误代码：%d\n", fateErr.Code)
        fmt.Printf("错误消息：%s\n", fateErr.Message)
    }
}
```

---

## 接口扩展策略

### 添加新数据

**扩展方式**：在 FateData 中添加新字段

**示例**：

```go
// 添加命宫信息
type FateData struct {
    SolarDate  string
    LunarDate  string
    Gender     int
    Bazi       *BaziInfo
    WuxingXiji *WuxingXijiInfo
    MingGong   *MingGongInfo  // 新增：命宫信息
}
```

---

### 添加新参数

**扩展方式**：在 FateInput 中添加新字段

**示例**：

```go
// 添加八字流派参数
type FateInput struct {
    BirthDate time.Time
    Gender    int
    IsLunar   bool
    Surname   string
    School    string  // 新增：八字流派（可选）
}
```

---

## 总结

FateAPI 是 chronos 的核心接口，使用 GetFateData(input) → FateData 的设计，一次性返回所有数据。FateInput 简化输入参数，FateData 封装所有计算结果，BaziInfo 和 WuxingXijiInfo 提供详细的八字和五行喜忌信息。

**核心接口**：GetFateData(input) → FateData
**设计原则**：一次性计算、数据封装、输入简化
**数据结构**：FateInput、FateData、BaziInfo、WuxingXijiInfo
**扩展策略**：在 FateData 和 FateInput 中添加新字段