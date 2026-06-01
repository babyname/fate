# fate 开发计划 v2

> 基于 2026-05-13 项目现状制定，反映已完成的重构工作（chronos/v2 升级、xorm 移除、sqlite3ent/sqlite3 集成、Character schema 重构）。

---

## 一、项目现状

### 已完成

| 里程碑 | 状态 | 说明 |
|--------|------|------|
| chronos/v2 升级 | ✅ | 移除 fate/chronos，本地 chronos/v2 + FateData + Bridge |
| xorm 依赖移除 | ✅ | 删除 sancai.go/wuxing.go(xorm版)，改用 wuxing.SanCai 纯计算 |
| SQLite 驱动替换 | ✅ | mattn/go-sqlite3 → sqlite3ent/sqlite3（纯Go，CGO_ENABLED=0） |
| Character schema 重构 | ✅ | char_type enum → 5个bool字段(is_simplified/is_traditional/is_kangxi/is_variant/is_ancient) |
| ent 代码重新生成 | ✅ | 新 schema 生成完毕，全量字段名更新 |
| dict 包同步更新 | ✅ | CharEntry/QueryFilter/DictIndex 适配新字段 |
| 编译+测试通过 | ✅ | dict 8/8, rating 4/4, analysis 2/2 |

| WuGeLucky 预计算查找表 | ✅ | [31][31][]WuGeResult + sync.Once，3.89ns/op，0 B alloc |
| 美名腾风格分析输出 | ✅ | 字信息/五格/三才/基础运/成功运/人际关系/周易卦象/四项评分 |
| 两种喜用神算法 | ✅ | 平衡用神法 + 格局用神法（用神/喜神/忌神/仇神/闲神） |

### 待完成

| 优先级 | 任务 | 说明 |
|--------|------|------|
| P0 | 数据迁移 | 旧DB(30K字+1M五格+124三才) → 新ent schema |
| P0 | 数据质量修复 | 旧数据五行覆盖率31%、name_science仅10条 |
| P1 | Filter 重构 | 配置化开关 + dict 索引集成 |
| P1 | Session 重构 | 集成 chronos/v2 + dict 索引 |
| P1 | Rating 统一 | rating/ vs naming/raters 合并 |
| P1 | 周易卦象数据完善 | yi包GuaXiang.XiangYue为空，需补充64卦大象/事业/经商/求名/婚恋/决策解读 |
| P1 | 生肖星座中文化 | 当前输出Dragon/Gemini，需映射为龙/双子座 |
| P1 | 简繁体对照表 | 当前TraditionalChar直接取原字，需建立简繁映射 |
| P2 | 诗词库 | 数据导入 + 出处追溯 | — |
| P2 | CLI 完善 | console/fetchdata/character 功能补全 | — |

---

## 二、开发阶段定义

### Phase 0: 数据层（当前阶段）

**目标**：建立完整、准确的数据基础，所有上层模块可依赖

#### 0.1 旧数据导出工具

**流程**：
1. 编写 `cmd/migratedb` 工具，读取旧 SQLite3 数据库
2. 逐表导出为 JSON/CSV 中间格式
3. 生成数据质量报告

**成果物**：
- [ ] `cmd/migratedb/main.go` — 迁移工具主程序
- [ ] `migratedb/export.go` — 旧数据导出逻辑
- [ ] `migratedb/transform.go` — 数据转换逻辑（旧字段 → 新字段）
- [ ] `migratedb/import.go` — 新数据库导入逻辑
- [ ] `migratedb/report.go` — 数据质量报告生成
- [ ] `data/export/character.json` — 导出的字表数据
- [ ] `data/export/wu_ge_lucky.json` — 导出的五格吉凶数据
- [ ] `data/export/wu_xing.json` — 导出的三才五行数据
- [ ] `data/export/quality_report.json` — 数据质量报告

**数据映射规则**：

| 旧字段 | 新字段 | 转换逻辑 |
|--------|--------|----------|
| ch | char | 直接映射 |
| pin_yin | pinyin | string → []string（按空格/逗号分割） |
| wu_xing | wu_xing | 直接映射 |
| stroke | simplified_stroke | 旧 stroke → simplified_stroke |
| kang_xi_stroke | kangxi_stroke | 直接映射 |
| science_stroke | science_stroke | 直接映射 |
| is_kang_xi | is_kangxi | 直接映射 |
| simple_total_stroke | simplified_stroke | 直接映射 |
| traditional_total_stroke | traditional_stroke | 直接映射 |
| regular | regular | 直接映射 |
| name_science | nameable | name_science=1 → nameable=true |
| lucky | (删除) | 旧字段不再使用 |
| traditional_character | simplified_of edge | 需查关联字建立 edge |
| variant_character | variant_of edge | 需查关联字建立 edge |
| (新增) | is_simplified | 根据简繁映射推导 |
| (新增) | is_traditional | 根据简繁映射推导 |
| (新增) | is_kangxi | 从 is_kang_xi 映射 |
| (新增) | is_variant | 从 variant_character 非空推导 |
| (新增) | is_ancient | 默认 false |
| (新增) | common_level | 根据 regular + 频率表推导 |
| (新增) | gender_hint | 默认空，后续补充 |

