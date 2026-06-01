# fate 目录结构

## 目录总览

fate 项目目录结构如下：

```
fate/
├── cmd/                    # 入口层（CLI命令）
│   ├── console/            # 交互式起名
│   │   ├── main.go         # 入口
│   │   ├── init.go         # 初始化
│   │   └── name.go         # 起名命令
│   ├── character/          # 字表管理
│   │   └── main.go         # 入口
│   ├── dictctl/            # 字典工具
│   │   └── main.go         # 入口
│   ├── seeddb/             # 数据库种子数据
│   │   └── main.go         # 入口
│   └── inspectdb/          # 数据库检查
│       └── main.go         # 入口
│
├── config/                 # 配置层（配置管理）
│   ├── config.go           # 配置加载和管理
│   ├── database.go         # 数据库配置
│   ├── log.go              # 日志配置
│   └── config_test.go      # 配置测试
│
├── dict/                   # 字典/字表（公开，可独立使用）
│   ├── dict.go             # CharEntry, MergeEntries, ValidateEntries
│   ├── index.go            # DictIndex, QueryFilter, Build
│   ├── kangxi_stroke.go    # 康熙笔画修正表
│   └── dict_test.go        # 字典测试
│
├── log/                    # 日志工具（公开）
│   ├── log.go              # Logger 接口
│   ├── file.go             # 文件日志
│   └── wrap.go             # 日志包装器
│
├── ent/                    # ent 生成代码（内部自动可见）
│   └── schema/             # Schema 定义
│       ├── character.go    # 字符 Schema
│       ├── version.go      # 版本 Schema
│       └── wu_xing.go      # 五行 Schema
│
├── internal/               # 内部实现（外部不可导入）
│   ├── bazi/               # 八字命理计算
│   │   ├── bazi.go         # BaZi, NewBaZi()
│   │   ├── xiyong.go       # 喜用神计算
│   │   ├── nayin.go        # 纳音
│   │   └── zodiac.go       # 生肖
│   │
│   ├── wuxing/             # 五行生克计算
│   │   ├── wu_xing.go      # 五行生克关系
│   │   └── san_cai.go      # 三才
│   │
│   ├── wuge/               # 五格计算
│   │   ├── wuge.go         # WuGe, CalcWuGe()
│   │   ├── dayan.go        # 大衍之数
│   │   └── result.go       # 计算结果
│   │
│   ├── zhouyi/             # 周易起卦
│   │   ├── zhouyi.go       # QiGua(), GuaYao
│   │   └── yao.go          # 爻
│   │
│   ├── filter/             # 过滤器系统
│   │   ├── filter.go       # Filter 接口和实现
│   │   └── option.go       # FilterOption
│   │
│   ├── rating/             # 评分系统
│   │   └── rating.go       # Rater, NameRating
│   │
│   ├── naming/             # 名字推荐
│   │   ├── naming.go       # Interface, Naming
│   │   ├── name.go         # Name, NameBasic
│   │   ├── raters.go       # 评分器
│   │   └── stroke.go       # 笔画
│   │
│   ├── session/            # 会话管理
│   │   ├── session.go      # Session 接口和实现
│   │   ├── input.go        # Input, Output
│   │   ├── cache.go        # 缓存
│   │   ├── filter_cache.go # 过滤缓存
│   │   ├── filter_cache2.go# 过滤缓存2
│   │   └── list.go         # 列表
│   │
│   ├── repository/         # 数据访问层
│   │   ├── repository.go   # Repository, New()
│   │   ├── character.go    # 字符查询/缓存
│   │   ├── wuxing.go       # 五行查询
│   │   ├── cache.go        # 内部缓存
│   │   └── log.go          # 日志
│   │
│   ├── database/           # 数据库连接层
│   │   └── database.go     # Builder, Client()
│   │
│   ├── seeddb/             # 种子数据
│   │   ├── seed.go         # 种子数据入口
│   │   ├── builtin_seed.go # 内置种子
│   │   ├── export.go       # 数据导出
│   │   ├── import.go       # 数据导入
│   │   ├── transform.go    # 数据转换
│   │   └── report.go       # 数据报告
│   │
│   └── analysis/           # 输出格式化
│       ├── analysis.go     # NameResult, FateAnalysis
│       ├── builder.go      # 输出构建器
│       ├── helpers.go      # 辅助函数
│       ├── scoring.go      # 评分输出
│       ├── sancai_data.go  # 三才数据
│       ├── zhouyi.go       # 周易卦象输出
│       ├── simplified_traditional.go # 简繁体
│       └── types.go        # 类型定义
│
├── example/                # 示例代码
│   └── create_a_name/      # 起名示例
│       └── main.go         # 示例入口
│
├── docs/                   # 文档目录
│   ├── architecture/       # 架构设计文档
│   ├── chronos-module/     # chronos 模块文档
│   ├── naming-module/      # naming 模块文档
│   ├── analysis-module/    # analysis 模块文档
│   ├── config-module/      # config 模块文档
│   ├── overview/           # 项目概述文档
│   ├── implementation/     # 实施计划文档
│   └── reference/          # 参考资料
│
├── data/                   # 数据目录
│   └── gua.data            # 卦象数据
│
├── fate.go                 # 🔓 公开 API 入口
├── version.go              # 🔓 版本号
├── go.mod                  # Go模块定义
├── go.sum                  # Go依赖锁定
├── Makefile                # Make构建文件
├── README.md               # 项目说明
└── LICENSE                 # 许可证
```

