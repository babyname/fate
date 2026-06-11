<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go" alt="Go Version" />
  <img src="https://img.shields.io/badge/version-4.2.0-blue" alt="Version" />
  <img src="https://img.shields.io/badge/license-MIT-green" alt="License" />
  <img src="https://img.shields.io/badge/module-v4-7B42BC" alt="Module v4" />
</p>

<p align="center">
  <h1 align="center">🌌 fate · 天命之名</h1>
  <p align="center">
    <strong>基于八字五行 · 三才五格 · 周易卦象 · 诗词出处的智能起名引擎</strong><br>
    <em>FATE — Celestial Naming Engine</em>
  </p>
</p>

<p align="center">
  <a href="README_EN.md">English</a> | 中文
</p>

---

## ✨ 功能亮点

<table>
<tr>
<td width="50%">

### 🎯 八字推命
- 四柱八字精密推算（年/月/日/时）
- 五行强弱分析 + 调候用神
- **双算法喜用神**：平衡用神法 & 格局用神法（10种格局）
- 生肖关系、纳音五行

</td>
<td width="50%">

### 📊 多维评分
- **五维评分体系**：文化印象 · 五行八字 · 生肖 · 五格数理 · 音韵
- 三才配置吉凶分析
- 81数大衍 O(1) 查找
- 64卦周易卦象解读

</td>
</tr>
<tr>
<td width="50%">

### 🌐 Web 界面
- 现代 React + Tailwind CSS 响应式 UI
- 实时流式生成 + 轮询进度
- 候选名册（表格/卡片双视图）
- 探索模式：随机抽样 + 条件过滤

</td>
<td width="50%">

### ⚡ 高性能
- Ent ORM + SQLite 高效数据层
- ExcellentTable 流式 Top-N 数据结构
- 异步并发生成（session 模式）
- 嵌入式 Web 前端（单二进制分发）

</td>
</tr>
</table>

---

## 📸 界面预览

> 💡 **提示**：Web 界面运行在 `http://localhost:18080`，启动服务器后即可访问。

| 首页 | 命名结果 |
|------|----------|
| ![首页](docs/screenshots/home.png) | ![结果](docs/screenshots/results.png) |

| 候选名册（表格） | 候选名册（卡片） |
|------------------|-------------------|
| ![表格](docs/screenshots/table.png) | ![卡片](docs/screenshots/cards.png) |

---

## 🚀 快速开始

### 环境要求

- **Go** 1.22+
- **Bun**（前端开发需要，生产部署无需）

### 安装 & 运行

```bash
# 克隆项目
git clone https://github.com/babyname/fate.git
cd fate

# 下载依赖
go mod download

# 初始化数据库（必需）
go run ./cmd/dbinit

# 启动 Web 服务
go run ./cmd/server
# 访问 http://localhost:18080
```

### CLI 命令行

```bash
# 编译
go build -o fate ./cmd/console

# 生成名字
./fate name -s 张 -b "2024/06/15 10:30" -g boy
```

<details>
<summary>📋 CLI 输出示例</summary>

```
═══════════════════════════════════════
           Top 10 推荐名字
═══════════════════════════════════════
  #1  张适抡  评分:90.1  等级:上上  五行:火火
      天格:凶 人格:吉 地格:半吉 总格:吉
      三才:木水水(大吉)

  #2  张粮屠  评分:90.1  等级:上上  五行:火火
  #3  张适晶  评分:90.1  等级:上上  五行:火火
  ...
───────────────────────────────────────
  全部候选名字共 10000 个
═══════════════════════════════════════
```

**详细分析**（`name detail`）：

```
═══════════════════════════════════════
  张适抡  评分: 90.1  等级: 上上
═══════════════════════════════════════
  【五格分析】
  天格:12(凶)  人格:29(吉)  地格:30(半吉)
  外格:13(吉)  总格:41(吉)
  三才: 木水水(大吉)

  【八字分析】
  四柱: 甲辰 庚午 庚戌 乙酉
  喜神: 火  忌神: 金

  【周易卦象】
  本卦: 水山蹇  变卦: 地山谦

  【评分明细】
  文化:86  五行:94  生肖:94  五格:89  音韵:89
═══════════════════════════════════════
```
</details>

### CLI 参数

| 参数 | 简写 | 说明 | 示例 |
|------|------|------|------|
| `--surname` | `-s` | 姓氏 | `-s 张` |
| `--born` | `-b` | 出生日期与时刻 | `-b "2024/06/15 10:30"` |
| `--gender` | `-g` | 性别 | `-g boy` / `-g girl` |
| `--xiyong` | | 喜用神算法 | `balance`（平衡法）/ `geju`（格局法） |
| `--strictness` | | 五格严格度 | `loose` / `moderate` / `strict` |
| `--filter` | `-f` | 过滤汉字 | `-f 病死穷` |
| `--output` | `-o` | 输出到文件 | `-o result.txt` |

