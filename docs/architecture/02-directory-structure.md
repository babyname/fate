# fate 目录结构

## 目录总览

fate 项目目录结构如下：

```
fate/
├── cmd/                    # 入口层（CLI命令、API服务）
│   ├── fate-cli/           # CLI命令行工具
│   │   └── main.go         # CLI入口
│   ├── fate-api/           # HTTP API服务
│   │   └── main.go         # API入口
│   │   └── server.go       # API服务器
│   │   └── handlers.go     # API处理器
│   └── fate-batch/         # 批量处理工具
│   │   └── main.go         # 批量入口
│   │   └── batch.go        # 批量处理逻辑
│
├── config/                 # 配置层（配置管理）
│   ├── config.go           # 配置加载和管理
│   ├── validator.go        # 配置验证
│   ├── default.go          # 默认配置
│   └── types.go            # 配置类型定义
│
├── chronos/                # 计算层（八字、五行喜忌）
│   ├── chronos.go          # chronos 入口（FateAPI）
│   ├── bazi.go             # 八字计算
│   ├── wuxing_xiji.go      # 五行喜忌分析
│   ├── bridge.go           # lunar-go 桥接
│   ├── lunar_adapter.go    # lunar-go 适配器
│   ├── types.go            # 类型定义（BaziInfo, WuxingXijiInfo等）
│   ├── constants.go        # 常量定义（天干、地支等）
│   └── errors.go           # 错误定义
│
├── naming/                 # 推荐层（名字筛选、评分、推荐）
│   ├── naming.go           # naming 入口
│   ├── filter.go           # 名字筛选
│   ├── rating.go           # 名字评分
│   ├── recommend.go        # 名字推荐
│   ├── database.go         # 汉字数据库
│   ├── strokes.go          # 笔画计算
│   ├── types.go            # 类型定义（NameInfo, RatingInfo等）
│   ├── constants.go        # 常量定义（评分权重等）
│   └── errors.go           # 错误定义
│
├── analysis/               # 输出层（格式化输出）
│   ├── analysis.go         # analysis 入口
│   ├── bazi_output.go      # 八字输出
│   ├── name_output.go      # 名字输出
│   ├── formatter.go        # 格式化器
│   ├── template.go         # 模板管理
│   ├── types.go            # 类型定义（OutputFormat等）
│   ├── constants.go        # 常量定义（输出模板等）
│   └── errors.go           # 错误定义
│
├── docs/                   # 文档目录
│   ├── overview/           # 项目概述文档
│   ├── architecture/       # 架构设计文档
│   ├── chronos-module/     # chronos 模块文档
│   ├── naming-module/      # naming 模块文档
│   ├── analysis-module/    # analysis 模块文档
│   ├── config-module/      # config 模块文档
│   ├── implementation/     # 实施计划文档
│   └── reference/          # 参考资料
│
├── data/                   # 数据目录
│   ├── characters.db       # 汉字数据库（SQLite）
│   ├── characters.json     # 汉字数据（JSON格式）
│   └── config.yaml         # 配置文件
│
├── test/                   # 测试目录
│   ├── chronos_test.go     # chronos 测试
│   ├── naming_test.go      # naming 测试
│   ├── analysis_test.go    # analysis 测试
│   └── config_test.go      # config 测试
│
├── scripts/                # 脚本目录
│   ├── build.sh            # 构建脚本
│   ├── test.sh             # 测试脚本
│   ├── deploy.sh           # 部署脚本
│
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
- 提供HTTP API服务入口
- 提供批量处理入口
- 协调各模块工作流程

**子目录**：

| 子目录 | 职责 | 文件 |
|-----|------|------|
| **fate-cli** | CLI命令行工具 | main.go（CLI入口） |
| **fate-api** | HTTP API服务 | main.go（API入口）、server.go（服务器）、handlers.go（处理器） |
| **fate-batch** | 批量处理工具 | main.go（批量入口）、batch.go（批量逻辑） |

**文件说明**：

- **main.go**：入口文件，解析命令行参数、调用各模块
- **server.go**：API服务器，HTTP服务器配置
- **handlers.go**：API处理器，处理HTTP请求
- **batch.go**：批量逻辑，批量处理起名任务

---

### 2. config 目录（配置层）

**职责**：
- 配置文件加载（YAML/JSON）
- 配置参数验证
- 配置参数管理
- 默认配置提供

**文件**：

| 文件 | 职责 | 说明 |
|-----|------|------|
| **config.go** | 配置加载和管理 | LoadConfig()、GetConfig() |
| **validator.go** | 配置验证 | ValidateConfig() |
| **default.go** | 默认配置 | DefaultConfig() |
| **types.go** | 配置类型定义 | Config结构体 |

---

### 3. chronos 目录（计算层）

**职责**：
- 八字计算（年柱、月柱、日柱、时柱）
- 五行喜忌分析（喜用五行、忌神五行）
- 数据提供（为 naming 和 analysis 提供数据）
- lunar-go 桥接（适配器模式）

**文件**：

| 文件 | 职责 | 说明 |
|-----|------|------|
| **chronos.go** | chronos 入口（FateAPI） | GetFateData() |
| **bazi.go** | 八字计算 | CalculateBazi() |
| **wuxing_xiji.go** | 五行喜忌分析 | CalculateWuxingXiji() |
| **bridge.go** | lunar-go 桥接 | 桥接逻辑 |
| **lunar_adapter.go** | lunar-go 适配器 | 适配 lunar-go API |
| **types.go** | 类型定义 | BaziInfo, WuxingXijiInfo, FateData等 |
| **constants.go** | 常量定义 | 天干、地支、五行等 |
| **errors.go** | 错误定义 | 错误类型 |

---

### 4. naming 目录（推荐层）

**职责**：
- 名字筛选（根据五行喜忌筛选）
- 名字评分（综合评分：五行、笔画、音韵）
- 名字推荐（排序推荐最佳名字）
- 名字生成（组合规则）

**文件**：

| 文件 | 职责 | 说明 |
|-----|------|------|
| **naming.go** | naming 入口 | FilterNames()、RateNames()、RecommendNames() |
| **filter.go** | 名字筛选 | 筛选逻辑 |
| **rating.go** | 名字评分 | 评分逻辑 |
| **recommend.go** | 名字推荐 | 推荐逻辑 |
| **database.go** | 汉字数据库 | 数据库查询 |
| **strokes.go** | 笔画计算 | 笔画计算逻辑 |
| **types.go** | 类型定义 | NameInfo, RatingInfo等 |
| **constants.go** | 常量定义 | 评分权重等 |
| **errors.go** | 错误定义 | 错误类型 |

---

### 5. analysis 目录（输出层）

**职责**：
- 八字解析输出（格式化八字信息）
- 名字解析输出（格式化名字信息）
- 格式化（文本、JSON、HTML等格式）
- 模板管理（输出模板设计）

**文件**：

| 文件 | 职责 | 说明 |
|-----|------|------|
| **analysis.go** | analysis 入口 | FormatOutput() |
| **bazi_output.go** | 八字输出 | FormatBaziOutput() |
| **name_output.go** | 名字输出 | FormatNameOutput() |
| **formatter.go** | 格式化器 | 格式化逻辑 |
| **template.go** | 模板管理 | 模板定义和管理 |
| **types.go** | 类型定义 | OutputFormat等 |
| **constants.go** | 常量定义 | 输出模板等 |
| **errors.go** | 错误定义 | 错误类型 |

---

### 6. docs 目录（文档目录）

**职责**：
- 项目文档管理
- 设计文档、实施文档、参考文档

**子目录**：

| 子目录 | 职责 | 文件数量 |
|-----|------|---------|
| **overview** | 项目概述文档 | 6个 |
| **architecture** | 架构设计文档 | 6个 |
| **chronos-module** | chronos 模块文档 | 12个 |
| **naming-module** | naming 模块文档 | 12个 |
| **analysis-module** | analysis 模块文档 | 10个 |
| **config-module** | config 模块文档 | 8个 |
| **implementation** | 实施计划文档 | 10个 |
| **reference** | 参考资料 | 10个 |

**总计**：74个文档

---

### 7. data 目录（数据目录）

**职责**：
- 数据文件管理
- 配置文件、数据库文件

**文件**：

| 文件 | 职责 | 说明 |
|-----|------|------|
| **characters.db** | 汉字数据库 | SQLite数据库 |
| **characters.json** | 汉字数据 | JSON格式数据 |
| **config.yaml** | 配置文件 | YAML格式配置 |

---

### 8. test 目录（测试目录）

**职责**：
- 测试文件管理
- 单元测试、集成测试

**文件**：

| 文件 | 职责 | 说明 |
|-----|------|------|
| **chronos_test.go** | chronos 测试 | chronos 单元测试 |
| **naming_test.go** | naming 测试 | naming 单元测试 |
| **analysis_test.go** | analysis 测试 | analysis 单元测试 |
| **config_test.go** | config 测试 | config 单元测试 |

---

### 9. scripts 目录（脚本目录）

**职责**：
- 构建脚本、测试脚本、部署脚本

**文件**：

| 文件 | 职责 | 说明 |
|-----|------|------|
| **build.sh** | 构建脚本 | 构建可执行文件 |
| **test.sh** | 测试脚本 | 运行测试 |
| **deploy.sh** | 部署脚本 | 部署应用 |

---

### 10. 根目录文件

**职责**：
- 项目配置、说明、许可证

**文件**：

| 文件 | 职责 | 说明 |
|-----|------|------|
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

fate 项目目录结构清晰，分为10个主要目录：cmd（入口层）、config（配置层）、chronos（计算层）、naming（推荐层）、analysis（输出层）、docs（文档目录）、data（数据目录）、test（测试目录）、scripts（脚本目录）、根目录文件。每个目录职责明确，文件命名规范，易于理解、易于开发、易于维护、易于测试、易于部署。

**核心目录**：cmd, config, chronos, naming, analysis（5个模块目录）
**辅助目录**：docs, data, test, scripts, 根目录文件
**目录结构意义**：易于理解、易于开发、易于维护、易于测试、易于部署