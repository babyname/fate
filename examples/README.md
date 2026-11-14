# FATE 示例

五个独立可运行的 Go 示例，全部基于 v4.3.0 的公开 API（`fate.New` → `Generate` / `GenerateStream`），不依赖测试框架或第三方 HTTP 客户端。旧版 `example/create_a_name/` 是 2022 年的空壳（引用的 `NewFate` / `MakeName` 等 API 已不存在），已删除。

## 前置条件

**所有示例都从仓库根目录运行**。`config.DefaultConfig()` 使用同目录下的 SQLite 文件 `fate`（`Database.Name = "fate"`）。首次运行前初始化一次：

```bash
go run ./cmd/dbinit
```

`http-client` 还需要先启动 Web 服务：

```bash
go run ./cmd/server -addr :18080
```

## 示例清单

| 目录 | 运行命令 | 演示的 API | 需要服务 |
| :--- | :--- | :--- | :--- |
| `quickstart/` | `go run ./examples/quickstart` | `fate.New` + `Generate` + `TopNames`，最小主链路（与 Web 服务同源） | 否 |
| `stream/` | `go run ./examples/stream` | `GenerateStream` 的增量批次、并发回调保护、唯一 `IsFinal` 终批 | 否 |
| `filters/` | `go run ./examples/filters` | `FilterOption` 与链式 `With*`：避字、笔画、繁简、喜用神、严格度，并校验避字生效 | 否 |
| `explore/` | `go run ./examples/explore` | `ExcellentTable.TopN / MarkShown / Explore`，`nil` 表示取尽 | 否 |
| `http-client/` | `go run ./examples/http-client` | 标准库调用 `/api/v1/generate → status → result → names → explore` | 是 |

`http-client` 可用 `FATE_BASE_URL` 覆盖服务地址，例如 `FATE_BASE_URL=http://127.0.0.1:18081 go run ./examples/http-client`。

## 契约要点

- HTTP 端点统一在 `/api/v1/*` 下；`generate` 是同步的，直接返回 `task_id` 与 `total`。
- `avoid_chars` / `require_chars` 是**逗号分隔字符串**，不是 JSON 数组。
- `explore` 的 `names` 恒为 JSON 数组；候选取尽时返回 `[]`，不会是 `null`。
- `FilterOption` 为 `nil` 时，`Generate` / `GenerateStream` 使用 `fate.DefaultFilterOption()`（有界、与 Web 表单一致），**不会**退回宽松的共享 `fate.DefaultFilter()`（后者候选组合无限膨胀、实测 10 分钟不收敛）。定制时以 `fate.DefaultFilterOption().WithAvoidCharacters(...)` 为基座，不要再写一份魔法数字。
- `GenerateStream` 的 `onBatch` 可能由多个 worker 并发调用，回调内共享状态必须自己加锁。

完整的 Web API 表、请求/响应字段与库集成说明见 [docs/DEVELOPMENT.md](../docs/DEVELOPMENT.md)。

## 自动化验收

```bash
bash scripts/verify-examples.sh
```

脚本会逐个真实运行五个示例（自动拉起本地服务、捕获退出码与 stdout 断言），并把日志与结论写入 `tests/acceptance/v4.3.0/evidence/`。
