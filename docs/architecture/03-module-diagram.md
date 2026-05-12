# fate 模块图

## 模块依赖关系图

使用 Mermaid 绘制模块依赖关系图：

```mermaid
graph TD
    CMD[cmd] --> CONFIG[config]
    CMD --> CHRONOS[chronos]
    CMD --> NAMING[naming]
    CMD --> ANALYSIS[analysis]
    
    CHRONOS --> CONFIG
    CHRONOS --> LUNAR[lunar-go]
    
    NAMING --> CONFIG
    NAMING --> CHRONOS
    
    ANALYSIS --> CONFIG
    ANALYSIS --> CHRONOS
    ANALYSIS --> NAMING
    
    CONFIG --> NOONE[无依赖]
    LUNAR --> NOONE2[无依赖]
    
    style CMD fill:#f9f,stroke:#333,stroke-width:4px
    style CONFIG fill:#9ff,stroke:#333,stroke-width:2px
    style CHRONOS fill:#9f9,stroke:#333,stroke-width:2px
    style NAMING fill:#ff9,stroke:#333,stroke-width:2px
    style ANALYSIS fill:#f99,stroke:#333,stroke-width:2px
    style LUNAR fill:#ccc,stroke:#333,stroke-width:1px
```

---

## 模块层次关系图

使用 Mermaid 绘制模块层次关系图：

```mermaid
graph TB
    subgraph 入口层
        CMD[cmd]
    end
    
    subgraph 配置层
        CONFIG[config]
    end
    
    subgraph 数据层
        CHRONOS[chronos]
        LUNAR[lunar-go]
    end
    
    subgraph 业务层
        NAMING[naming]
    end
    
    subgraph 输出层
        ANALYSIS[analysis]
    end
    
    CMD --> CONFIG
    CMD --> CHRONOS
    CMD --> NAMING
    CMD --> ANALYSIS
    
    CHRONOS --> CONFIG
    CHRONOS --> LUNAR
    
    NAMING --> CONFIG
    NAMING --> CHRONOS
    
    ANALYSIS --> CONFIG
    ANALYSIS --> CHRONOS
    ANALYSIS --> NAMING
```

---

## 模块调用关系图

### 单次起名调用关系

使用 Mermaid 绘制单次起名调用关系图：

```mermaid
sequenceDiagram
    participant User as 用户
    participant CLI as cmd/cli
    participant Config as config
    participant Chronos as chronos
    participant Lunar as lunar-go
    participant Naming as naming
    participant Analysis as analysis
    
    User->>CLI: 输入日期、性别
    CLI->>Config: LoadConfig()
    Config-->>CLI: 配置参数
    
    CLI->>Chronos: GetFateData(input)
    Chronos->>Lunar: Lunar.FromSolar(date)
    Lunar-->>Chronos: LunarDate
    Chronos->>Lunar: Lunar.GetEightChar()
    Lunar-->>Chronos: EightChar
    Chronos->>Chronos: CalculateWuxingXiji()
    Chronos-->>CLI: FateData
    
    CLI->>Naming: FilterNames(FateData)
    Naming->>Naming: 查询汉字数据库
    Naming-->>CLI: FilteredNames
    
    CLI->>Naming: RateNames(FilteredNames)
    Naming-->>CLI: RatedNames
    
    CLI->>Analysis: FormatOutput(FateData, RatedNames)
    Analysis->>Analysis: FormatBaziOutput()
    Analysis->>Analysis: FormatNameOutput()
    Analysis-->>CLI: FormattedOutput
    
    CLI-->>User: 输出结果
```

---

### 批量起名调用关系

使用 Mermaid 绘制批量起名调用关系图：

```mermaid
sequenceDiagram
    participant User as 用户
    participant Batch as cmd/batch
    participant Config as config
    participant Chronos as chronos
    participant Naming as naming
    participant Analysis as analysis
    
    User->>Batch: 上传批量数据（Excel/CSV）
    Batch->>Config: LoadConfig()
    Config-->>Batch: 配置参数
    
    loop 每个请求
        Batch->>Chronos: GetFateData(input)
        Chronos-->>Batch: FateData
        
        Batch->>Naming: FilterNames(FateData)
        Naming-->>Batch: FilteredNames
        
        Batch->>Naming: RateNames(FilteredNames)
        Naming-->>Batch: RatedNames
        
        Batch->>Analysis: FormatOutput(FateData, RatedNames)
        Analysis-->>Batch: FormattedOutput
    end
    
    Batch-->>User: 批量输出结果
```

