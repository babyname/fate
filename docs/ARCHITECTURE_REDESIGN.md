# Architecture Redesign

**Date**: 2026-06-11

## 问题诊断

### 1. 繁简诗词匹配问题

当前 `import-poetry` 的 `extractCharRefs` 直接从诗词原文提取汉字：

```go
if !unicode.Is(unicode.Han, r) { continue }
refs = append(refs, &charRef{Char: string(r), ...})
```

问题：
- 诗句 "國破山河在" 中，`國` 被索引为 `國`（繁体）
- 用户搜 `国`（简体） → 查 `PoemChar.Char = '国'` → 查不到
- 近 40% 的 chinese-poetry 数据是繁体文本

### 2. fate 架构混乱

```
internal/  ← 什么都有：
├── analysis/      算法？结果组装？
├── bazi/          算法
├── chronosfate/   算法
├── naming/        算法（但依赖 repository）
├── rating/        算法
├── wuge/          算法
├── wuxing/        算法
├── zhouyi/        算法
├── repository/    数据库访问
├── seeddb/        数据导入
├── http/          Web
├── dict/          工具
├── filter/        工具
├── session/       业务逻辑
└── log/           工具
```

**cmd/ 同样混乱**：
- `server` / `console` 等命令混在一起，无清晰分层
- 命令无区分，无构建标签

---

## 架构设计方案

### 核心原则：四层分离

```
┌──────────────────────────────────────────────┐
│                  Server 层                    │  HTTP/RPC 入口、路由、中间件
├──────────────────────────────────────────────┤
│                 Service 层                    │  业务流程编排、session、analysis
├──────────────────────────────────────────────┤
│  Engine 层 (纯算法)     │   Data 层 (数据访问) │  算法不碰 DB · 数据不碰算法
├──────────────────────────────────────────────┤
│                      Model 层                 │  ent schema + 生成代码
└──────────────────────────────────────────────┘
```

### 规则

1. **Engine 层禁止 import ent** — 输入/输出都是纯 Go struct
2. **Data 层禁止 import engine** — 只做数据存取，不含业务逻辑
3. **Service 层** 同时依赖 engine 和 data，负责编排
4. **Server 层** 只依赖 service，不直接碰 engine/data

### fate 目标架构

```
fate/
├── ent/                    ← Model 层：ent schema + 生成 (不变)
├── config/                 ← 配置 (不变)
├── engine/                 ← 算法层：零DB依赖，纯computation
│   ├── bazi/
│   ├── chronosfate/
│   ├── wuxing/
│   ├── wuge/
│   ├── zhouyi/
│   ├── naming/             ← Name generation + rating (依赖 repository 接口)
│   └── rating/             ← 5维评分
├── data/                   ← 数据层：只做存取
│   ├── repository/         ← DB 查询 (从 internal/repository 迁移)
│   ├── dict/               ← 笔画/部首/拼音查找
│   └── seeddb/             ← JSON → DB 导入
├── service/                ← 业务层：编排 engine + data
│   ├── analysis/           ← 组装 NameResult (结合 engine 输出 + character 数据)
│   ├── session/            ← 异步命名 session
│   └── filter/             ← 筛选器
├── server/                 ← Web 层
│   ├── handler/            ← HTTP handler (从 internal/http 迁移)
│   └── middleware/
├── cmd/                    ← 命令
│   ├── server/
│   ├── console/
│   └── dbinit/
├── resources/
└── fate.go                 ← 顶层 API facade
```

---

## 繁简诗词匹配方案

### 方案对比

| 方案 | 原理 | 优点 | 缺点 |
|------|------|------|------|
| A: 双份索引 | 对每个字符存 简+繁 两条PoemChar | 查询快，无需JOIN | 数据膨胀 1.3x |
| B: 变体映射表 | char_variant 表存映射 → 查询时JOIN | 存储小 | 查询慢，需两步 |
| C: 运行时 OpenCC | 查询时实时转简繁 | 存储最小 | 每次查询都转，慢 |
| **D: 规范字索引** | 存储时统一转简体，查询也是简体 | 存储精炼，查询简单 | 丢失原文 | 

### 推荐方案：D（规范字索引）+ 保留原文

**流程：**

