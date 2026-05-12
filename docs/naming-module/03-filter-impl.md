# 名字筛选实现

## filter.go 实现

```go
package naming

func FilterNames(xiWuxing []string, surname string) ([]NameInfo, error) {
    // 查询数据库
    candidates, err := queryDatabase(xiWuxing)
    if err != nil {
        return nil, err
    }
    
    // 篮选姓氏匹配
    filtered := filterBySurname(candidates, surname)
    
    return filtered, nil
}
```

---

## 辅助函数

```go
func filterBySurname(names []NameInfo, surname string) []NameInfo {
    filtered := []NameInfo{}
    for _, name := range names {
        if matchSurname(name, surname) {
            filtered = append(filtered, name)
        }
    }
    return filtered
}
```

---

## 实现要点

| 要点 | 说明 |
|-----|------|
| 数据库查询 | SQLite 查询汉字库 |
| 姓氏匹配 | 根据姓氏笔画筛选 |

---

## 总结

名字筛选实现包括数据库查询、姓氏匹配、五行筛选。

**核心实现**：FilterNames()
**辅助函数**：filterBySurname()