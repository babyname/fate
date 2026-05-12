# 名字推荐设计

## 推荐原理

综合筛选和评分，推荐最佳名字列表。

---

## RecommendNames() 接口

```go
package naming

func RecommendNames(fateData *FateData) (*RatedNames, error) {
    // 1. 筛选名字
    xiWuxing := fateData.WuxingXiji.XiWuxing
    surname := fateData.Surname
    names, err := FilterNames(xiWuxing, surname)
    if err != nil {
        return nil, err
    }
    
    // 2. 评分名字
    rated, err := RateNames(names, fateData)
    if err != nil {
        return nil, err
    }
    
    // 3. 推荐前10名
    rated.Names = rated.Names[:10]
    
    return rated, nil
}
```

---

## 推荐要点

| 要点 | 说明 |
|-----|------|
| 篮选名字 | 根据五行喜忌筛选 |
| 评分名字 | 综合评分排序 |
| 推荐数量 | 推荐前10名 |

---

## 总结

名字推荐综合筛选和评分，推荐前10名最佳名字。

**核心接口**：RecommendNames()
**推荐流程**：筛选 → 评分 → 推荐