**测试验证**：
- [ ] 导出数据行数与旧库一致（character=30060, wu_ge_lucky=1081344, wu_xing=124）
- [ ] 每个字段的转换逻辑单元测试
- [ ] 数据质量报告生成无报错

#### 0.2 数据质量修复

**流程**：
1. 分析数据质量报告，识别缺失/错误数据
2. 补充五行属性（当前覆盖率 31%→目标 90%+）
3. 补充康熙笔画修正（基于 dict/kangxi_stroke.go 已有 200+ 条）
4. 补充常用字等级（common_level）
5. 交叉验证：Unihan 数据 vs 旧库数据

**成果物**：
- [ ] `data/fix/wuxing_fix.json` — 五行属性补丁
- [ ] `data/fix/kangxi_fix.json` — 康熙笔画补丁
- [ ] `data/fix/common_level.json` — 常用字等级数据
- [ ] `migratedb/fix.go` — 数据修复逻辑
- [ ] 更新后的数据质量报告（覆盖率达标）

**数据来源**：
- Unihan 数据库（拼音、笔画、简繁映射）
- 康熙字典电子版（笔画验证）
- 通用规范汉字表（常用字等级）
- 旧数据库已有数据（交叉验证）

**测试验证**：
- [ ] 五行覆盖率 ≥ 90%（当前 31%）
- [ ] 康熙笔画覆盖率 ≥ 95%（当前 58%）
- [ ] 常用字等级覆盖率 ≥ 80%（当前 0%）
- [ ] 无矛盾数据（同一字不同来源数据冲突）

#### 0.3 新数据库初始化

**流程**：
1. 使用 ent 的 auto migration 创建新表结构
2. 导入修复后的数据
3. 建立 edge 关系（简繁关联、异体字关联）
4. 创建索引
5. 验证数据完整性

**成果物**：
- [ ] `fate_new.db` — 新 SQLite3 数据库文件
- [ ] `database/seed.go` — 数据库种子数据初始化
- [ ] `database/migrate.go` — 数据库迁移逻辑
- [ ] 数据完整性验证报告

**测试验证**：
- [ ] `CGO_ENABLED=0 go build ./...` 编译通过
- [ ] 新库可被 ent client 正常打开和查询
- [ ] 所有 ent CRUD 操作单元测试通过
- [ ] character 表行数 ≥ 30000
- [ ] wu_ge_lucky 表行数 = 1081344
- [ ] wu_xing 表行数 = 124
- [ ] edge 关系正确（简繁互查、异体字关联）

---

### Phase 1: 核心逻辑层

**目标**：Filter/Session/Rating 三个核心模块重构完成，可端到端完成起名流程

#### 1.1 Filter 系统重构

**流程**：
1. 定义 FilterConfig（配置化开关）
2. 重写 filter.go 使用 dict.DictIndex 替代 ent 查询
3. 实现 FilterChain 模式（可组合、可配置）
4. 集成五行喜忌过滤

**成果物**：
- [ ] `filter/config.go` — FilterConfig 定义
- [ ] `filter/chain.go` — FilterChain 实现
- [ ] `filter/wuxing_filter.go` — 五行喜忌过滤器
- [ ] `filter/stroke_filter.go` — 笔画过滤器
- [ ] `filter/sancai_filter.go` — 三才过滤器
- [ ] `filter/zodiac_filter.go` — 生肖过滤器
- [ ] 重构后的 `filter.go`
- [ ] `filter/filter_test.go` — 过滤器测试

**测试验证**：
- [ ] 单个过滤器单元测试
- [ ] FilterChain 组合测试
- [ ] 性能测试：过滤 30000 字 < 100ms
- [ ] 与旧 filter.go 结果对比（回归测试）

#### 1.2 Session 起名流程重构

**流程**：
1. 重写 Session 使用 chronos/v2 的 GetFateData()
2. 集成 dict.DictIndex 作为字库查询
3. 集成 FilterChain 作为过滤引擎
4. 集成 Rating 作为评分引擎
5. 实现 generate → filter → rate → recommend 完整流程

