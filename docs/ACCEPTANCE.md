# FATE 强制自动化验收规范

## 0. 为什么

人工在浏览器里点一遍无法证明正确性：看不到真实状态码、响应体与服务端日志，也无法复现。本仓库任何代码改动、修复、配置变更都**必须自带可复现的自动化验收**，并把证据落盘。范式参考 `orig-hub/docs/ACCEPTANCE.md`。

## 1. 触发条件

以下任一情形即触发：Go 源码改动、`web/` 源码改动、HTTP 契约改动、发布配置（`.goreleaser.yml` / `release.yml` / `ci.yml`）改动、示例与文档契约改动、版本变更。

## 2. 硬指标（全部满足才判 PASS）

1. `python scripts/check_repo.py` 全绿；`--self-test` 必中反例与必不中反例全绿。
2. `go build ./...`、`go vet ./...`、`go test ./... -count=1` 退出码 0。
3. 涉及 `Generate` / filter 的改动：额外 `go test -tags integration . -run TestGenerateNilFilterCompletes -v`，必须在 60s 内返回非空结果。
4. 涉及 Web 契约的改动：`bash scripts/verify-release.sh` 全绿，含 8 个端点的**真实状态码 + 响应体**、`names` 必须为 JSON 数组、SPA fallback 契约。
5. 涉及示例或文档的改动：`bash scripts/verify-examples.sh` 全绿，5 个示例逐个真实 `go run` 并断言 stdout。
6. 涉及发布的改动：`bash scripts/verify-release-assets.sh` 全绿，含 Release API 响应、资产清单、sha256 逐一校验、workflow 结论、CI 同版 `goreleaser check`。
7. UI / 截图改动：真实浏览器渲染截图（`scripts/capture-screenshots.mjs`），数字断言不能替代亲眼看图。

## 3. 证据落盘

统一写到 `tests/acceptance/<version>/evidence/`：每个断言一条 `.log` / `.json` / `.status` / `.body`，汇总报告为 `report.md` 或 `*-verification.md`。证据文件随代码一起提交，`.gitignore` 已为 `tests/acceptance/**/evidence/*.log` 开放例外。

证据必须包含：真实退出码、真实 HTTP 状态码与响应体、服务端日志、断言的期望值与实际值。只写「PASS」不算证据。

## 4. 交付约束

- 汇报只给四段：依据哪条规则 / 改了什么 / 验收结果 / 唯一待拍板项。过程性技术观察写进项目记忆或 issue，不进汇报。
- 看不到效果就不能声称完成：附带可查看的证据路径。
- 结论带证据，避免「可能 / 也许」。
- 用户侧只做验收，不运行任何开发工具命令；dev 服务、构建、脚本全生命周期由自动化脚本管理。

## 5. 禁止项

- 禁止用 `grep` 不到 / 编译不报错 当作「已修复」：检测器必须先用必中反例自证能命中。
- 禁止把瞬时态（连接超时、端口未起、依赖漂移）当成产品缺陷下结论；先稳定环境再测。
- 禁止 only-green 断言（永远为真的条件）、禁止跳过失败用例、禁止 `--no-verify`。
- 禁止用「我这边复现不出」对抗用户观察到的事实；先改进观测方法。
