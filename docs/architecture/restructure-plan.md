# fate 项目结构重组方案 v2

> 注意：此文档描述的功能中，transfer/（数据迁移）、cmd/fetchdata（数据抓取）已迁移至 qiming 项目。

> 在 v1 基础上增加分组分层和内外部访问控制

---

## 一、分组架构

项目分为 **4 个分组**，用 `internal/` 控制 external 不可见：

```
fate/
├── (根包)                      # 🔓 公开 API — 外部用户唯一入口
│
├── config/                     # 🔓 公开 — 配置结构体（外部可自定义配置）
│
├── internal/                   # 🔒 内部实现 — 外部不可导入
│   ├── bazi/                   #   八字命理计算
│   ├── wuxing/                 #   五行生克计算
│   ├── wuge/                   #   五格计算
│   ├── zhouyi/                 #   周易起卦
│   ├── filter/                 #   过滤器系统
│   ├── rating/                 #   评分系统
│   ├── naming/                 #   名字推荐
│   ├── session/                #   会话管理
│   ├── repository/             #   数据访问层
│   ├── database/               #   数据库连接
│   └── analysis/               #   输出格式化
│
├── dict/                       # 🔓 公开 — 字典/字表（可独立使用）
├── ent/                        # 🔒 ent 生成代码（internal 自动可见）
├── cmd/                        # CLI 工具
├── log/                        # 🔓 公开 — 日志工具
└── transfer/                   # 🔒 数据迁移（仅 cmd 使用）
```

### 访问控制规则