**成果物**：
- [ ] 重构后的 `session.go`
- [ ] 重构后的 `session_wuge.go`
- [ ] `session_test.go` — 端到端测试
- [ ] 起名流程时序图

**测试验证**：
- [ ] 端到端测试：输入出生信息 → 输出推荐名字列表
- [ ] 性能测试：完整起名流程 < 1秒
- [ ] 边界测试：单姓/复姓、单名/双名
- [ ] 异常测试：无效输入、空结果

#### 1.3 Rating 评分系统统一

**流程**：
1. 分析 rating/ 和 naming/raters 的功能重叠
2. 统一评分接口设计
3. 合并到 rating/ 包
4. naming/ 包只保留推荐逻辑

**成果物**：
- [ ] `rating/rater.go` — 统一 Rater 接口
- [ ] `rating/wuxing.go` — 五行评分器
- [ ] `rating/bihua.go` — 笔画评分器
- [ ] `rating/yinyun.go` — 音韵评分器
- [ ] `rating/composite.go` — 组合评分器
- [ ] `rating/rating_test.go` — 评分器测试
- [ ] 重构后的 `naming/naming.go`（仅推荐逻辑）

**测试验证**：
- [ ] 各评分器单元测试
- [ ] 组合评分权重测试
- [ ] 评分结果一致性验证
- [ ] 评分范围 0-100 验证

---

### Phase 2: 应用层

**目标**：诗词库、Analysis 输出、CLI 工具完善

#### 2.1 诗词库 + 出处追溯

**流程**：
1. 确定诗词数据源（唐诗三百首/宋词三百首/诗经/楚辞）
2. 设计 Poem/PoemChar schema（已有，需验证）
3. 编写数据导入工具
4. 实现出处追溯查询
5. 集成到起名推荐结果

**成果物**：
- [ ] `cmd/importpoem/main.go` — 诗词导入工具
- [ ] `poem/poem.go` — 诗词查询接口
- [ ] `poem/trace.go` — 出处追溯逻辑
- [ ] 诗词数据库文件
- [ ] `poem/poem_test.go` — 诗词查询测试

**测试验证**：
- [ ] 诗词数据导入完整性
- [ ] 出处追溯准确性
- [ ] 集成到推荐结果的端到端测试

#### 2.2 Analysis 输出模块完善

**流程**：
1. 定义输出模板格式
2. 实现文本/JSON/HTML 输出
3. 集成八字解析 + 名字推荐
4. 添加评分图表

**成果物**：
- [ ] `analysis/template.go` — 输出模板
- [ ] `analysis/formatter_text.go` — 文本格式化
- [ ] `analysis/formatter_json.go` — JSON 格式化
- [ ] `analysis/formatter_html.go` — HTML 格式化
- [ ] `analysis/analysis_test.go` — 输出测试

**测试验证**：
- [ ] 各格式输出正确性
- [ ] 模板渲染无错误
- [ ] 中文编码正确

#### 2.3 CLI 工具完善

**流程**：
1. 完善 console 命令（起名交互流程）
2. 完善 character 命令（字表管理）
3. 添加 config 命令（配置管理）

**成果物**：
- [ ] 重构后的 `cmd/console/`
- [ ] 重构后的 `cmd/character/`
- [ ] 新增 `cmd/config/`

**测试验证**：
- [ ] 各命令功能测试
- [ ] 命令行参数解析测试
- [ ] 交互流程测试

---

## 三、阶段依赖关系

```
Phase 0 (数据层)
  ├── 0.1 旧数据导出 ──→ 0.2 数据质量修复 ──→ 0.3 新数据库初始化
  │
  ▼
Phase 1 (核心逻辑层) [依赖 Phase 0 完成]
  ├── 1.1 Filter 重构 ──┐
  ├── 1.2 Session 重构 ──┤── 1.3 Rating 统一
  │                      │
  ▼                      ▼
Phase 2 (应用层) [依赖 Phase 1 完成]
  ├── 2.1 诗词库
  ├── 2.2 Analysis 完善
  └── 2.3 CLI 完善
```

**关键路径**：0.1 → 0.2 → 0.3 → 1.1 → 1.2 → 1.3

---

## 四、质量门禁

### 编译门禁

```bash
CGO_ENABLED=0 go build ./...
```

必须零错误通过。

### 测试门禁

```bash
CGO_ENABLED=0 go test ./...
```

所有测试必须通过，不允许 skip。

### 覆盖率要求

| 模块 | 最低覆盖率 |
|------|-----------|
| dict | 80% |
| rating | 80% |
| filter | 80% |
| session | 70% |
| wuxing | 80% |
| analysis | 70% |

### 数据质量门禁

