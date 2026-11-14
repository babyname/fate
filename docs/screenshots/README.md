# 界面截图

README 引用的截图实体全部保存在本目录，均来自**真实运行中的 Web 界面**（v4.3.0），不是设计稿或占位图。

| 文件 | 内容 |
|------|------|
| `home.png` | 首页：姓氏、出生日期时刻、性别、诗词模式输入与「开始命名」 |
| `results.png` | 结果页：精选好名、更多推荐、候选名册入口 |
| `table.png` | 表格视图：候选名逐项对比，含五格吉凶列与查同名入口 |
| `cards.png` | 卡片视图：五格吉凶与查同名入口 |

## 重新生成

截图由可复现脚本捕获，使用 Chromium DevTools Protocol，无需额外测试框架（Node 22+ 自带 WebSocket）：

```bash
# 1. 启动服务（另开终端）
go run ./cmd/server

# 2. 捕获四张截图（需本机安装 Chrome 或 Edge）
node scripts/capture-screenshots.mjs
```

脚本会真实调用 `POST /api/v1/generate` 创建任务，再依次渲染首页、结果页、表格视图与卡片视图。可用环境变量覆盖：`FATE_BASE`（服务地址）、`CHROME_PATH`（浏览器路径）、`CDP_PORT`（调试端口）。

如果界面结构变化导致脚本找不到「卡片」按钮或渲染条件不满足，脚本会直接失败并打印原因，避免静默产出空白截图。
