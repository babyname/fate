# fate 依赖管理

## 依赖总览

fate 的依赖分为两类：

| 依赖类型 | 数量 | 说明 |
|---------|------|------|
| **核心依赖** | 3个 | Go 运行时、chronos/v2、yi |
| **可选依赖** | 2个 | MySQL 驱动、YAML 解析器 |

> **注意**：Web API、诗词模块等商业功能已迁移至 qiming 项目，相关依赖（Web 框架等）已从 fate 中移除。

---

## go.mod 设计

### go.mod 内容

```go
module github.com/babyname/fate

go 1.25.0

require (
    entgo.io/ent v0.11.9
    github.com/babyname/chronos/v2 v2.0.7
    github.com/babyname/yi v1.0.2
    github.com/go-sql-driver/mysql v1.7.0
    github.com/spf13/cobra v1.6.1
    github.com/sqlite3ent/sqlite3 v1.50.0
    golang.org/x/exp v0.0.0-20221230185412-738e83a70c30
    golang.org/x/sync v0.20.0
    gopkg.in/yaml.v3 v3.0.1
)

replace (
    github.com/babyname/chronos/v2 => ../chronos
    github.com/babyname/yi => ../yi
)
```

---

### 模块依赖说明

| 模块 | 版本 | 用途 | 必选/可选 |
|-----|------|------|----------|
| **Go** | 1.25.0+ | Go 运行时 | 必选 |
| **chronos/v2** | v2.0.7+ | 八字计算、五行喜忌分析（含 lunar-go 桥接） | 必选 |
| **yi** | v1.0.2+ | 周易卦象数据 | 必选 |
| **ent** | v0.11.9+ | ORM 代码生成 | 必选 |
| **sqlite3ent/sqlite3** | v1.50.0+ | 纯 Go SQLite 驱动（CGO_ENABLED=0） | 必选 |
| **cobra** | v1.6.1+ | CLI 命令框架 | 必选 |
| **mysql** | v1.7.0+ | MySQL 驱动 | 可选 |
| **yaml.v3** | v3.0.1+ | YAML 解析器 | 必选 |

---

## chronos/v2 版本管理

### chronos/v2 简介

**项目地址**：`github.com/babyname/chronos/v2`（本地路径 `../chronos`）

**版本**：v2.0.7+

**功能**：
- 八字计算（年柱、月柱、日柱、时柱）
- 五行喜忌分析（喜用神、忌神）
- FateData 一次性数据返回
- Bridge 层桥接 lunar-go

**说明**：chronos/v2 封装了 lunar-go（`github.com/6tail/lunar-go`），fate 通过 chronos/v2 间接使用 lunar-go，不再直接依赖。

---

### yi 简介

**项目地址**：`github.com/babyname/yi`（本地路径 `../yi`）

**版本**：v1.0.2+

**功能**：
- 周易六十四卦数据
- 卦象查询（GuaXiang）
- 爻辞数据

---

### lunar-go 简介

**项目地址**：https://github.com/6tail/lunar-go

**版本**：v1.3.0+（chronos/v2 的间接依赖）

**功能**：
- 公历转农历、农历转公历
- 八字计算（年柱、月柱、日柱、时柱）
- 节气计算
- 星座、生肖计算
- 纳音、十神、藏干计算

---

### chronos/v2 版本锁定

**锁定策略**：
- go.mod 中通过 replace 指令指向本地 chronos 目录
- 避免意外升级导致 API 变化

**锁定示例**：

```go
// go.mod
replace (
    github.com/babyname/chronos/v2 => ../chronos
    github.com/babyname/yi => ../yi
)
```

**升级策略**：
- 测试新版本 API 是否兼容
- 更新 chronos 本地代码
- 更新 go.mod replace 指令

---

### lunar-go API 依赖

**依赖的 API**（通过 chronos/v2 间接使用）：

| API | 用途 | 说明 |
|-----|------|------|
| `calendar.Lunar` | 农历日期 | 公历转农历 |
| `calendar.Solar` | 公历日期 | 农历转公历 |
| `calendar.EightChar` | 八字计算 | 四柱干支 |
| `Lunar.GetJieQi()` | 节气计算 | 获取节气 |
| `Lunar.GetShengXiao()` |生肖 | 生肖 |
| `Lunar.GetConstellation()` | 星座 | 星座 |

---

### chronos/v2 桥接设计

**桥接模式**：

```
internal/bazi → chronos/v2 → lunar-go
```

**桥接层职责**：
- chronos/v2 封装 lunar-go API
- 提供 FateData 一次性数据返回
- 隔离 lunar-go API 变化

**桥接示例**：

```go
// chronos/v2 内部桥接
package chronos

import "github.com/6tail/lunar-go/calendar"

func GetFateData(input *FateInput) (*FateData, error) {
    solar := calendar.NewSolarFromDate(input.BirthDate)
    lunar := solar.GetLunar()
    eightChar := lunar.GetEightChar()
    // ... 计算 FateData
}
```

---

## 第三方依赖

### SQLite 依赖

**模块**：github.com/sqlite3ent/sqlite3

**版本**：v1.50.0+