---

## 🔌 Web API

启动服务器后可用：

| Method | Endpoint | 说明 |
|--------|----------|------|
| `GET` | `/health` | 健康检查 |
| `POST` | `/api/generate` | 异步生成名字（返回 task_id） |
| `GET` | `/api/task/:id` | 查询任务状态 |
| `GET` | `/api/names/:taskId` | 分页获取名字列表 |
| `GET` | `/api/name-detail` | 名字详细分析 |
| `GET` | `/api/explore/:taskId` | 随机探索名字 |

### 生成请求示例

```bash
curl -X POST http://localhost:18080/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "surname": "张",
    "born": "2024/06/15 10:30",
    "sex": "boy",
    "poetry_mode": 2,
    "avoid_chars": ["病", "死"]
  }'
```

---

## 🧩 代码集成

```bash
go get github.com/babyname/fate/v4
```

```go
package main

import (
    "fmt"
    "time"

    "github.com/babyname/fate/v4"
    "github.com/babyname/fate/v4/config"
)

func main() {
    cfg := config.DefaultConfig()
    f, _ := fate.New(cfg)

    filter := fate.NewFilter(fate.FilterOption{
        CharacterFilter:     true,
        CharacterFilterType: fate.CharacterFilterTypeDefault,
        MinStroke:           3,
        MaxStroke:           18,
        RegularFilter:       true,
        DaYanFilter:         true,
        WuXingFilter:        true,
        AvoidCharacters:     []string{"病", "死", "穷"},
    })

    s := f.NewSessionWithFilter(filter)
    born, _ := time.Parse("2006/01/02 15:04", "2024/06/15 10:30")
    input := &fate.Input{
        Last: [2]string{"张", ""},
        Born: born,
        Sex:  fate.SexBoy,
    }

    s.Start(input)
    s.Wait()

    // 获取结果
    output := input.Output()
    for _, nr := range output.TopNames() {
        fmt.Printf("%s - %.1f (%s)\n", nr.FullName, nr.Score, nr.Grade)
    }

    // ExcellentTable: 流式 Top-N 访问
    table := output.GetExcellentTable()
    entries := table.TopN(10)
    for _, e := range entries {
        fmt.Printf("%s%s - %.1f has_poetry=%v\n", e.Char1, e.Char2, e.Score, e.HasPoetry)
    }

    // Explore: 随机采样 + 条件过滤
    explored := table.Explore(10, func(e fate.ExcellentEntry) bool {
        return e.HasPoetry
    })
    for _, e := range explored {
        fmt.Printf("%s%s - %.1f\n", e.Char1, e.Char2, e.Score)
    }
}
```

### 核心 API

| 类型 | 说明 |
|------|------|
| `fate.New(cfg)` | 创建引擎实例 |
| `fate.NewSessionWithFilter(filter)` | 创建带筛选器的会话 |
| `fate.NewFilter(option)` | 创建筛选器 |
| `fate.FilterOption` | 筛选选项（笔画/五行/大衍/性别/避用字等） |
| `fate.Input` | 输入（姓氏/生日/性别） |
| `fate.Output` | 输出（TopNames/ExcellentTable/CharMap） |
| `fate.ExcellentTable` | 流式 Top-N 结构（TryPush/Finalize/TopN/Explore） |

---

## 🏛️ 架构

```
┌─────────────────────────────────────────────┐
│                  fate CLI / Server           │
├─────────────────────────────────────────────┤
│  cmd/console     │  cmd/server               │
│  (CLI 工具)     │  (HTTP API + 嵌入式前端)   │
├─────────────────────────────────────────────┤
│               internal/http                  │
│               (API 路由 & 中间件)             │
├─────────────────────────────────────────────┤
│                  fate 核心引擎                │
│  ┌──────────┬──────────┬──────────────────┐ │
│  │ Session  │ Filter   │ ExcellentTable   │ │
│  │ 异步生成 │ 多级筛选  │ 流式Top-N        │ │
│  ├──────────┼──────────┼──────────────────┤ │
│  │ Bazi     │ Wuge     │ Zhouyi           │ │
│  │ 八字推算 │ 五格数理  │ 周易卦象          │ │
│  └──────────┴──────────┴──────────────────┘ │
├─────────────────────────────────────────────┤
│              chronos/v2 (八字引擎)            │
│              yi (周易引擎)                    │
│              Ent ORM (数据层)                 │
├─────────────────────────────────────────────┤
│              SQLite3 (数据存储)               │
│              character.json (汉字数据库)       │
│              chinese-poetry (诗词库)          │
└─────────────────────────────────────────────┘
```