```
诗词原文 "國破山河在"
    │
    ▼
extractCharRefs("國破山河在")
    ├── char="國"  → normalize("國") → "国" → PoemChar{char:"国", context:"國破山河在"}
    ├── char="破"  → normalize("破") → "破" → PoemChar{char:"破", ...}
    ├── char="山"  → normalize("山") → "山" → PoemChar{char:"山", ...}
    ├── char="河"  → normalize("河") → "河" → PoemChar{char:"河", ...}
    └── char="在"  → normalize("在") → "在" → PoemChar{char:"在", ...}
```

**关键设计：**

```go
// pipeline/mapping/mapping.go

// CharacterMapping 简繁/异体映射表
type CharacterMapping struct {
    Simplified string   // 规范简体字
    Variants   []string // 所有变体（繁体、异体）
}

// Normalizer 规范字转换器
type Normalizer struct {
    variantToCanonical map[string]string // 國→国, 後→后, ...
}

func NewNormalizer(jsonPath string) (*Normalizer, error) {
    // 从 character.json 或 OpenCC 数据加载映射
    // character.json 已有: character.simplified_of_char, character.variant_of_char
}

// Normalize 将任意变体转换为规范简体
func (n *Normalizer) Normalize(char string) string {
    if canonical, ok := n.variantToCanonical[char]; ok {
        return canonical
    }
    return char // 已是规范字，直接返回
}

// Expand 展开一个规范字的所有变体（用于查询）
func (n *Normalizer) Expand(char string) []string {
    // 用于：用户查"国"时需要同时查 国/國/囯/...
}
```

**PoemChar schema 调整：**

```go
// PoemChar - 新增字段
type PoemChar struct {
    ent.Schema
}

func (PoemChar) Fields() []ent.Field {
    return []ent.Field{
        field.Int("id"),
        field.Int("poem_id"),
        field.String("char").NotEmpty().Comment("规范简体字"),
        field.String("original_char").Optional().Comment("原文中的字（繁/简/异）"),
        field.Int("position"),
        field.String("sentence").Optional(),
        field.String("context").Optional(),
    }
}
```

**查询时：**

```go
// service/poetry/query.go
func (s *Service) FindCharPoetry(char string) ([]PoetrySource, error) {
    // char 已经是简体（来自fate character）
    // 直接查规范字索引即可
    return s.client.PoemChar.Query().
        Where(poemchar.CharEQ(char)).
        WithPoem().
        All(ctx)
}
```

为什么选方案 D：
1. **character.json 已经以简体为主**，fate 的命名字库都是简体
2. **查询简单**：单表查询，无需 JOIN
3. **存储精炼**：一个汉字一条索引
4. **原文不丢**：`original_char` 和 `context` 保留，前端展示时可用原文
5. **规范字表已有**：从 character.json 的 `simplified_of_char` + `variant_of_char` 即可构建

### 规范字表生成流程

```
character.json (fate)
    ├── 8200+ 汉字
    ├── has simplified_of_char / variant_of_char 关系
    │
    ▼
pipeline/mapping/builder.go
    → 解析 character.json
    → 构建 variantToCanonical map
    → 输出 mapping.json (供 import-poetry 使用)
    │
    ▼
import-poetry 使用 mapping.json
    → 每个提取的汉字先 normalize
    → 写入 PoemChar.char = 规范简体
    → 写入 PoemChar.original_char = 原文
```

---

## 实施路径

### Phase 1: 繁简诗词匹配

1. **构建映射表**
   - 源：`character.json` 的 `simplified_of_char` + `variant_of_char` 字段

2. **重构导入工具**
   - `extractCharRefs` 中每个汉字先 `normalize()` 再存储
   - `PoemChar` schema 加 `original_char` 字段 (ALTER TABLE 兼容)

3. **验证**
   - 查 `国` 能匹配到含 `國` 的诗词
   - 查 `乐` 能匹配到含 `樂` 的诗词

### Phase 2: fate 架构重整

1. **提取 engine 层** — `internal/{bazi,chronosfate,wuxing,wuge,zhouyi,naming,rating}` → `engine/`
2. **提取 data 层** — `internal/{repository,dict,seeddb}` → `data/`
3. **提取 service 层** — `internal/{analysis,session,filter}` → `service/`
4. **提取 server 层** — `internal/http` → `server/handler`
5. **清理 cmd**: 移除冗余命令
