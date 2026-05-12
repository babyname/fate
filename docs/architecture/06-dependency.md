# fate 依赖管理

## 依赖总览

fate 的依赖分为两类：

| 依赖类型 | 数量 | 说明 |
|---------|------|------|
| **核心依赖** | 2个 | Go 运行时、lunar-go |
| **可选依赖** | 3个 | SQLite、MySQL、YAML 解析器 |

---

## go.mod 设计

### go.mod 内容

```go
module github.com/godcong/fate

go 1.21

require (
    github.com/6tail/lunar-go v1.3.0
    github.com/mattn/go-sqlite3 v1.14.16
    gopkg.in/yaml.v3 v3.0.1
)

require (
    github.com/go-sql-driver/mysql v1.7.1  // optional
)
```

---

### 模块依赖说明

| 模块 | 版本 | 用途 | 必选/可选 |
|-----|------|------|----------|
| **Go** | 1.21+ | Go 运行时 | 必选 |
| **lunar-go** | v1.3.0+ | 农历计算库 | 必选 |
| **go-sqlite3** | v1.14.16+ | SQLite 驱动 | 必选（单机） |
| **mysql** | v1.7.1+ | MySQL 驱动 | 可选（企业） |
| **yaml.v3** | v3.0.1+ | YAML 解析器 | 必选 |

---

## lunar-go 版本管理

### lunar-go 简介

**项目地址**：https://github.com/6tail/lunar-go

**版本**：v1.3.0+

**功能**：
- 公历转农历、农历转公历
- 八字计算（年柱、月柱、日柱、时柱）
- 节气计算
- 星座、生肖计算
- 纳音、十神、藏干计算

---

### lunar-go 版本锁定

**锁定策略**：
- go.mod 中锁定 lunar-go 版本
- 避免意外升级导致 API 变化

**锁定示例**：

```go
// go.mod
require (
    github.com/6tail/lunar-go v1.3.0  // 锁定版本
)
```

**升级策略**：
- 测试新版本 API 是否兼容
- 更新桥接层适配新 API
- 升级版本并更新 go.mod

---

### lunar-go API 依赖

**依赖的 API**：

| API | 用途 | 说明 |
|-----|------|------|
| `calendar.Lunar` | 农历日期 | 公历转农历 |
| `calendar.Solar` | 公历日期 | 农历转公历 |
| `calendar.EightChar` | 八字计算 | 四柱干支 |
| `Lunar.GetJieQi()` | 节气计算 | 获取节气 |
| `Lunar.GetShengXiao()` |生肖 | 生肖 |
| `Lunar.GetConstellation()` | 星座 | 星座 |

---

### lunar-go 桥接设计

**桥接模式**：

```
chronos → BridgeLayer → lunar-go
```

**桥接层职责**：
- 适配 lunar-go API
- 转换数据格式
- 隔离 API 变化

**桥接示例**：

```go
// bridge.go
package chronos

import "github.com/6tail/lunar-go/calendar"

// 桥接 lunar-go
func BridgeLunar(date time.Time) (*calendar.Lunar, error) {
    solar := calendar.NewSolarFromDate(date)
    lunar := solar.GetLunar()
    return lunar, nil
}

// 适配 EightChar API
func GetEightChar(lunar *calendar.Lunar) *calendar.EightChar {
    return lunar.GetEightChar()
}
```

---

## 第三方依赖

### SQLite 依赖

**模块**：github.com/mattn/go-sqlite3

**版本**：v1.14.16+

**用途**：
- 汉字数据库存储
- 名字筛选数据查询

**依赖方式**：

```go
// go.mod
require (
    github.com/mattn/go-sqlite3 v1.14.16
)
```

**使用示例**：

```go
// database.go
package naming

import "github.com/mattn/go-sqlite3"

func QueryCharacters(wuxing string) ([]CharacterInfo, error) {
    db, err := sql.Open("sqlite3", "data/characters.db")
    if err != nil {
        return nil, err
    }
    // ...
}
```

---

### MySQL 依赖

**模块**：github.com/go-sql-driver/mysql

**版本**：v1.7.1+

**用途**：
- 企业版汉字数据库存储
- 多用户数据管理

**依赖方式**：

```go
// go.mod
require (
    github.com/go-sql-driver/mysql v1.7.1  // optional
)
```

**使用示例**：

```go
// database.go
package naming

import "github.com/go-sql-driver/mysql"

func QueryCharacters(wuxing string) ([]CharacterInfo, error) {
    db, err := sql.Open("mysql", "user:password@/database")
    if err != nil {
        return nil, err
    }
    // ...
}
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
- API 变化可能导致 chronos 模块失效

**应对策略**：

1. **版本锁定**：
   - go.mod 锁定 lunar-go 版本
   - 避免意外升级
   
2. **桥接隔离**：
   - chronos 使用桥接模式隔离 lunar-go
   - API 变化时只需修改桥接层
   
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
   - 企业版使用 MySQL
   - 提高性能

---

### Go 版本变化风险

**风险**：
- Go 版本升级可能导致语法变化
- Go 版本升级可能导致行为变化

**应对策略**：

1. **版本锁定**：
   - go.mod 指定 Go 1.21+
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

### 升级示例（lunar-go）

**升级步骤**：

```bash
# 1. 测试新版本兼容性
go get github.com/6tail/lunar-go@v1.4.0

# 2. 测试 API 是否兼容
go test ./chronos/...

# 3. 更新桥接层（如需要）
# 修改 bridge.go

# 4. 更新 go.mod
go mod tidy

# 5. 测试验证
go test ./...

# 6. 发布
git add go.mod go.sum
git commit -m "升级 lunar-go v1.4.0"
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
go mod init github.com/godcong/fate

# 添加依赖
go get github.com/6tail/lunar-go@v1.3.0

# 更新依赖
go get -u github.com/6tail/lunar-go

# 清理依赖
go mod tidy

# 验证依赖
go mod verify

# 查看依赖
go list -m all
```

---

## 总结

fate 依赖管理包括核心依赖（Go 1.21+、lunar-go v1.3.0+）和可选依赖（SQLite、MySQL、YAML解析器）。通过版本锁定、桥接隔离、测试验证等策略，降低依赖风险，确保项目稳定运行。

**核心依赖**：Go 1.21+、lunar-go v1.3.0+、go-sqlite3 v1.14.16+、yaml.v3 v3.0.1+
**可选依赖**：mysql v1.7.1+（企业版）
**依赖风险**：lunar-go API变化、SQLite性能瓶颈、Go版本变化
**应对策略**：版本锁定、桥接隔离、测试验证、优化索引、缓存优化、升级MySQL