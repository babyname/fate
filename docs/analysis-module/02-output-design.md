# 输出设计

## 输出格式

支持多种输出格式：文文本、JSON、HTML。

---

## FormatReport() 接口

```go
package analysis

func FormatReport(fateData *FateData, ratedNames *RatedNames) (string, error) {
    var report strings.Builder
    
    report.WriteString(FormatBazi(fateData))
    report.WriteString(FormatWuxingXiji(fateData))
    report.WriteString(FormatNames(ratedNames))
    
    return report.String(), nil
}
```

---

## 输出要点

| 要点 | 说明 |
|-----|------|
| 支持格式 | 文文本、JSON、HTML |
| 拼接顺序 | 八字 → 五行 → 名字 |

---

## 总结

输出设计支持文文本、JSON、HTML格式，按八字→五行→名字顺序拼接。

**支持格式**：文文本、JSON、HTML