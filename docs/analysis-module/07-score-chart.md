# 评分图表格式化

## 评分图表设计

| 评分项 | 格式 |
|-------|------|
| 五行评分 | ★★★★★ (5星制) |
| 笔画评分 | 分数/等级 |
| 综合评分 | 百分制 |

---

## 实现示例

```go
func FormatScoreChart(name RatedName) string {
    stars := ""
    fullStars := int(name.Score / 20)
    for i := 0; i < fullStars; i++ {
        stars += "★"
    }
    return fmt.Sprintf("%s (%.1f分) %s", name.Name.Char, name.Score, stars)
}
```

---

## 总结

评分图表使用星级制可视化评分结果。

**核心接口**：FormatScoreChart()
**图表格式**：★星级