# 高级示例

## 自定义筛选

```go
func customFilter(names []NameInfo, wuxing []string) []NameInfo {
    var filtered []NameInfo
    for _, name := range names {
        for _, w := range wuxing {
            if name.Wuxing == w {
                filtered = append(filtered, name)
                break
            }
        }
    }
    return filtered
}
```

---

## 自定义评分

```go
func customRate(names []NameInfo, fateData *FateData) RatedNames {
    rated := RatedNames{}
    for _, name := range names {
        score := calculateScore(name, fateData)
        rated.Names = append(rated.Names, RadiateName{
            Name:  name,
            Score: score,
        })
    }
    return rated
}
```

---

## 总结

高级示例包括自定义筛选和自定义评分。

**示例内容**：自定义筛选、自定义评分