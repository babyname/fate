# fate 接口设计

## 接口设计总览

fate 的接口设计遵循以下原则：
- **最小接口原则**：每个模块提供最小必要接口
- **接口隔离原则**：避免接口污染，接口职责单一
- **依赖倒置原则**：高层模块依赖抽象接口，不依赖具体实现
- **错误处理原则**：统一错误类型，明确错误返回

---

## 接口设计原则

### 1. 最小接口原则

**原则**：
每个模块只提供必要接口，避免过度暴露内部实现。

**示例**：

```go
// chronos 只提供一个核心接口
package chronos

// 核心接口：一次性返回所有数据
func GetFateData(input *FateInput) (*FateData, error)

// 不暴露内部接口：
// - CalculateBazi()（内部）
// - CalculateWuxingXiji()（内部）
// - BridgeLunar()（内部）
```

**意义**：
- 简化模块使用
- 避免接口复杂化
- 易于理解和维护

---

### 2. 接口隔离原则

**原则**：
每个接口职责单一，避免接口污染。

**示例**：

```go
// naming 提供三个职责单一的接口
package naming

// 职责单一：筛选名字
func FilterNames(xiWuxing []string) ([]NameInfo, error)

// 职责单一：评分名字
func RateNames(names []NameInfo, fateData *FateData) (*RatedNames, error)

// 聨责单一：推荐名字
func RecommendNames(fateData *FateData) (*RatedNames, error)

// 不设计一个臃肿接口：
// func ProcessNames(...)（包含筛选、评分、推荐）
```

**意义**：
- 职责清晰
- 易于理解和使用
- 易于测试和维护

---

### 3. 依赖倒置原则

**原则**：
高层模块依赖抽象接口，不依赖具体实现。

**示例**：

```go
// 定义抽象接口
type ChronosInterface interface {
    GetFateData(input *FateInput) (*FateData, error)
}

// naming 依赖抽象接口，不依赖具体实现
type Naming struct {
    chronos ChronosInterface
}

// 注入具体实现
func NewNaming(chronos ChronosInterface) *Naming {
    return &Naming{chronos: chronos}
}
```

**意义**：
- 易于替换实现
- 易于测试（使用 mock）
- 易于扩展

---

### 4. 错误处理原则

**原则**：
统一错误类型，明确错误返回。

**示例**：

```go
// 定义统一错误类型
package errors

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
    ErrCodeFilterNames     = 3001  // 名字筛选错误
)

// 统一错误返回
func GetFateData(input *FateInput) (*FateData, error) {
    if err := ValidateInput(input); err != nil {
        return nil, &FateError{
            Code:    ErrCodeInputInvalid,
            Message: err.Error(),
            Module:  "chronos",
        }
    }
    // ...
}
```

**意义**：
- 错误类型统一
- 错误信息明确
- 易于定位和处理

---

## 接口命名规范

### 函数命名规范

**规范**：
- 使用动词开头（Get、Calculate、Filter、Rate、Recommend）
- 使用驼峰命名（CamelCase）
- 简洁明确

**示例**：

| 函数名 | 职责 | 说明 |
|-----|------|------|
| GetFateData | 获取 fate 数据 | 核心接口 |
| CalculateBazi | 计算八字 | 内部接口 |
| CalculateWuxingXiji | 计算五行喜忌 | 内部接口 |
| FilterNames | 筛选名字 | naming 接口 |
| RateNames | 评分名字 | naming 接口 |
| RecommendNames | 推荐名字 | naming 接口 |
| FormatOutput | 格式化输出 | analysis 接口 |

---

### 参数命名规范

**规范**：
- 使用名词（Input、Data、Names）
- 使用驼峰命名（CamelCase）
- 简洁明确

**示例**：

| 参数名 | 类型 | 说明 |
|-----|------|------|
| input | *FateInput | fate 输入 |
| fateData | *FateData | fate 数据 |
| xiWuxing | []string | 喜用五行 |
| names | []NameInfo | 名字列表 |
| ratedNames | *RatedNames | 评分名字 |

---

### 返回值命名规范

**规范**：
- 主要返回值在前，错误在后
- 返回值类型明确

**示例**：

| 返回值 | 类型 | 说明 |
|-----|------|------|
| fateData | *FateData | fate 数据 |
| err | error | 错误 |
| ratedNames | *RatedNames | 评分名字 |
| output | *FormattedOutput | 格式化输出 |

---

## 模块接口定义

### chronos 模块接口

**核心接口**：

