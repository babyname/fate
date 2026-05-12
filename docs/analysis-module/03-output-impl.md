# 输出实现

## analysis.go 实现

```go
package analysis

func FormatReport(fateData *FateData, ratedNames *RatedNames) (string, error) {
    bazi := FormatBazi(fateData)
    wuxing := FormatWuxingXiji(fateData)
    names := FormatNames(ratedNames)
    
    return bazi + "\n" + wuxing + "\n" + names, nil
}
```

---

## 辅助函数

```go
func FormatBazi(fateData *FateData) string {
    return fmt.Sprintf("八字: %s", dateToString(fateData.Bazi))
}

func FormatWuxingXiji(fateData *FateData) string {
    return fmt.Sprintf("喜用五行: %s", strings.Join(fateData.WuxingXiji.XiWuxing, ", "))
}
```

---

## 总结

输出实现调用 FormatBazi()、FormatWuxingXiji()、FormatNames() 拼接输出。

**核心实现**：FormatReport()
**辅助函数**：FormatBazi()、FormatWuxingXiji()