---

## 目录职责说明

### 1. cmd 目录（入口层）

**职责**：
- 提供CLI命令入口
- 提供起名交互入口
- 协调各模块工作流程

**子目录**：

| 子目录 | 职责 | 文件 |
|-----|------|------|
| **console** | 交互式起名 | main.go（入口）、init.go（初始化）、name.go（起名命令） |
| **character** | 字表管理 | main.go（入口） |
| **dictctl** | 字典工具 | main.go（入口） |
| **seeddb** | 数据库种子数据 | main.go（入口） |
| **inspectdb** | 数据库检查 | main.go（入口） |

---

### 2. config 目录（配置层）

**职责**：
- 配置文件加载（YAML）
- 数据库配置管理
- 日志配置管理
- 默认配置提供

**文件**：

| 文件 | 职责 | 说明 |
|-----|------|------|
| **config.go** | 配置加载和管理 | LoadConfig()、GetConfig() |
| **database.go** | 数据库配置 | DatabaseConfig |
| **log.go** | 日志配置 | 日志配置参数 |
| **config_test.go** | 配置测试 | 配置单元测试 |

---

### 3. internal 目录（内部实现）

**职责**：
- 核心计算逻辑（八字、五行、五格、周易）
- 业务逻辑（过滤、评分、推荐、会话）
- 数据访问（repository、database）
- 输出格式化（analysis）

**子目录**：

| 子目录 | 职责 | 说明 |
|-----|------|------|
| **bazi** | 八字命理计算 | BaZi, XiYong, NaYin, Zodiac |
| **wuxing** | 五行生克计算 | WuXing, SanCai |
| **wuge** | 五格计算 | WuGe, DaYan |
| **zhouyi** | 周易起卦 | QiGua, GuaYao |
| **filter** | 过滤器系统 | Filter, FilterOption |
| **rating** | 评分系统 | Rater, NameRating |
| **naming** | 名字推荐 | Interface, Naming, Name |
| **session** | 会话管理 | Session, Input, Output |
| **repository** | 数据访问层 | Repository, 缓存 |
| **database** | 数据库连接层 | Builder, Client() |
| **seeddb** | 种子数据 | 导入/导出/转换 |
| **analysis** | 输出格式化 | NameResult, FateAnalysis |

---

### 4. dict 目录（字典/字表）

**职责**：
- 字条数据管理（CharEntry）
- 内存索引构建和查询（DictIndex）
- 康熙笔画修正表

**文件**：

| 文件 | 职责 | 说明 |
|-----|------|------|
| **dict.go** | 字条数据 | CharEntry, MergeEntries, ValidateEntries |
| **index.go** | 内存索引 | DictIndex, QueryFilter, Build |
| **kangxi_stroke.go** | 康熙笔画修正 | GetScienceStrokeCorrection |

---

### 5. log 目录（日志工具）

**职责**：
- 日志接口定义
- 文件日志实现
- 日志包装器

**文件**：

| 文件 | 职责 | 说明 |
|-----|------|------|
| **log.go** | Logger 接口 | 日志接口定义 |
| **file.go** | 文件日志 | 文件日志实现 |
| **wrap.go** | 日志包装器 | 日志包装 |

---

### 6. ent 目录（ORM 生成代码）

**职责**：
- ent schema 定义
- ent 自动生成的 CRUD 代码

**Schema 文件**：

| 文件 | 职责 | 说明 |
|-----|------|------|
| **character.go** | 字符 Schema | Character 表定义 |
| **version.go** | 版本 Schema | Version 表定义 |
| **wu_xing.go** | 五行 Schema | WuXing 表定义 |

---

### 7. docs 目录（文档目录）

**职责**：
- 项目文档管理
- 设计文档、实施文档、参考文档

**子目录**：