```go
package chronos

// 核心接口：一次性返回所有数据
func GetFateData(input *FateInput) (*FateData, error)

// 接口类型定义
type FateInput struct {
    BirthDate time.Time
    Gender    int
    IsLunar   bool
    Surname   string
}

type FateData struct {
    SolarDate  string
    LunarDate  string
    Gender     int
    Bazi       *BaziInfo
    WuxingXiji *WuxingXijiInfo
    Dayun      *DayunInfo
}
```

**内部接口**（不对外暴露）：

```go
// 内部接口：计算八字
func CalculateBazi(lunar *calendar.Lunar) (*BaziInfo, error)

// 内部接口：计算五行喜忌
func CalculateWuxingXiji(bazi *BaziInfo) (*WuxingXijiInfo, error)

// 内部接口：桥接 lunar-go
func BridgeLunar(date time.Time) (*calendar.Lunar, error)
```

**接口职责**：
- GetFateData：核心接口，一次性返回所有数据
- CalculateBazi：内部接口，计算八字信息
- CalculateWuxingXiji：内部接口，计算五行喜忌信息
- BridgeLunar：内部接口，桥接 lunar-go

---

### naming 模块接口

**核心接口**：

```go
package naming

// 核心接口：筛选名字
func FilterNames(xiWuxing []string) ([]NameInfo, error)

// 核心接口：评分名字
func RateNames(names []NameInfo, fateData *FateData) (*RatedNames, error)

// 核心接口：推荐名字（组合筛选和评分）
func RecommendNames(fateData *FateData) (*RatedNames, error)

// 接口类型定义
type NameInfo struct {
    Name     string
    Strokes  int
    Wuxing   string
    Score    float64
    Analysis string
}

type RatedNames struct {
    Names []NameInfo
}
```

**内部接口**（不对外暴露）：

```go
// 内部接口：查询汉字数据库
func QueryCharacters(wuxing string) ([]CharacterInfo, error)

// 内部接口：计算笔画
func CalculateStrokes(name string) int

// 内部接口：计算五行评分
func CalculateWuxingScore(name string, xiWuxing []string) float64
```

**接口职责**：
- FilterNames：筛选名字（根据五行喜忌）
- RateNames：评分名字（综合评分）
- RecommendNames：推荐名字（组合筛选和评分）
- QueryCharacters：内部接口，查询汉字数据库
- CalculateStrokes：内部接口，计算笔画
- CalculateWuxingScore：内部接口，计算五行评分

---

### analysis 模块接口

**核心接口**：

```go
package analysis

// 核心接口：格式化输出
func FormatOutput(fateData *FateData, ratedNames *RatedNames) (*FormattedOutput, error)

// 辅助接口：格式化八字输出
func FormatBaziOutput(bazi *BaziInfo) string

// 辅助接口：格式化名字输出
func FormatNameOutput(ratedNames *RatedNames) string

// 接口类型定义
type FormattedOutput struct {
    BaziOutput string
    NameOutput string
    Format     string
}

type OutputFormat string

const (
    OutputFormatText  OutputFormat = "text"
    OutputFormatJSON  OutputFormat = "json"
    OutputFormatHTML  OutputFormat = "html"
)
```

**内部接口**（不对外暴露）：

```go
// 内部接口：格式化器
func formatText(data interface{}) string

// 内部接口：模板管理
func loadTemplate(format OutputFormat) string
```

**接口职责**：
- FormatOutput：核心接口，格式化输出
- FormatBaziOutput：辅助接口，格式化八字输出
- FormatNameOutput：辅助接口，格式化名字输出
- formatText：内部接口，格式化器
- loadTemplate：内部接口，模板管理

---

### config 模块接口

**核心接口**：

```go
package config

// 核心接口：加载配置
func LoadConfig(path string) (*Config, error)

// 核心接口：验证配置
func ValidateConfig(config *Config) error

// 核心接口：获取配置
func GetConfig() *Config

// 接口类型定义
type Config struct {
    DatabasePath string
    OutputFormat string
    MaxNames     int
    ScoreWeights ScoreWeights
}

type ScoreWeights struct {
    Wuxing float64
    Strokes float64
    Tone    float64
}
```

**接口职责**：
- LoadConfig：核心接口，加载配置文件
- ValidateConfig：核心接口，验证配置参数
- GetConfig：核心接口，获取配置参数

---

## 接口依赖关系

### 模块接口依赖图

使用 Mermaid 绘制接口依赖图：

```mermaid
graph LR
    subgraph cmd接口
        Main[main()]
    end
    
    subgraph config接口
        LoadConfig[LoadConfig()]
        GetConfig[GetConfig()]
    end
    
    subgraph chronos接口
        GetFateData[GetFateData()]
    end
    
    subgraph naming接口
        FilterNames[FilterNames()]
        RateNames[RateNames()]
        RecommendNames[RecommendNames()]
    end
    
    subgraph analysis接口
        FormatOutput[FormatOutput()]
    end
    
    Main --> LoadConfig
    Main --> GetFateData
    Main --> RecommendNames
    Main --> FormatOutput
    
    RecommendNames --> FilterNames
    RecommendNames --> RateNames
    
    FormatOutput --> GetFateData
    FormatOutput --> RecommendNames
```

