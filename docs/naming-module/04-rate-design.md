# 名字评分设计

## 评分原理

综合评分名字，包括五行评分、笔画评分、音韵评分。

---

## 评分权重

| 评分项 | 权重 |
|-----|------|
| 五行评分 | 40% |
| 笔画评分 | 30% |
| 音韵评分 | 30% |

---

## RateNames() 接口

```go
package naming

func RateNames(names []NameInfo, fateData *FateData) (*RatedNames, error) {
    rated := []RatedName{}
    
    for _, name := range names {
        // 计算五行评分
        wuxingScore := calcWuxingScore(name, fateData)
        
        // 计算笔画评分
        bihuaScore := calcBihuaScore(name)
        
        // 计算音韵评分
        yinyunScore := calcYinyunScore(name)
        
        // 综合评分
        totalScore := wuxingScore*0.4 + bihuaScore*0.3 + yinyunScore*0.3
        
        rated = append(rated, RatedName{
            Name: name,
            Score: totalScore,
        })
    }
    
    // 按评分排序
    sort.Slice(rated, func(i, j int) bool {
        return rated[i].Score > rated[j].Score
    })
    
    return &RatedNames{Names: rated}, nil
}
```

---

## 总结

名字评分综合评分五行、笔画、音韵，权重分别为40%、30%、30%。

**核心接口**：RateNames()
**评分权重**：五行40%、笔画30%、音韵30%