### 目录结构

```
fate/
├── cmd/                # 命令行入口
│   ├── console/        #   主 CLI（name / dbinit / dictctl）
│   ├── server/         #   Web 服务
│   ├── dbinit/         #   数据库初始化
│   ├── character/      #   汉字数据工具
│   ├── dictctl/        #   字典管理工具
│   └── tools/          #   开发调试工具
├── config/             # 配置定义
├── data/               # 种子数据
│   ├── chinese-poetry/ #   诗词库
│   └── seed/           #   基础数据
├── docs/               # 文档 & 截图
├── ent/                # Ent Schema & 生成代码
├── internal/           # 内部包
│   ├── chronosfate/    #   chronos 适配层
│   ├── http/           #   HTTP API
│   └── ...
├── model/              # 数据模型
├── resources/          # 嵌入式资源
│   ├── character.json  #   汉字元数据库
│   └── static/         #   Web 前端构建产物
├── scripts/            # 构建脚本
└── web/                # React 前端源码
    └── src/
        ├── components/ # UI 组件
        ├── pages/      # 页面
        ├── lib/        # API 客户端 & 工具
        └── store/      # 状态管理 (Zustand)
```

---

## 📐 五维评分体系

| 维度 | 权重 | 评分依据 |
|------|:----:|----------|
| 🎨 **文化印象** | 20% | 常用度、正体字、字义丰富度、文化关联 |
| 🔥 **五行八字** | 25% | 名字五行与日主喜用神匹配度 |
| 🐉 **生肖** | 10% | 生肖五行与名字五行的生克关系 |
| 📏 **五格数理** | 25% | 天格/人格/地格/外格/总格 81 数吉凶 |
| 🎵 **音韵** | 20% | 姓名声韵调和谐度、平仄组合 |

---

## 🌿 喜用神算法

### 平衡用神法
基于日主强弱，通过同党/异党力量对比取用神：

| 日主 | 用神 | 喜神 | 忌神 |
|------|------|------|------|
| 强 | 克我（官杀） | 我克 + 生我 | 同我（比劫） |
| 弱 | 生我（印星） | 同我（比劫） | 克我 + 我生 |

### 格局用神法
先定格局（正官/七杀/食神/伤官/正财/偏财/建禄/阳刃等 10 种），再按格局取用神。

---

## 🔀 分支策略

| 分支 | 用途 | go.mod module |
|------|------|---------------|
| `main` | v4.x 开发主线 | `github.com/babyname/fate/v4` |
| `v3` | v3.x Bugfix 维护 | `github.com/babyname/fate` |

> ⚠️ **历史原因**：v3 分支的模块路径不含 `/v3`（早期 Go module 规范未强制要求）。
> 从 v4 开始引入 `/v4` 后缀，未来 v5 也将统一使用 `/v5`。

### 版本标签

| 系列 | 最新 | 标签 |
|------|------|------|
| v4.x | v4.2.0 | `v4.0.0` |
| v3.x | v3.6.0 | `v3.0.0` ~ `v3.6.0`（共 14 个发布） |
| v2.x | v2.0.2 | `v2.0.0` ~ `v2.0.2` |

---

## 🛠️ 技术栈

| 层 | 技术 | 说明 |
|----|------|------|
| 语言 | Go 1.22+ | 核心引擎 & 服务 |
| 前端 | React 19 + Tailwind CSS 3 + Rsbuild | Web 界面 |
| 状态 | Zustand | 前端状态管理 |
| UI | Radix UI + Lucide Icons | 无障碍组件库 |
| ORM | Ent | 类型安全的 Go ORM |
| 存储 | SQLite3 (modernc) | 嵌入式数据库，零依赖 |
| 八字 | chronos/v2 | 八字推算引擎 |
| 卦象 | yi | 周易六十四卦 |

---

## 🧪 开发

```bash
# 前端开发（热更新）
cd web && bun install && bun run dev

# 后端开发
go run ./cmd/server

# 构建生产版本
cd web && bun run build    # 前端构建 → resources/static/
go build -o fate ./cmd/console    # CLI 工具
go build -o fate-server ./cmd/server  # 服务器（含嵌入式前端）

# 运行测试
go test ./...

# 代码检查
golangci-lint run
```

---

## 📄 License

[MIT](LICENSE) © babyname