| 分组 | Go 可见性 | 谁可以使用 | 说明 |
|------|----------|-----------|------|
| 🔓 根包 | `package fate` | 所有人 | 外部用户唯一依赖 |
| 🔓 config | `package config` | 所有人 | 外部可自定义配置 |
| 🔓 dict | `package dict` | 所有人 | 字典库可独立使用 |
| 🔓 log | `package log` | 所有人 | 日志工具 |
| 🔒 internal/* | `internal` 机制 | 仅 fate 模块内部 | 实现细节，不暴露 |
| 🔒 ent/ | 不导出给外部 | 仅 internal 包 | ent 生成代码 |
| 🔒 transfer/ | 不导出给外部 | ~~仅 cmd/~~ 已迁移至 qiming |

---

## 二、完整目录结构

```
fate/
├── cmd/                            # CLI 工具
│   ├── console/                    #   交互式起名
│   ├── character/                  #   字表管理
│   ├── dictctl/                    #   字典工具
│   ├── inspectdb/                  #   数据库检查
│   └── seeddb/                     #   数据库种子数据
│
├── config/                         # 🔓 公开 — 配置
│   ├── config.go                   #   Config, LoadConfig()
│   ├── database.go                 #   DatabaseConfig
│   └── filter.go                   #   FilterConfig（从 martial.go 迁入）
│
├── dict/                           # 🔓 公开 — 字典/字表
│   ├── dict.go                     #   CharEntry, MergeEntries, ValidateEntries
│   ├── index.go                    #   DictIndex, QueryFilter, Build
│   └── kangxi_stroke.go            #   康熙笔画修正表 + GetScienceStrokeCorrection
│
├── log/                            # 🔓 公开 — 日志
│   ├── log.go                      #   Logger 接口
│   ├── file.go                     #   文件日志
│   └── wrap.go                     #   日志包装器
│
├── ent/                            # 🔒 ent 生成代码
│   └── schema/                     #   Schema 定义
│
├── internal/                       # 🔒 内部实现
│   ├── bazi/                       #   八字命理（纯计算）
│   │   ├── bazi.go                 #     BaZi, NewBaZi(), NewBaZiFromBridge()
│   │   ├── xiyong.go               #     XiYong, 喜用神计算
│   │   ├── nayin.go                #     NaYin, 纳音
│   │   ├── zodiac.go               #     Zodiac, GetZodiac()
│   │   └── table.go                #     天干地支强度表（内部数据）
│   │
│   ├── wuxing/                     #   五行计算（纯计算）
│   │   ├── wu_xing.go              #     五行生克关系
│   │   └── san_cai.go              #     SanCai, NewSanCai()
│   │
│   ├── wuge/                       #   五格计算（纯计算）
│   │   ├── wuge.go                 #     WuGe, CalcWuGe()
│   │   └── dayan.go                #     大衍之数（从 dayan/ 迁入）
│   │
│   ├── zhouyi/                     #   周易（纯计算）
│   │   └── zhouyi.go               #     QiGua(), GuaYao
│   │
│   ├── filter/                     #   过滤器系统
│   │   ├── filter.go               #     Filter 接口, filter 实现
│   │   └── option.go               #     FilterOption, CharacterFilterType
│   │
│   ├── rating/                     #   评分系统
│   │   ├── rater.go                #     Rater, NameRating, NewRater()
│   │   ├── wuxing.go               #     五行评分器
│   │   ├── bihua.go                #     笔画评分器
│   │   └── yinyun.go               #     音韵评分器
│   │
│   ├── naming/                     #   名字推荐
│   │   ├── naming.go               #     Interface, Naming, RecommendNames()
│   │   └── name.go                 #     Name, NameBasic, FirstName
│   │
│   ├── session/                    #   会话管理
│   │   ├── session.go              #     Session 接口, session 实现
│   │   ├── input.go                #     Input, Output
│   │   └── state.go                #     SessionState 枚举
│   │
│   ├── repository/                 #   数据访问层（替代 model/）
│   │   ├── repository.go           #     Repository, New()
│   │   ├── character.go            #     字符查询/缓存
│   │   ├── wugelucky.go            #     五格吉凶查询/缓存
│   │   └── cache.go                #     内部缓存结构
│   │
│   ├── database/                   #   数据库连接层
│   │   └── database.go             #     Builder, Client()
│   │
│   └── analysis/                   #   输出格式化
│       └── analysis.go             #     NameResult, FateAnalysis, Formatter
│
├── transfer/                       # 🔒 数据迁移（已迁移至 qiming，仅保留历史参考）
│   ├── config.go
│   └── transfer.go
│
├── fate.go                         # 🔓 公开 API 入口
├── version.go                      # 🔓 版本号
└── go.mod
```

---

## 三、每个包的导出控制

### 3.1 🔓 根包 `fate`（公开 API）

**原则**：外部用户只通过此包使用 fate。所有内部类型通过接口和 DTO 暴露。

| 导出项 | 类型 | 说明 |
|--------|------|------|
| `Fate` | interface | 主入口接口 |
| `New(cfg)` | func | 构造函数 |
| `Session` | interface | 会话接口（隐藏实现） |
| `NameResult` | struct | 推荐结果 DTO |
| `NameOption` | struct | 起名选项 DTO |
| `Version` | const | 版本号 |

**不导出**（内部实现细节）：
- `fateImpl` struct（小写，接口隐藏实现）
- 所有对 internal/ 包的引用

```go
package fate

type Fate interface {
    NewSession(opt *NameOption) Session
}

type Session interface {
    Start() error
    Stop() error
    Results() <-chan NameResult
}

type NameResult struct {
    FullName    string
    GivenName   string
    Strokes     string
    WuXing      string
    Score       float64
    Grade       string
    Interpret   string
    Source      string
}

type NameOption struct {
    LastName    [2]string
    BirthDate   time.Time
    Gender      int
    IsLunar     bool
    FilterType  string
}
```

### 3.2 🔓 `config`（公开）

| 导出 | 说明 |
|------|------|
| `Config` | 主配置结构体 |
| `DatabaseConfig` | 数据库配置 |
| `FilterConfig` | 过滤器配置（替代 Martial） |
| `LoadConfig(path)` | 加载配置 |
| `DefaultConfig()` | 默认配置 |

**不导出**：
- `loadFromYAML()` — 内部加载逻辑
- `validate()` — 内部验证

### 3.3 🔓 `dict`（公开，可独立使用）

| 导出 | 说明 |
|------|------|
| `CharEntry` | 字条数据结构 |
| `DictIndex` | 内存索引 |
| `QueryFilter` | 查询过滤条件 |
| `NewDictIndex()` | 创建索引 |
| `MergeEntries()` | 合并字条 |
| `ValidateEntries()` | 验证字条 |
| `GetScienceStrokeCorrection()` | 康熙笔画修正 |
| `ApplyKangxiCorrections()` | 批量修正 |

**不导出**：
- `kangxiStrokeMap` — 内部修正表数据
- `byChar`, `byWuXing` 等 — 内部索引字段

### 3.4 🔒 `internal/bazi`（内部）

| 导出（模块内可见） | 说明 |
|-------------------|------|
| `BaZi` | 八字结构体 |
| `NewBaZi()` | 构造函数 |
| `XiYong` | 喜用神 |
| `NaYin` | 纳音 |
| `Zodiac` | 生肖 |
| `GetZodiac()` | 获取生肖 |

**不导出**：
- `tiangan`, `dizhi` — 天干地支强度表
- `wuXingTianGan`, `wuXingDiZhi` — 五行映射表
- `calcXiYong()`, `calcSimilar()`, `calcHeterogeneous()` — 内部计算
- `diIndex`, `tianIndex` — 内部索引

### 3.5 🔒 `internal/wuxing`（内部）

| 导出（模块内可见） | 说明 |
|-------------------|------|
| `SanCai` | 三才结构体 |
| `NewSanCai()` | 构造函数 |
| `WuXingOf()` | 查五行属性 |
| `IsGenerating()` | 五行相生判断 |
| `IsOvercoming()` | 五行相克判断 |

**不导出**：
- `shengMap`, `keMap` — 生克映射表
- `sanCaiFortune` — 三才吉凶表

### 3.6 🔒 `internal/wuge`（内部）

| 导出（模块内可见） | 说明 |
|-------------------|------|
| `WuGe` | 五格结构体 |
| `CalcWuGe()` | 计算五格 |
| `DaYan()` | 大衍之数查表 |

**不导出**：
- `dayanData` — 81 个大衍数值
- `tianGe()`, `renGe()`, `diGe()` — 内部计算函数

### 3.7 🔒 `internal/filter`（内部）

| 导出（模块内可见） | 说明 |
|-------------------|------|
| `Filter` | 过滤器接口 |
| `FilterOption` | 过滤选项 |
| `NewFilter()` | 构造函数 |
| `DefaultFilter()` | 默认过滤器 |
| `CharacterFilterType` | 笔画类型枚举 |

**不导出**：
- `filter` struct — 实现细节
- `checkSkipStrokeNumberScope()` — 内部过滤方法
- `characterTypeFilter*()` — 内部过滤函数

### 3.8 🔒 `internal/rating`（内部）

| 导出（模块内可见） | 说明 |
|-------------------|------|
| `Rater` | 评分器 |
| `NameRating` | 评分结果 |
| `NewRater()` | 构造函数 |

**不导出**：
- `rateWuXing()`, `rateBiHua()`, `rateYinYun()` — 内部评分方法
- `scoreToGrade()` — 内部转换
- `wuxingWeight`, `bihuaWeight`, `yinyunWeight` — 内部权重

### 3.9 🔒 `internal/naming`（内部）

| 导出（模块内可见） | 说明 |
|-------------------|------|
| `Interface` | 推荐接口 |
| `Name` | 名字结构体 |
| `NameBasic` | 名字基本信息 |
| `FirstName` | 名字类型 |

**不导出**：
- `Naming` struct — 实现细节
- `commonLevel1` — 内部常用字表

### 3.10 🔒 `internal/session`（内部）

| 导出（模块内可见） | 说明 |
|-------------------|------|
| `Session` | 会话接口 |
| `Input` | 输入参数 |
| `Output` | 输出结果 |
| `SessionState` | 会话状态枚举 |

**不导出**：
- `session` struct — 实现细节
- `generate()` — 内部生成逻辑
- `startOutput()` — 内部输出逻辑

### 3.11 🔒 `internal/repository`（内部）

| 导出（模块内可见） | 说明 |
|-------------------|------|
| `Repository` | 数据访问层 |
| `New()` | 构造函数 |

**不导出**：
- `charCache`, `luckyCache` — 内部缓存
- `CharQuery`, `StrokeQuery` — 内部查询类型
- `getCharactersCached()` — 内部缓存查询

### 3.12 🔒 `internal/database`（内部）

| 导出（模块内可见） | 说明 |
|-------------------|------|
| `Builder` | 数据库构建器 |
| `New()` | 构造函数 |

**不导出**：
- `builder` struct — 实现细节

---

## 四、依赖关系图

```
                    ┌─────────────┐
                    │  fate(根包)  │  🔓 公开 API
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
        ┌─────▼─────┐ ┌───▼───┐ ┌──────▼──────┐
        │  config   │ │ dict  │ │   session   │  🔒 internal
        │    🔓     │ │  🔓   │ │             │
        └───────────┘ └───────┘ └──────┬──────┘
                                       │
                    ┌──────────────────┼──────────────────┐
                    │                  │                  │
              ┌─────▼─────┐    ┌───────▼───────┐   ┌─────▼─────┐
              │  filter   │    │    naming     │   │ repository│
              └─────┬─────┘    └───────┬───────┘   └─────┬─────┘
                    │                  │                  │
           ┌───────┼───────┐    ┌─────┼─────┐           │
           │       │       │    │           │           │
     ┌─────▼──┐ ┌──▼──┐ ┌─▼──┐ ┌▼────┐ ┌───▼───┐  ┌───▼───┐
     │ wuxing │ │wuge │ │bazi│ │rating│ │ bazi  │  │  ent  │
     └────────┘ └──┬──┘ └──┬─┘ └──┬───┘ └───┬───┘  └───────┘
                   │       │       │          │
                ┌──▼──┐ ┌──▼──┐ ┌──▼──┐  ┌───▼───┐
                │dayan│ │chron│ │wuxng│  │chronos│
                └─────┘ │/v2  │ └─────┘  │  /v2  │
                        └─────┘          └───────┘
```

**依赖规则**：
- 🔓 公开包（根包、config、dict、log）不依赖 🔒 internal 包
- 🔒 internal 包可依赖 🔓 公开包
- 🔒 internal 包之间有严格的单向依赖
- `repository` 是唯一依赖 `ent` 的包
- 纯计算包（bazi, wuxing, wuge, zhouyi）无数据库依赖

---

## 五、导出类型转换策略

根包需要将 internal 类型转换为公开 DTO，避免暴露内部实现：

```go
// fate.go
package fate

import (
    "github.com/babyname/fate/config"
    "github.com/babyname/fate/internal/session"
    "github.com/babyname/fate/internal/naming"
    "github.com/babyname/fate/internal/repository"
    "github.com/babyname/fate/internal/database"
)

type Fate interface {
    NewSession(opt *NameOption) Session
}

type fateImpl struct {
    cfg  *config.Config
    repo *repository.Repository
}

func New(cfg *config.Config) (Fate, error) {
    builder := database.New(cfg.Database)
    client, err := builder.Client()
    if err != nil {
        return nil, err
    }
    return &fateImpl{
        cfg:  cfg,
        repo: repository.New(client),
    }, nil
}

func (f *fateImpl) NewSession(opt *NameOption) Session {
    sess := session.New(f.repo, convertOption(opt))
    return &sessionAdapter{sess: sess}
}

// sessionAdapter 将 internal/session.Session 适配为公开 Session 接口
type sessionAdapter struct {
    sess *session.Session
}

func (a *sessionAdapter) Start() error {
    return a.sess.Start()
}

func (a *sessionAdapter) Results() <-chan NameResult {
    ch := make(chan NameResult, 128)
    go func() {
        defer close(ch)
        for {
            name, ok := a.sess.NextName()
            if !ok {
                return
            }
            ch <- NameResult{
                FullName:  name.String(),
                GivenName: name.GivenName(),
                WuXing:    name.WuXing(),
                Score:     name.Score(),
                Grade:     name.Grade(),
            }
        }
    }()
    return ch
}
```

---

## 六、internal/ 的 Go 机制

Go 的 `internal/` 目录机制：
- `internal/` 下的包只能被其**父目录树**内的代码导入
- 对于 `github.com/babyname/fate/internal/bazi`，只有 `fate/` 及其子目录可以导入
- 外部项目（如 `github.com/other/project`）**无法**导入这些包

这意味着：
- ✅ `fate/fate.go` 可以导入 `fate/internal/bazi`
- ✅ `fate/cmd/console/` 可以导入 `fate/internal/session`
- ❌ 外部项目无法导入 `fate/internal/bazi`
- ❌ 外部项目无法导入 `fate/internal/repository`

---

## 七、重组前后对比

| 指标 | 重组前 | 重组后 |
|------|--------|--------|
| 根包 .go 文件数 | 26 | 2 |
| 🔓 公开包 | 0（全部暴露） | 4（根包, config, dict, log） |
| 🔒 内部包 | 0（全部暴露） | 9（bazi, wuxing, wuge, zhouyi, filter, rating, naming, session, repository） |
| 外部可导入的包 | 全部 | 仅 4 个公开包 |
| 直接操作 ent 的包 | 4 | 1（repository） |
| 纯计算包（无DB依赖） | 3 | 5（bazi, wuxing, wuge, zhouyi, rating） |
| 根包导出类型 | 20+ | ~6（Fate, Session, NameResult, NameOption, Sex, Version） |

---

## 八、执行步骤

### Step 1: 删除废弃代码
- 删除 service/, statik/, fate 二进制, fate_db.zip
- 删除 filter_option_enum.go, session_enum.go

### Step 2: 创建 internal/ 目录结构
- 创建 internal/bazi/, internal/wuxing/, internal/wuge/, internal/zhouyi/
- 创建 internal/filter/, internal/rating/, internal/naming/, internal/session/
- 创建 internal/repository/, internal/database/, internal/analysis/

### Step 3: 迁移纯计算包（低风险）
- bazi.go + xiyong.go + nayin.go + zodiac.go → internal/bazi/
- wuxing/ → internal/wuxing/
- session_wuge.go + dayan/ → internal/wuge/
- yao.go + zhouyi.go → internal/zhouyi/

### Step 4: 迁移业务包（中风险）
- filter.go + filter_option.go → internal/filter/
- rating/ → internal/rating/（合并 naming/raters.go）
- naming/ → internal/naming/（合并 name.go + name_stroke.go）
- session.go + io.go → internal/session/

### Step 5: 重组数据层（中风险）
- model/ → internal/repository/
- cache/ → 合并到 internal/session/
- database/ → internal/database/

### Step 6: 精简根包（低风险）
- 根包仅保留 fate.go + version.go
- fate.go 实现 Fate/Session/NameResult 公开接口
- 删除根包 log.go

### Step 7: 更新所有 import 路径
- 全局替换 `github.com/babyname/fate/model` → `github.com/babyname/fate/internal/repository`
- 全局替换 `github.com/babyname/fate/cache` → `github.com/babyname/fate/internal/session`
- 等等

### Step 8: 验证
- `CGO_ENABLED=0 go build ./...`
- `go test ./...`
- 确认外部无法导入 internal/ 包