| 子目录 | 职责 |
|-----|------|
| **architecture** | 架构设计文档 |
| **chronos-module** | chronos 模块文档 |
| **naming-module** | naming 模块文档 |
| **analysis-module** | analysis 模块文档 |
| **config-module** | config 模块文档 |
| **overview** | 项目概述文档 |
| **implementation** | 实施计划文档 |
| **reference** | 参考资料 |

---

### 8. 根目录文件

**职责**：
- 项目公开 API 入口、版本号、配置

**文件**：

| 文件 | 职责 | 说明 |
|-----|------|------|
| **fate.go** | 公开 API 入口 | Fate, Session, NameResult |
| **version.go** | 版本号 | Version |
| **go.mod** | Go模块定义 | 模块名称、依赖 |
| **go.sum** | Go依赖锁定 | 依赖版本锁定 |
| **Makefile** | Make构建文件 | 构建规则 |
| **README.md** | 项目说明 | 项目介绍 |
| **LICENSE** | 许可证 | 许可证声明 |

---

## 文件命名规范

### Go 文件命名规范

**原则**：
- 文件名与职责对应
- 文件名简洁明确
- 文件名使用小写字母
- 文件名使用下划线分隔

**示例**：

| 文件名 | 职责 | 说明 |
|-----|------|------|
| chronos.go | chronos 入口 | 模块入口文件 |
| bazi.go | 八字计算 | 八字计算逻辑 |
| wuxing_xiji.go | 五行喜忌分析 | 五行喜忌分析逻辑 |
| types.go | 类型定义 | 类型定义 |
| constants.go | 常量定义 | 常量定义 |
| errors.go | 错误定义 | 错误定义 |

---

### 文档文件命名规范

**原则**：
- 文档名与内容对应
- 文档名使用编号排序
- 文档名简洁明确
- 文档名使用小写字母
- 文档名使用连字符分隔

**示例**：

| 文档名 | 内容 | 说明 |
|-----|------|------|
| 01-project-intro.md | 项目简介 | 编号排序 |
| 02-core-workflow.md | 核心工作流程 | 编号排序 |
| 03-user-stories.md | 用户场景和需求 | 编号排序 |

---

### 数据文件命名规范

**原则**：
- 数据名与内容对应
- 数据名使用小写字母
- 数据名使用点号分隔扩展名

**示例**：

| 数据名 | 内容 | 说明 |
|-----|------|------|
| characters.db | 汉字数据库 | SQLite数据库 |
| characters.json | 汉字数据 | JSON格式数据 |
| config.yaml | 配置文件 | YAML格式配置 |

---

## 目录结构的意义

### 易于理解

**意义**：
- 目录结构清晰，易于理解项目结构
- 目录职责明确，易于理解每个目录的作用
- 文件命名规范，易于理解每个文件的职责

---

### 易于开发

**意义**：
- 目录结构清晰，易于开发新功能
- 模块化目录，易于并行开发
- 文件命名规范，易于找到对应文件

---

### 易于维护

**意义**：
- 目录结构清晰，易于维护项目
- 模块化目录，易于修改和扩展
- 文件命名规范，易于维护文件

---

### 易于测试

**意义**：
- 测试目录独立，易于管理测试文件
- 测试文件与源文件对应，易于找到测试文件
- 测试命名规范，易于理解测试内容

---

### 易于部署

**意义**：
- 脚本目录独立，易于管理脚本
- 构建脚本、部署脚本明确，易于部署
- 配置文件独立，易于配置管理

---

## 目录结构扩展策略

### 添加新模块

**策略**：
- 在根目录添加新模块目录（如 logging/）
- 在 docs 目录添加新模块文档（如 logging-module/）
- 在 test 目录添加新模块测试（如 logging_test.go）

---

### 添加新功能

**策略**：
- 在对应模块目录添加新文件
- 在 docs 目录添加新功能文档
- 在 test 目录添加新功能测试

---

### 添加新文档

**策略**：
- 在对应 docs 子目录添加新文档
- 遵循编号排序规则
- 遵循命名规范

---

## 总结

fate 项目目录结构清晰，分为8个主要部分：cmd（入口层）、config（配置层）、internal（内部实现）、dict（字典/字表）、log（日志工具）、ent（ORM 生成代码）、docs（文档目录）、根目录文件。每个目录职责明确，文件命名规范，易于理解、易于开发、易于维护、易于测试。

**核心目录**：cmd, config, internal, dict, ent（核心模块）
**公开包**：根包(fate.go), config, dict, log（外部可导入）
**内部包**：internal/*（外部不可导入，包含 bazi, wuxing, wuge, zhouyi, filter, rating, naming, session, repository, database, seeddb, analysis）
**目录结构意义**：易于理解、易于开发、易于维护、易于测试