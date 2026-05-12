# 名字推荐实现

## naming.go 实现

```go
package naming

func RecommendNames(fateData *FateData) (*RatedNames, error) {
    xiWuxing := fateData.WuxingXiji.XiWuxing
    surname := fateData.Surname
    
    names, err := FilterNames(xiWuxing, surname)
    if err != nil {
        return nil, err
    }
    
    rated, err := RateNames(names, fateData)
    if err != nil {
        return nil, err
    }
    
    if len(rated.Names) > 10 {
        rated.Names = rated.Names[:10]
    }
    
    return rated, nil
}
```

---

## 实现要点

| 要点 | 说明 |
|-----|------|
| 调用筛选 | FilterNames() |
| 调用评分 | RateNames() |
| 推荐数量 | 前10名 |

---

## 总结

名字推荐实现调用筛选和评分，推荐前10名最佳名字。

**核心实现**：RecommendNames()
**调用流程**：筛选 → 评分 → 推荐