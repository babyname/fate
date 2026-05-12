# analysis 模块介绍

## 模块定位

**analysis** 是 fate 项目的**输出格式化模块**，负责：
- 格式化输出八字分析结果
- 生成五行喜忌报告
- 格式化名字推荐列表

---

## 核心职责

| 职责 | 说明 |
|-----|------|
| **八字格式化** | 格式化四柱、天干地支 |
| **五行格式化** | 格式化五行喜忌结果 |
| **名字格式化** | 格式化名字推荐列表 |

---

## 核心接口

```go
package analysis

func FormatBazi(fateData *FateData) (string, error)
func FormatWuxingXiji(fateData *FateData) (string, error)
func FormatNames(ratedNames *RatedNames) (string, error)
func FormatReport(fateData *FateData, ratedNames *RatedNames) (string, error)
```

---

## 总结

analysis 模块负责格式化输出八字分析、五行喜忌和名字推荐。

**核心职责**：格式化输出
**核心接口**：FormatReport()、FormatBazi()、FormatWuxingXiji()