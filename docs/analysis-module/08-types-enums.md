# 分析模块类型定义

## types.go

```go
package analysis

type Report struct {
    Bazi    string
    Wuxing  string
    Names   string
}

type OutputFormat int

const (
    TextFormat OutputFormat = iota
    JSONFormat
    HTMLFormat
)
```

---

## 类型要点

| 类型 | 说明 |
|-----|------|
| Report | 分析报告（八字+五行+名字） |
| OutputFormat | 输出格式枚举 |

---

## 总结

类型定义包括 Report 和 OutputFormat，用于存储分析报告和输出格式。

**核心类型**：Report、OutputFormat