**用途**：
- 汉字数据库存储
- 名字筛选数据查询
- 纯 Go 实现，支持 CGO_ENABLED=0

**依赖方式**：

```go
// go.mod
require (
    github.com/sqlite3ent/sqlite3 v1.50.0
)
```

**说明**：已从 `mattn/go-sqlite3`（CGO 依赖）迁移至 `sqlite3ent/sqlite3`（纯 Go），支持交叉编译。

---

### MySQL 依赖

**模块**：github.com/go-sql-driver/mysql

**版本**：v1.7.0+

**用途**：
- 可选的 MySQL 数据库驱动
- 多用户数据管理

**依赖方式**：

```go
// go.mod
require (
    github.com/go-sql-driver/mysql v1.7.0  // optional
)
```

---

### YAML 依赖

**模块**：gopkg.in/yaml.v3

**版本**：v3.0.1+

**用途**：
- 配置文件解析
- 参数配置管理

**依赖方式**：

```go
// go.mod
require (
    gopkg.in/yaml.v3 v3.0.1
)
```

**使用示例**：

```go
// config.go
package config

import "gopkg.in/yaml.v3"

func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

---

## 依赖风险和应对

### lunar-go API 变化风险

**风险**：
- lunar-go 版本升级可能导致 API 变化
- API 变化可能影响 chronos/v2 模块

**应对策略**：

1. **间接依赖**：
   - fate 不直接依赖 lunar-go，通过 chronos/v2 间接使用
   - lunar-go API 变化由 chronos/v2 屏蔽

2. **桥接隔离**：
   - chronos/v2 使用桥接模式隔离 lunar-go
   - API 变化时只需修改 chronos/v2

3. **测试验证**：
   - 升级前测试新版本 API 是否兼容
   - 验证八字计算结果是否正确

---

### SQLite 性能瓶颈风险

**风险**：
- SQLite 性能瓶颈可能影响批量处理
- 单机版性能限制

**应对策略**：

1. **优化索引**：
   - 为汉字数据库添加索引
   - 提高查询速度
   
2. **缓存优化**：
   - 缓存常用查询结果
   - 减少重复查询
   
3. **升级 MySQL**：
   - 需要时使用 MySQL
   - 提高性能

---

### Go 版本变化风险

**风险**：
- Go 版本升级可能导致语法变化
- Go 版本升级可能导致行为变化

**应对策略**：

1. **版本锁定**：
   - go.mod 指定 Go 1.25.0+
   - 避免意外升级
   
2. **测试验证**：
   - 升级前测试新版本是否兼容
   - 验证功能是否正常

---

## 依赖升级策略

### 升级流程

**升级流程**：

```
1. 测试新版本兼容性 → 2. 更新桥接层 → 3. 更新 go.mod → 4. 测试验证 → 5. 发布
```

---

### 升级示例（chronos/v2）

**升级步骤**：

```bash
# 1. 更新 chronos 本地代码
cd ../chronos
git pull

# 2. 测试 API 是否兼容
cd ../fate
go test ./internal/bazi/...

# 3. 更新 go.mod（如需要）
go mod tidy

# 4. 测试验证
go test ./...

# 5. 发布
git add go.mod go.sum
git commit -m "升级 chronos/v2"
git push
```

---

## 依赖管理最佳实践

### 最佳实践清单

| 最佳实践 | 说明 |
|---------|------|
| **版本锁定** | go.mod 锁定依赖版本 |
| **桥接隔离** | 使用桥接模式隔离第三方依赖 |
| **测试验证** | 升级前测试兼容性 |
| **最小依赖** | 只依赖必要的库 |
| **定期更新** | 定期检查依赖更新 |
| **文档记录** | 记录依赖版本和用途 |

---

### 依赖管理命令

**常用命令**：

```bash
# 初始化 go.mod
go mod init github.com/babyname/fate

# 添加依赖
go get github.com/babyname/chronos/v2@v2.0.7

# 更新依赖
go get -u github.com/babyname/chronos/v2

# 清理依赖
go mod tidy

# 验证依赖
go mod verify

# 查看依赖
go list -m all
```

---

## 总结

fate 依赖管理包括核心依赖（Go 1.25.0+、chronos/v2 v2.0.7+、yi v1.0.2+、ent v0.11.9+、sqlite3ent/sqlite3 v1.50.0+、cobra v1.6.1+、yaml.v3 v3.0.1+）和可选依赖（mysql v1.7.0+）。通过本地 replace 指令、桥接隔离、测试验证等策略，降低依赖风险，确保项目稳定运行。

**核心依赖**：Go 1.25.0+、chronos/v2 v2.0.7+、yi v1.0.2+、ent v0.11.9+、sqlite3ent/sqlite3 v1.50.0+、cobra v1.6.1+、yaml.v3 v3.0.1+
**可选依赖**：mysql v1.7.0+
**间接依赖**：lunar-go v1.3.0+（通过 chronos/v2）
**依赖风险**：lunar-go API变化、SQLite性能瓶颈、Go版本变化
**应对策略**：间接依赖隔离、桥接隔离、测试验证、优化索引、缓存优化