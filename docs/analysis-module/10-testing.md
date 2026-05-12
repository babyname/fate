# 分析模块测试设计

## 测试范围

| 测试类型 | 测试内容 |
|---------|---------|
| 单元测试 | FormatReport() 完整报告输出 |
| 格式测试 | JSON、HTML 格式验证 |

---

## 测试示例

```go
func TestFormatReport(t *testing.T) {
    fateData := mockFateData()
    ratedNames := mockRatedNames()
    
    report, err := FormatReport(fateData, ratedNames)
    assert.NoError(t, err)
    assert.Contains(t, report, "八字")
    assert.Contains(t, report, "喜用五行")
}
```

---

## 总结

测试设计包括完整报告输出测试和格式验证。

**测试要点**：FormatReport() 完整输出、JSON/HTML 格式输出