# chronos 模块介绍

## 模块定位

**chronos** 是 fate 项目的**核心计算模块**，负责：
- 公历/农历转换
- 八字计算（年柱、月柱、日柱、时柱）
- 五行喜忌分析（喜用五行、忌神五行）
- 为 naming 和 analysis 模块提供数据

---

## 核心职责

### 主要职责

| 职责 | 说明 | 重要性 |
|-----|------|-------|
| **八字计算** | 计算四柱干支、五行、纳音、十神、藏干、旬空 | 核心 |
| **五行喜忌分析** | 分析日主强弱、推导喜用五行、忌神五行 | 核心 |
| **数据提供** | 为 naming 和 analysis 模块提供完整数据 | 核心 |
| **lunar-go 桥接** | 适配 lunar-go API，转换数据格式 | 重要 |

### 不负责的职责

| 职责 | 说明 | 由谁负责 |
|-----|------|---------|
| 名字筛选 | 根据五行喜忌筛选名字 | naming 模块 |
| 名字评分 | 综合评分名字 | naming 模块 |
| 格式化输出 | 格式化八字解析输出 | analysis 模块 |

---

## 模块边界

### 输入边界

**输入数据**：FateInput

```go
type FateInput struct {
    BirthDate time.Time  // 出生日期（公历）
    Gender    int        // 性别（1男，0女）
    IsLunar   bool       // 是否农历日期
    Surname   string     // 姓氏（可选）
}
```

**输入来源**：cmd 模块（CLI、API、批量）

---

### 输出边界

**输出数据**：FateData

```go
type FateData struct {
    SolarDate  string            // 公历日期
    LunarDate  string            // 农历日期
    Gender     int               // 性别
    Bazi       *BaziInfo         // 八字信息
    WuxingXiji *WuxingXijiInfo   // 五行喜忌信息
    Dayun      *DayunInfo        // 大运信息（可选）
}
```

**输出去向**：naming 模块、analysis 模块

---

### 边界原则

1. **最小暴露**：只暴露 GetFateData() 核心接口
2. **数据封装**：内部计算逻辑不暴露
3. **桥接隔离**：lunar-go 依赖通过桥接层隔离

---

## 模块依赖

### 内部依赖

| 依赖模块 | 用途 | 说明 |
|---------|------|------|
| config | 配置管理 | 配置参数 |

### 外部依赖

| 依赖库 | 用途 | 说明 |
|-------|------|------|
| lunar-go | 农历计算 | 公历转农历、八字计算 |
| Go 标准库 | 时间处理 | time.Time |

---

## 核心接口

### GetFateData()

**接口定义**：

```go
func GetFateData(input *FateInput) (*FateData, error)
```

**职责**：一次性返回所有计算数据

**特点**：
- 一次性计算，避免多次调用
- 返回完整数据，包含八字、五行喜忌等
- 简化 naming 和 analysis 模块使用

---

## 模块文件结构

```
chronos/
├── chronos.go          # 核心入口（GetFateData）
├── bazi.go             # 八字计算（CalculateBazi）
├── wuxing_xiji.go      # 五行喜忌分析（CalculateWuxingXiji）
├── bridge.go           # lunar-go 桥接（BridgeLunar）
├── lunar_adapter.go    # lunar-go 适配器（适配 API）
├── types.go            # 类型定义（FateInput、FateData等）
├── constants.go        # 常量定义（天干、地支、五行等）
└── errors.go           # 错误定义（FateError）
```

---

## 模块特点

### 1. 专为 fate 优化

**优化点**：
- 简化输出：输出格式更适合起名工具（数组格式）
- 强类型：使用强类型枚举（Constellation、Zodiac、SolarTerm）
- 一次性计算：GetFateData() 一次性返回所有数据

---

### 2. 桥接 lunar-go

**桥接原因**：
- lunar-go 功能完善，但输出格式不适合起名工具
- lunar-go API 变化风险，需要隔离

