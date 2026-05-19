# Fate 起名系统 — 功能规划

> **注意**：本文档中描述的以下功能已迁移至 **qiming** 项目：
> - 古诗词/名言名句/名著来源取名（internal/poetry、ent/schema/poem.go、poem_char.go）
> - 浏览页面/Web API（internal/api、cmd/server）
> - 数据抓取（cmd/fetchpoem、cmd/fetchdata）
> - 报告生成（cmd/gen_report、internal/analysis/format.go、report.go）
> - 星座详情（internal/analysis/constellation.go）
> - 周易详情（internal/analysis/zhouyi_data.go）
> - 数据迁移（cmd/transfer、transfer/）
>
> fate 现仅保留核心起名计算能力。以下内容保留作为历史参考。

## 一、古诗词/名言名句/名著来源取名

### 1.1 目标
用户可以选择"诗词取名"模式，系统从古诗词、名言名句、经典名著中提取适合起名的字词，结合八字五行和五格数理进行推荐。

### 1.2 数据模型（已有）
- `poem` 表：存储诗词/名言/名著条目（title, author, dynasty, content, keywords, tags, type）
- `poem_char` 表：存储诗词中每个可起名字的汉字（char, position, sentence, context）
- `type` 枚举：shi(诗), ci(词), fu(赋), jing(经), other(其他)

### 1.3 需要新增的功能

#### 1.3.1 数据导入
- [ ] `cmd/seeddb` 增加诗词数据导入（唐诗三百首、宋词三百首、诗经、楚辞、论语、道德经等）
- [ ] 支持从 JSON/CSV 批量导入诗词数据
- [ ] 自动分词提取可起名字的汉字，写入 poem_char 表
- [ ] 名言名句数据源：增广贤文、菜根谭、围炉夜话等
- [ ] 名著数据源：红楼梦、西游记、三国演义、水浒传中适合起名的字词

#### 1.3.2 诗词取名逻辑
- [ ] 新增 `internal/poetry` 包，提供诗词取名功能
- [ ] `PoetryNamer` 接口：
  - `FindByChar(char string) []PoemCharResult` — 查找包含指定字的诗词
  - `FindByWuXing(wuxing string, limit int) []PoemCharResult` — 按五行查找适合的字
  - `FindPair(wuxing1, wuxing2 string) []PoemPairResult` — 查找同句/同联中五行匹配的字对
  - `GetSentence(poemID int, position int) string` — 获取字的上下文句子
- [ ] 诗词取名评分维度：
  - 出处权威性（诗经 > 唐诗 > 宋词 > 其他）
  - 句意美好度（含"吉"关键词加分）
  - 五行匹配度（与八字喜用神匹配）
  - 上下文连贯性（两字来自同一句/同一联加分）

#### 1.3.3 集成到生成流程
- [ ] `FilterOption` 新增 `PoetryMode` 字段：none(关闭), prefer(优先), only(仅诗词)
- [ ] session.generate() 中，当 PoetryMode 开启时：
  1. 先从 poem_char 表查询五行匹配的字
  2. 优先组合同句/同联的字对
  3. 再与普通生成结果合并/替换
- [ ] `NameResult` 新增 `PoetrySource` 字段：
  ```go
  type PoetrySource struct {
      Title    string // 诗名
      Author   string // 作者
      Dynasty  string // 朝代
      Sentence string // 出处句子
      Type     string // 体裁
  }
  ```

#### 1.3.4 命令行支持
- [ ] `name generate --poetry=prefer` — 优先诗词取名
- [ ] `name generate --poetry=only` — 仅从诗词中取名
- [ ] `name detail` 输出中显示诗词来源信息

### 1.4 实施步骤
1. 创建 `internal/poetry` 包，实现查询逻辑
2. 在 `internal/repository` 添加 Poem/PoemChar 的查询方法
3. 扩展 FilterOption 和 session.generate()
4. 扩展 NameResult 和 analysis 包
5. 更新 cmd/console/name.go
6. 准备种子数据并导入

---

## 二、浏览页面

### 2.1 目标
提供 Web 界面，用户可以：
- 在线输入生日/姓氏/性别，生成名字
- 浏览 Top10 详细分析
- 查看所有候选名字
- 点击名字查看完整分析报告
- 浏览古诗词来源

### 2.2 技术方案

#### 方案A：嵌入式 Web 服务（推荐）
- 使用 Go 标准库 `net/http` 或 `gin`/`echo` 提供 API
- 前端使用嵌入式静态文件（embed.FS）
- 单二进制部署，无需额外依赖

#### 方案B：独立前端 + API
- 后端提供 RESTful API
- 前端独立部署（React/Vue）

### 2.3 API 设计
```
POST /api/generate
  Body: { "surname": "张", "born": "2024/06/15 10:30", "sex": "boy", "poetry_mode": "prefer" }
  Response: { "top_names": [...], "total": 1234, "task_id": "xxx" }

GET /api/generate/:task_id/status
  Response: { "state": "finished", "progress": 100 }

GET /api/generate/:task_id/names?page=1&size=20
  Response: { "names": [...], "total": 1234 }

GET /api/name/detail?surname=张&char1=瑞&char2=霖&born=2024/06/15 10:30&sex=boy
  Response: { "name_result": {...} }

GET /api/poetry/search?char=瑞
  Response: { "results": [...] }
```

### 2.4 前端页面
- [ ] 首页：输入表单（姓氏、生日、性别、选项）
- [ ] 结果页：Top10 卡片 + 全部名字列表（分页）
- [ ] 详情页：完整分析报告（五格、三才、周易、诗词来源）
- [ ] 诗词浏览页：按朝代/作者/关键词浏览

### 2.5 实施步骤
1. 选择 Web 框架（建议 gin + embed.FS）
2. 创建 `internal/api` 包，实现 RESTful API
3. 创建 `cmd/server/main.go`，启动 Web 服务
4. 实现前端页面（HTML/CSS/JS，嵌入到二进制中）
5. 集成名字生成和详情查询 API
6. 添加诗词搜索 API

---

## 三、优先级排序

| 优先级 | 功能 | 预计工作量 | 状态 |
|--------|------|-----------|------|
| P0 | 诗词取名核心逻辑（internal/poetry） | 中 | 已迁移至 qiming |
| P0 | 诗词数据导入 | 中 | 已迁移至 qiming |
| P1 | 集成到生成流程 | 小 | 已迁移至 qiming |
| P1 | Web API 层 | 中 | 已迁移至 qiming |
| P2 | 前端页面 | 大 | 已迁移至 qiming |
| P2 | 诗词浏览功能 | 小 | 已迁移至 qiming |
