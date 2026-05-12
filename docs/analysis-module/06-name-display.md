# 名字展示格式化

## FormatNames() 实现

```go
package analysis

func FormatNames(ratedNames *RatedNames) string {
    if len(ratedNames.Names) == 0 {
        return "推荐名字: 无"
    }
    
    var names strings.Builder
    names.WriteString("推荐名字:\n")
    
    for i, ratedName := range ratedNames.Names {
        names.WriteString(fmt.Sprintf("%d. %s (评分: %.2f)\n",
            i+1, ratedName.Name.Char, ratedName.Score))
    }
    
    return names.String()
}
```

---

## 格式说明

| 内容 | 说明 |
|-----|------|
| 名字列表 | 按评分排序 |
| 评分显示 | 保留兩位小数 |

---

## 总结

名字展示格式化输出推荐名字列表，按评分排序。

**核心接口**：FormatNames()
**输出格式**：编号. 名字 (评分: xx.xx)