---

### API服务调用关系

使用 Mermaid 绘制 API 服务调用关系图：

```mermaid
sequenceDiagram
    participant Client as 客户端
    participant API as cmd/api
    participant Config as config
    participant Chronos as chronos
    participant Naming as naming
    participant Analysis as analysis
    
    Client->>API: HTTP POST /api/fate
    API->>Config: LoadConfig()
    Config-->>API: 配置参数
    
    API->>Chronos: GetFateData(input)
    Chronos-->>API: FateData
    
    API->>Naming: FilterNames(FateData)
    Naming-->>API: FilteredNames
    
    API->>Naming: RateNames(FilteredNames)
    Naming-->>API: RatedNames
    
    API->>Analysis: FormatOutput(FateData, RatedNames)
    Analysis-->>API: FormattedOutput
    
    API-->>Client: HTTP Response (JSON)
```

---

## 模块接口关系图

### chronos 接口关系

```mermaid
graph LR
    subgraph chronos
        FateAPI[GetFateData]
        BaziAPI[CalculateBazi]
        WuxingAPI[CalculateWuxingXiji]
        Bridge[BridgeLunar]
    end
    
    subgraph lunar-go
        Lunar[Lunar]
        Solar[Solar]
        EightChar[EightChar]
    end
    
    subgraph naming
        NamingInput[FateData]
    end
    
    subgraph analysis
        AnalysisInput[FateData]
    end
    
    naming --> FateAPI
    analysis --> FateAPI
    
    FateAPI --> BaziAPI
    FateAPI --> WuxingAPI
    
    BaziAPI --> Bridge
    Bridge --> Lunar
    Bridge --> Solar
    Bridge --> EightChar
    
    Lunar --> Solar
    EightChar --> Lunar
```

---

### naming 接口关系

```mermaid
graph LR
    subgraph naming
        FilterAPI[FilterNames]
        RateAPI[RateNames]
        RecommendAPI[RecommendNames]
        Database[汉字数据库]
    end
    
    subgraph chronos
        FateData[FateData]
    end
    
    subgraph analysis
        NameOutput[名字输出]
    end
    
    FilterAPI --> FateData
    RateAPI --> FilterAPI
    RecommendAPI --> RateAPI
    
    FilterAPI --> Database
    RateAPI --> Database
    
    NameOutput --> RecommendAPI
```

---

### analysis 接口关系

```mermaid
graph LR
    subgraph analysis
        FormatAPI[FormatOutput]
        BaziOutput[FormatBaziOutput]
        NameOutput[FormatNameOutput]
        Formatter[Formatter]
        Template[Template]
    end
    
    subgraph chronos
        FateData[FateData]
    end
    
    subgraph naming
        RatedNames[RatedNames]
    end
    
    subgraph output
        TextOutput[文本输出]
        JSONOutput[JSON输出]
        HTMLOutput[HTML输出]
    end
    
    FormatAPI --> FateData
    FormatAPI --> RatedNames
    
    FormatAPI --> BaziOutput
    FormatAPI --> NameOutput
    
    BaziOutput --> Formatter
    NameOutput --> Formatter
    
    Formatter --> Template
    
    Formatter --> TextOutput
    Formatter --> JSONOutput
    Formatter --> HTMLOutput
```

---

## 模块边界图

### 模块边界定义

使用 Mermaid 绘制模块边界图：

```mermaid
graph TB
    subgraph cmd边界
        CMD[cmd模块]
        CMD_IN[输入：用户参数]
        CMD_OUT[输出：用户结果]
    end
    
    subgraph config边界
        CONFIG[config模块]
        CONFIG_IN[输入：配置文件]
        CONFIG_OUT[输出：配置参数]
    end
    
    subgraph chronos边界
        CHRONOS[chronos模块]
        CHRONOS_IN[输入：日期、性别]
        CHRONOS_OUT[输出：FateData]
    end
    
    subgraph naming边界
        NAMING[naming模块]
        NAMING_IN[输入：FateData]
        NAMING_OUT[输出：RatedNames]
    end
    
    subgraph analysis边界
        ANALYSIS[analysis模块]
        ANALYSIS_IN[输入：FateData, RatedNames]
        ANALYSIS_OUT[输出：FormattedOutput]
    end
    
    CMD_IN --> CMD
    CMD --> CMD_OUT
    
    CONFIG_IN --> CONFIG
    CONFIG --> CONFIG_OUT
    
    CHRONOS_IN --> CHRONOS
    CHRONOS --> CHRONOS_OUT
    
    NAMING_IN --> NAMING
    NAMING --> NAMING_OUT
    
    ANALYSIS_IN --> ANALYSIS
    ANALYSIS --> ANALYSIS_OUT
    
    CMD --> CONFIG
    CMD --> CHRONOS
    CMD --> NAMING
    CMD --> ANALYSIS
    
    CHRONOS --> NAMING
    CHRONOS --> ANALYSIS
    
    NAMING --> ANALYSIS
```

