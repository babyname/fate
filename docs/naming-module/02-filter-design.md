# 名字筛选设计

## 筛选原理

根据五行喜忌筛选名字，只保留五行符合的名字。

---

## 筛选流程

```
查询数据库 → 匹配五行 → 筛选姓氏 → 返回候选名字
```

---

## FilterNames() 接口

```go
package naming

func FilterNames(xiWuxing []string, surname string) ([]NameInfo, error) {
    // 1. 查询数据库获取候选名字
    candidates := queryDatabase(xiWuxing)
    
    // 2. 筛选姓氏匹配的名字
    filtered := filterBySurname(candidates, surname)
    
    // 3. 返回
    return filtered, nil
}
```

---

## 篩选要点

| 要点 | 说明 |
|-----|------|
| 五行匹配 | 只保留喜用五行名字 |
| 姓氏筛选 | 根据姓氏笔画筛选 |
| 数据库查询 | SQLite 查询汉字库 |

---

## 总结

名字筛选根据五行喜忌筛选名字，通过数据库查询、五行匹配、姓氏筛选，返回候选名字列表。

**核心接口**：FilterNames()
**筛选流程**：查询 → 匹配 → 篮选 → 返回