---

### 接口调用顺序

**单次起名接口调用顺序**：

```
1. config.LoadConfig()           → 加载配置
2. config.GetConfig()            → 获取配置
3. chronos.GetFateData()         → 计算 fate 数据
4. naming.RecommendNames()       → 推荐名字
   → naming.FilterNames()        → 筛选名字
   → naming.RateNames()          → 评分名字
5. analysis.FormatOutput()       → 格式化输出
6. 输出到用户
```

---

## 接口错误处理

### 错误类型定义

```go
package errors

// 统一错误类型
type FateError struct {
    Code    int
    Message string
    Module  string
}

// 实现 error 接口
func (e *FateError) Error() string {
    return fmt.Sprintf("[%s] Error %d: %s", e.Module, e.Code, e.Message)
}

// 错误代码定义
const (
    // 输入错误（1000-1099）
    ErrCodeInputInvalid    = 1001
    ErrCodeDateRange       = 1002
    ErrCodeGenderInvalid   = 1003
    
    // 计算错误（2000-2099）
    ErrCodeCalculateBazi   = 2001
    ErrCodeCalculateWuxing = 2002
    ErrCodeLunarConvert    = 2003
    
    // 推荐错误（3000-3099）
    ErrCodeFilterNames     = 3001
    ErrCodeRateNames       = 3002
    ErrCodeDatabaseQuery   = 3003
    
    // 输出错误（4000-4099）
    ErrCodeFormatOutput    = 4001
    ErrCodeTemplateLoad    = 4002
)
```

---

### 错误返回示例

```go
// chronos 错误返回
func GetFateData(input *FateInput) (*FateData, error) {
    // 输入验证
    if err := validateInput(input); err != nil {
        return nil, &FateError{
            Code:    ErrCodeInputInvalid,
            Message: "输入参数无效",
            Module:  "chronos",
        }
    }
    
    // 日期范围验证
    if input.BirthDate.Year() < 1900 {
        return nil, &FateError{
            Code:    ErrCodeDateRange,
            Message: "日期范围错误：1900-2100年",
            Module:  "chronos",
        }
    }
    
    // 八字计算
    bazi, err := calculateBazi(input.BirthDate)
    if err != nil {
        return nil, &FateError{
            Code:    ErrCodeCalculateBazi,
            Message: "八字计算错误",
            Module:  "chronos",
        }
    }
    
    // ...
}
```

---

### 错误处理策略

**错误处理策略**：

| 错误类型 | 处理策略 | 说明 |
|---------|---------|------|
| 输入错误 | 直接返回错误 | 输入无效，不继续处理 |
| 计算错误 | 降级处理或返回错误 | 计算失败，尝试降级或返回错误 |
| 推荐错误 | 扩大范围或返回错误 | 筛选失败，扩大范围或返回错误 |
| 输出错误 | 返回错误 | 格式化失败，返回错误 |

**降级处理示例**：

```go
// 五行喜忌分析失败时降级处理
func CalculateWuxingXiji(bazi *BaziInfo) (*WuxingXijiInfo, error) {
    wuxing, err := calculateWuxing(bazi)
    if err != nil {
        // 降级处理：返回默认值
        return DefaultWuxingXiji(), nil
    }
    return wuxing, nil
}

// 默认值
func DefaultWuxingXiji() *WuxingXijiInfo {
    return &WuxingXijiInfo{
        DayGan:       "甲",
        DayWuxing:    "木",
        XiWuxing:     []string{"水", "木"},
        JiWuxing:     []string{"金", "土"},
        Analysis:     "无法分析五行喜忌，使用默认值",
        SuggestWuxing: "水木",
    }
}
```

---

## 总结

fate 接口设计遵循最小接口原则、接口隔离原则、依赖倒置原则、错误处理原则。每个模块提供核心接口和内部接口，接口职责清晰，接口命名规范，接口依赖明确，错误处理统一。

**核心接口**：
- chronos：GetFateData()
- naming：FilterNames()、RateNames()、RecommendNames()
- analysis：FormatOutput()
- config：LoadConfig()、ValidateConfig()、GetConfig()

**接口设计原则**：最小接口、接口隔离、依赖倒置、错误处理
**接口命名规范**：动词开头、驼峰命名、简洁明确
**错误处理原则**：统一错误类型、明确错误返回、降级处理