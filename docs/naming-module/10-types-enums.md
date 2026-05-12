# 命名模块类型定义

## types.go

```go
package naming

type NameInfo struct {
    Char   string
    Wuxing string
    Bihua  int
}

type RatedName struct {
    Name  NameInfo
    Score float64
}

type RatedNames struct {
    Names []RatedName
}
```

---

## 类型要点

| 类型 | 说明 |
|-----|------|
| NameInfo | 汉字信息（字、五行、笔画） |
| RatedName | 评分名字（名字+评分） |
| RatedNames | 评分名字列表 |

---

## 总结

类型定义包括 NameInfo、RatedName、RatedNames，用于存储汉字信息和评分结果。

**核心类型**：NameInfo、RatedName、RatedNames