# 名字评分实现

## rate.go 实现

```go
package naming

func RateNames(names []NameInfo, fateData *FateData) (*RatedNames, error) {
    rated := []RatedName{}
    
    for _, name := range names {
        wuxingScore := calcWuxingScore(name, fateData)
        bihuaScore := calcBihuaScore(name)
        yinyunScore := calcYinyunScore(name)
        
        totalScore := wuxingScore*0.4 + bihuaScore*0.3 + yinyunScore*0.3
        
        rated = append(rated, RatedName{Name: name, Score: totalScore})
    }
    
    sort.Slice(rated, func(i, j int) bool {
        return rated[i].Score > rated[j].Score
    })
    
    return &RatedNames{Names: rated}, nil
}
```

---

## 辅助函数

```go
func calcWuxingScore(name NameInfo, fateData *FateData) float64 {
    xiWuxing := fateData.WuxingXiji.XiWuxing
    matchCount := 0
    for _, wuxing := range name.Wuxing {
        if contains(xiWuxing, wuxing) {
            matchCount++
        }
    }
    return float64(matchCount) / float64(len(name.Wuxing)) * 100
}
```

---

## 总结

名字评分实现包括五行评分、笔画评分、音韵评分，综合评分后排序。

**核心实现**：RateNames()
**辅助函数**：calcWuxingScore()、calcBihuaScore()、calcYinyunScore()