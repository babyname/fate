# 基础示例

## 示例1：新生儿起名

```go
import "github.com/godcong/fate/chronos"

input := &chronos.FateInput{
    Birthday: "2024-01-01 08:30",
    Surname:  "张",
}

data, _ := chronos.GetFateData(input)
rated, _ := naming.RecommendNames(data)
report, _ := analysis.FormatReport(data, rated)
fmt.Println(report)
```

---

## 示例2：使用配置

```go
cfg, _ := config.Load("config.yaml")
fmt.Println(cfg.Chonos.DataRangeStart)
```

---

## 总结

基础示例包括新生儿起名和使用配置。

**示例内容**：新手起名、配置加载