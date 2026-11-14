# 更新日志

## v4.3.0 (2026-11-15)

**修复**

- 五格「半吉」按独立档位计分，不再被 `Contains("吉")` 误判为「吉」加分；补有效回归断言（PR #145）。
- `/api/v1/explore` 在候选耗尽时返回 JSON 数组 `[]`，不再是 `null`。
- `fate.ExcellentTable` 改为指针别名，`namesvc.go` 与 `handler.go` 不再复制含 `sync.RWMutex` 的值；`go vet ./...` 归零。
- Web 前端 API 基址与 5 个端点对齐后端 `/api/v1/*`；`avoid_chars` / `require_chars` 数组转逗号串。
- 前端类型检查归零：移除未使用的 `setStoreTop10`、不存在的 `@/components/ui/label` 引用及未使用 import。
- 服务端启动信息与实际注册的 8 个端点一致。

**功能**

- 候选名册与详情展示五格吉凶；新增官方同名人数查询（复制全名、打开官方入口、结果仅存本机浏览器）。
- 探索交互：分批加载、错误重试、耗尽状态提示；分页与 `NameResult` 类型收敛。

**清理**

- 删除 37 个历史前端 bundle（约 1.6 MiB）与磁盘上约 133 MiB 编译/临时产物；`rsbuild` 开启 `cleanDistPath`，构建时自动清空产物目录。
- 删除零引用的第二套评分实现 `internal/naming/{naming,raters,stroke}.go`、死包 `internal/bazi` / `internal/zhouyi`、`internal/session` 的缓存/链表死代码、孤儿 `cmd/querychar` 与 `generate.sh`。
- 收敛 `resources/embed.go` 双份嵌入；删除 `rating.go` 3 个死函数、`ExcellentTable` 2 个死方法及只读的 `index` 映射、`config.DatabaseConfig` 废弃别名；`go mod tidy` 降级 `golang.org/x/text` 为间接依赖。
- 下线两个没有后端实现的 UI 入口（姓名测分、诗词搜索）及其死类型，避免点击后 404；恢复计划见 [#146](https://github.com/babyname/fate/issues/146)。
- 修正 `.gitignore` 误忽略 `go.sum` 的规则，忽略 `.workbuddy/`。

**文档与验收**

- `README.md` / `README_EN.md` 重写为面向用户的入口；开发者内容迁至 `docs/DEVELOPMENT.md`。
- 补齐 `docs/screenshots/` 的 4 张真实界面截图，脚本可复现：`scripts/capture-screenshots.mjs`。
- 新增全量验收脚本 `scripts/verify-release.sh`，证据落盘 `tests/acceptance/v4.3.0/evidence/`。
- 版本号统一：`VERSION` / `DB_VERSION` / `version.go` = `4.3.0` / `v4.3.0`，徽章改为动态 release badge。

**已知遗留**：公开 API 兼容边界（`fate.Session`、流式 API）、config YAML 契约、集成测试对未跟踪 `fate_test.db` 的依赖、CI 门禁缺失 —— 见 [#146](https://github.com/babyname/fate/issues/146)。

完整分析见 [#147](https://github.com/babyname/fate/issues/147)。

---

## 早期版本计划（历史记录，未作为发布版本）

| 版本 | 计划内容 |
|---|---|
| v0.1.0 | 初始化项目结构、chronos 模块、五行喜忌算法 |
| v0.2.0 | naming 起名模块、名字筛选与评分 |
| v0.3.0 | analysis 输出模块、JSON/HTML 输出 |
| v1.0.0 | 模块集成、配置文件、完整测试 |
