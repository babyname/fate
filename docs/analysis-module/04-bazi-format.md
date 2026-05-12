# 八字格式化

## FormatBazi() 实现

```go
package analysis

func FormatBazi(fateData *FateData) string {
    return fmt.Sprintf("八字: %s %s %s %s",
        dateToString(fateData.Bazi.Year),
        dateToString(fateData.Bazi.Month),
        dateToString(fateData.Bazi.Day),
        dateToString(fateData.Bazi.Hour),
    )
}
```

---

## 格式说明

| 柱位 | 内容 |
|-----|------|
| 年柱 | 天干+地支 |
| 月柱 | 天干+地支 |
| 日柱 | 天干+地支 |
| 时柱 | 天干+地支 |

---

## 总结

八字格式化输出四柱天干地支，格式为"天干+地支"。

**核心接口**：FormatBazi()
**输出格式**：年柱 月柱 日柱 时柱