---

## 模块数据流转图

### 数据流转路径

```mermaid
graph LR
    Input[用户输入] --> Config[配置参数]
    Input --> ChronosInput[chronos输入]
    
    Config --> Chronos
    Config --> Naming
    Config --> Analysis
    
    ChronosInput --> Chronos
    Chronos --> LunarData[lunar-go数据]
    LunarData --> Chronos
    
    Chronos --> FateData[FateData]
    
    FateData --> Naming
    Naming --> Database[汉字数据库]
    Database --> Naming
    
    Naming --> RatedNames[RatedNames]
    
    FateData --> Analysis
    RatedNames --> Analysis
    
    Analysis --> FormattedOutput[格式化输出]
    
    FormattedOutput --> Output[用户输出]
```

---

## 模块组件图

### chronos 模块组件

```mermaid
graph TB
    subgraph chronos
        FateAPI[FateAPI入口]
        BaziCalculator[八字计算器]
        WuxingCalculator[五行喜忌计算器]
        Bridge[桥接层]
        Adapter[适配器]
        Types[类型定义]
        Constants[常量定义]
        Errors[错误定义]
    end
    
    subgraph lunar-go
        Lunar[Lunar]
        Solar[Solar]
        EightChar[EightChar]
    end
    
    FateAPI --> BaziCalculator
    FateAPI --> WuxingCalculator
    
    BaziCalculator --> Bridge
    Bridge --> Adapter
    Adapter --> Lunar
    Adapter --> Solar
    Adapter --> EightChar
    
    WuxingCalculator --> Types
    WuxingCalculator --> Constants
    
    FateAPI --> Errors
    BaziCalculator --> Errors
    WuxingCalculator --> Errors
```

---

### naming 模块组件

```mermaid
graph TB
    subgraph naming
        NamingAPI[Naming入口]
        Filter[筛选器]
        Rating[评分器]
        Recommender[推荐器]
        Database[数据库]
        Strokes[笔画计算]
        Types[类型定义]
        Constants[常量定义]
        Errors[错误定义]
    end
    
    NamingAPI --> Filter
    NamingAPI --> Rating
    NamingAPI --> Recommender
    
    Filter --> Database
    Filter --> Strokes
    
    Rating --> Database
    Rating --> Strokes
    
    Recommender --> Filter
    Recommender --> Rating
    
    NamingAPI --> Errors
    Filter --> Errors
    Rating --> Errors
```

---

### analysis 模块组件

```mermaid
graph TB
    subgraph analysis
        AnalysisAPI[Analysis入口]
        BaziOutput[八字输出]
        NameOutput[名字输出]
        Formatter[格式化器]
        Template[模板]
        Types[类型定义]
        Constants[常量定义]
        Errors[错误定义]
    end
    
    AnalysisAPI --> BaziOutput
    AnalysisAPI --> NameOutput
    
    BaziOutput --> Formatter
    NameOutput --> Formatter
    
    Formatter --> Template
    
    AnalysisAPI --> Errors
    BaziOutput --> Errors
    NameOutput --> Errors
```

---

## 总结

fate 模块图包括：
- 模块依赖关系图（5个核心模块 + lunar-go依赖）
- 模块层次关系图（入口层、配置层、数据层、业务层、输出层）
- 模块调用关系图（单次起名、批量起名、API服务）
- 模块接口关系图（chronos、naming、analysis接口）
- 模块边界图（明确模块边界）
- 模块数据流转图（数据流转路径）
- 模块组件图（各模块内部组件）

**核心图**：模块依赖关系图、模块调用关系图、模块接口关系图
**辅助图**：模块层次关系图、模块边界图、模块数据流转图、模块组件图