| 指标 | 最低要求 |
|------|---------|
| 五行覆盖率 | ≥ 90% |
| 康熙笔画覆盖率 | ≥ 95% |
| 常用字等级覆盖率 | ≥ 80% |
| 拼音覆盖率 | ≥ 75% |

---

## 五、文档模板

### 5.1 模块开发文档模板

每个模块开发时，应产出以下文档：

```markdown
# {模块名} 模块

## 1. 概述
- 模块职责
- 在系统中的位置（上层/下层依赖）
- 核心接口

## 2. 设计
- 数据结构定义
- 算法说明
- 接口设计
- 错误处理策略

## 3. 实现
- 文件清单
- 关键实现说明
- 配置项

## 4. 测试
- 测试用例清单
- 测试数据说明
- 覆盖率
- 性能基准

## 5. 变更记录
| 日期 | 变更内容 | 影响范围 |
|------|---------|---------|
```

### 5.2 数据迁移文档模板

```markdown
# {表名} 数据迁移

## 1. 源数据
- 来源：旧数据库 / 外部数据源
- 表名：{旧表名}
- 行数：{N}
- 关键字段：{字段列表}

## 2. 目标数据
- 目标表：{新表名}
- Schema 版本：{ent schema 版本}
- 关键字段：{字段列表}

## 3. 字段映射
| 源字段 | 目标字段 | 转换规则 | 备注 |
|--------|---------|---------|------|

## 4. 数据质量
- 源数据问题：{问题描述}
- 修复策略：{修复方案}
- 修复后覆盖率：{N%}

## 5. 验证
- [ ] 行数一致
- [ ] 字段值范围正确
- [ ] 关联关系完整
- [ ] 无矛盾数据
```

### 5.3 测试报告模板

```markdown
# {模块名} 测试报告

## 1. 测试概要
- 测试日期：{YYYY-MM-DD}
- 测试范围：{模块/功能}
- 测试环境：Go {version}, CGO_ENABLED=0

## 2. 测试结果
| 测试项 | 用例数 | 通过 | 失败 | 跳过 |
|--------|-------|------|------|------|

## 3. 覆盖率
| 包 | 覆盖率 |
|----|-------|

## 4. 性能基准
| 操作 | 耗时 | 数据量 |
|------|------|-------|

## 5. 发现的问题
| 编号 | 严重程度 | 描述 | 状态 |
|------|---------|------|------|

## 6. 结论
- 是否通过质量门禁
- 遗留问题
```

### 5.4 阶段完成检查清单

```markdown
# Phase {N} 完成检查清单

## 代码
- [ ] 所有计划文件已创建/修改
- [ ] CGO_ENABLED=0 go build ./... 通过
- [ ] go vet ./... 无警告
- [ ] 所有测试通过
- [ ] 覆盖率达标

## 数据
- [ ] 数据迁移完成
- [ ] 数据质量达标
- [ ] 数据完整性验证通过

## 文档
- [ ] 模块开发文档已更新
- [ ] 数据迁移文档已更新
- [ ] 测试报告已生成

## 集成
- [ ] 与上下游模块接口对齐
- [ ] 端到端测试通过
- [ ] 性能指标达标

## 交付
- [ ] git push 成功
- [ ] 代码已合并到主分支
```

---

## 六、风险与应对

| 风险 | 概率 | 影响 | 应对策略 |
|------|------|------|---------|
| 旧数据五行覆盖率低 | 高 | 高 | 多数据源交叉补充 + 规则推导 |
| ent schema 变更导致迁移困难 | 中 | 高 | 中间格式解耦 + 版本化迁移 |
| 五格吉凶表数据量巨大(1M+) | 低 | 中 | 批量导入 + 事务优化 |
| 简繁关联数据不完整 | 中 | 中 | Unihan 数据补充 + 人工校验 |
| 诗词数据版权问题 | 低 | 高 | 使用公有领域数据源 |

---

## 七、术语表

| 术语 | 说明 |
|------|------|
| 五格 | 天格、人格、地格、外格、总格 |
| 三才 | 天、人、地三才五行组合 |
| 大衍之数 | 81个固定数值，用于五格吉凶判断 |
| 康熙笔画 | 康熙字典中的笔画数，起名专用 |
| 姓名学笔画 | 基于康熙字典，含部首变形修正的笔画数 |
| science_stroke | 姓名学笔画（= 康熙笔画 + 部首修正） |
| 喜用神 | 八字中对日主有利的五行 |
| 忌神 | 八字中对日主不利的五行 |
| 调候用神 | 根据月令气候调整的用神 |
| 日主强弱 | 日柱天干在八字中的旺衰程度 |