**桥接设计**：
- 使用适配器模式适配 lunar-go API
- 转换数据格式为 fate 专用格式
- 隔离 lunar-go API 变化

---

### 3. 自己实现五行喜忌分析

**原因**：
- lunar-go 不提供五行喜忌分析功能
- 五行喜忌分析是起名工具的核心需求

**实现**：
- 自己实现五行喜忌分析算法
- 使用数据表辅助计算
- 提供详细的喜忌分析说明

---

## 与其他模块的关系

### 与 naming 模块的关系

**数据提供**：
- chronos 为 naming 提供 FateData
- naming 使用 WuxingXijiInfo 筛选名字

**依赖关系**：
- naming 依赖 chronos 的 FateData
- naming 不依赖 chronos 内部实现

---

### 与 analysis 模块的关系

**数据提供**：
- chronos 为 analysis 提供 FateData
- analysis 使用 BaziInfo 格式化输出

**依赖关系**：
- analysis 依赖 chronos 的 FateData
- analysis 不依赖 chronos 内部实现

---

### 与 config 模块的关系

**配置使用**：
- chronos 使用 config 的配置参数
- 配置参数影响计算行为

**依赖关系**：
- chronos 依赖 config 的配置参数
- chronos 不依赖 config 内部实现

---

## 模块开发要点

### 开发原则

1. **最小暴露**：只暴露核心接口，内部实现不暴露
2. **桥接隔离**：lunar-go 依赖通过桥接层隔离
3. **数据封装**：FateData 封装所有计算结果
4. **错误处理**：统一错误类型，明确错误返回

---

### 开发步骤

**步骤1：桥接 lunar-go**
- 实现 lunar_adapter.go
- 适配 lunar-go API
- 转换数据格式

**步骤2：实现八字计算**
- 实现 bazi.go
- 使用 lunar-go 获取八字数据
- 转换为 BaziInfo

**步骤3：实现五行喜忌分析**
- 实现 wuxing_xiji.go
- 自己实现五行喜忌分析算法
- 使用数据表辅助计算

**步骤4：实现核心接口**
- 实现 chronos.go
- 组装八字计算和五行喜忌分析
- 返回 FateData

---

## 模块测试要点

### 测试范围

| 测试类型 | 测试内容 | 说明 |
|---------|---------|------|
| 单元测试 | GetFateData() | 核心接口测试 |
| 单元测试 | CalculateBazi() | 八字计算测试 |
| 单元测试 | CalculateWuxingXiji() | 五行喜忌分析测试 |
| 单元测试 | BridgeLunar() | lunar-go 桥接测试 |
| 集成测试 | chronos + naming | 模块集成测试 |
| 集成测试 | chronos + analysis | 模块集成测试 |

---

### 测试数据

**测试用例**：
- 已知八字案例（验证正确性）
- 边界日期（1900年、2100年）
- 错误输入（无效日期、无效性别）

---

## 模块性能要点

### 性能优化

**优化点**：
- 缓存 lunar-go 的常用计算结果
- 预加载节气时刻数据
- 使用 goroutine 并发处理批量任务

---

### 性能目标

| 性能指标 | 目标值 | 说明 |
|---------|-------|------|
| 单次计算响应时间 | < 1秒 | GetFateData() 响应时间 |
| 批量计算（100个） | < 30秒 | 批量处理响应时间 |

---

## 总结

chronos 是 fate 的核心计算模块，负责八字计算、五行喜忌分析、数据提供。模块专为 fate 优化，使用桥接模式隔离 lunar-go，自己实现五行喜忌分析。核心接口 GetFateData() 一次性返回所有数据，简化 naming 和 analysis 模块使用。

**核心职责**：八字计算、五行喜忌分析、数据提供
**核心接口**：GetFateData(input) → FateData
**模块特点**：专为 fate 优化、桥接 lunar-go、自己实现五行喜忌分析