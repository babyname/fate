# AGENTS.md — FATE 协作约定

本文件是仓库规范真源。范式参考 `orig-hub`：规范必须落成**可执行门禁**，写不进脚本的条款只是愿望。`scripts/check_repo.py` 是这些规则的判定实现，`--self-test` 用必中反例证明每个检测器真能命中。

## 0. 规范层级（冲突时按此裁决）

1. 本文件（可执行判据）
2. `docs/ACCEPTANCE.md`（验收标准）
3. `docs/DEVELOPMENT.md`（当前 Web / 库 API 契约）
4. `README.md` / `README_EN.md`（面向用户）
5. `docs/architecture/`、`docs/analysis-module/`、`docs/chronos-module/`、`docs/config-module/`、`docs/implementation/`、`docs/PHASE3*.md`、`docs/ROADMAP.md` —— **历史设计档案**，反映迁出时的设想，不构成当前契约。源码与门禁优先。

向第 5 层文档追加内容前先问：这条规则能不能写成 `check_repo.py` 的一个断言？能就写进门禁，不要再加一份会腐烂的散文。

## 1. 流程

- 分支：`main` 唯一源。开工先对齐 `origin/main`；交付 = 本地 main 提交 + 普通 push。禁止 force push、禁止改写已发布 tag。
- 提交信息：`<type>(<scope>): <subject>`，**英文**、**LF**（`git config core.autocrlf input`）。禁止 `--no-verify`。门禁以 `1f0cfa3` 为生效点，历史消息只统计不阻断。
- 验收：任何代码改动必须自带可复现自动化验收，捕获真实退出码 / 状态码 / 响应体 / 日志并落盘 `tests/acceptance/<version>/evidence/`。**不把验收甩给人工点测**。
- 发布：push `v*` tag 触发 `release.yml`；`prepare-db` 生成 DB 包，GoReleaser 构建二进制，`changelog.use: github-native` 自动生成 notes。禁止手写 release body、禁止人工 curated notes 文件。

## 2. 算法契约（不可静默改动）

| 规则 | 判据 |
|---|---|
| 五格半吉计分 | `internal/rating/rating.go` 精确 `switch`：`"吉"`/`"大吉"` = `+15*weight`，`"半吉"` = `+8*weight`，其余 `-5*weight`。禁止退回 `strings.Contains(lucky, "吉")` |
| 有界筛选 | `GenerateRequest.Filter == nil` → `fate.DefaultFilterOption()`（`filteroption.go`，唯一真相源）。**禁止**把共享 `DefaultFilter()` 作为 `Generate` 的默认值：其笔画无上界，候选笛卡尔积实测 10 分钟不收敛（2022s CPU / 1.1 GiB）；有界默认实测 6.9s |
| 候选表 | `fate.ExcellentTable` 是指针别名，内部含 `sync.RWMutex`，**禁止按值复制**（`go vet` copylock 门禁） |
| 探索取尽 | `/api/v1/explore` 的 `names` 恒为 JSON 数组；耗尽返回 `[]`，不得是 `null` |
| 地区表述 | 中国台湾地区 / 中国香港 / 中国澳门，不得省略国家前缀 |

## 3. 实现约束

- **HTTP 契约**：`internal/http/handler.go` 注册 7 个 `/api/v1/*` 端点；`avoid_chars` / `require_chars` 是**逗号分隔字符串**不是数组；旧 `/api/*` 端点不存在（未知路径由 SPA fallback 返回 200 HTML）。文档必须列出全部 7 条路由。
- **静态资源**：`resources/embed.go` 只允许 `staticFS` 一个 `//go:embed all:static`；不要恢复第二份 `//go:embed *`（会把临时备份打进二进制）。
- **前端产物**：Rsbuild 输出到被 git 跟踪的 `resources/static`，文件名是 content hash。`cleanDistPath: true` 必须保留；`web/bun.lock` 必须入库，安装依赖用 `bun install --frozen-lockfile`。门禁强制「每个被跟踪的 bundle 都被 `index.html` 引用」——历史上正是缺这条，js 文件从 7 累积到 36。
- **发布构建**：`CGO_ENABLED=0`。驱动是纯 Go 的 `sqlite3ent/sqlite3`（底层 `modernc.org/sqlite`），开 CGO 会让 `linux_arm64` 交叉编译触发 `runtime/cgo` 汇编错误。
- **GoReleaser v2**：`archives.formats`（不是废弃的 `format`）、`changelog.use: github-native`。配置改动后用 CI 同版 `goreleaser check` 验证。
- **测试**：`go test ./...` 需要 `fate_test.db`；真实 DB / 长耗时用例放 `//go:build integration`，由验收脚本显式 `go test -tags integration` 运行，不进入默认测试集。

## 4. 样式与标准

- Go：`gofmt` / `goimports`；新增或改动的 `.go` 文件**禁止 UTF-8 BOM**（历史 23 个带 BOM 的文件是欠账，改动到就顺手清除）；`go vet ./...` 零错误。
- TypeScript：`tsc --noEmit` 零错误；不得保留没有后端实现的 UI 入口（`NameScoreTab` / `PoetrySearchTab` 已下线）。
- 文档：不写虚构 API；示例必须能 `go run` 通过；用户文档中文、提交信息英文。
- 版本单一真相源：`VERSION` = `DB_VERSION` = `version.go` 的 `const Version`（门禁校验）。DB 资产按 `DB_VERSION` 打成 `fate-v<db>.zip` + `.sha256`。
- 验收证据日志由 `.gitignore` 的 `!tests/acceptance/**/evidence/*.log` 例外放行，其他 `*.log` 一律不得入库。

## 5. 门禁与自动化

| 入口 | 作用 |
|---|---|
| `python scripts/check_repo.py` | 全量静态门禁（污染 / 生成物 / 契约 / 版本 / BOM / gofmt / build / vet） |
| `python scripts/check_repo.py --staged` | pre-commit 暂存区门禁 |
| `python scripts/check_repo.py --ci` | CI 门禁（含最近一次提交范围与提交信息） |
| `python scripts/check_repo.py --self-test` | 检测器必中反例 + 必不中反例，证明门禁没有假阴性 |
| `bash scripts/verify-release.sh` | 发布级回归：Go / tsc / 前端构建 / 8 个 HTTP 端点真实响应 |
| `bash scripts/verify-examples.sh` | 5 个示例真实运行 + HTTP wire 证据 |
| `bash scripts/verify-release-assets.sh` | 已发布 Release 的 API / 资产 / sha256 / workflow / GoReleaser 配置 |

本地钩子：`git config core.hooksPath scripts/git-hooks`（`pre-commit` 是 `check_repo.py --staged` 的薄封装）。

CI 在 `.github/workflows/ci.yml`，触发 `main` 与 PR。**旧的 `build_linux.yml` / `build_macos.yml` / `build_windows.yml` 只触发 `master`，且会创建 `auto_build` Release 并写 Go 1.19 + CGO=1：它们是历史遗留，禁止「顺手」改成 `main` 触发**——激活即产生垃圾 Release。替代方案是统一走 `release.yml` + 本 CI。

## 6. 待拍板项规范

向用户抛待拍板项必须三段式，只给编号或短语视为汇报缺陷：①**因果**（现状 → 根因 → 不决定要付的代价）②**处置**（动哪个文件 / 有几种做法 / 影响面与回归风险）③**请求**（请拍的那一句：选 A/B、是否批准）。疑问句一律当作开放待定项，不得记为「已确认」。
