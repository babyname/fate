# 数据库实现

## database.go 实现

```go
package naming

import "database/sql"
import _ "github.com/mattn/go-sqlite3"

func queryDatabase(xiWuxing []string) ([]NameInfo, error) {
    db, err := sql.Open("sqlite3", "hanzi.db")
    if err != nil {
        return nil, err
    }
    defer db.Close()
    
    query := "SELECT char, wuxing, bihua FROM hanzi WHERE wuxing IN (?, ?)"
    rows, err := db.Query(query, xiWuxing[0], xiWuxing[1])
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    names := []NameInfo{}
    for rows.Next() {
        var name NameInfo
        rows.Scan(&name.Char, &name.Wuxing, &name.Bihua)
        names = append(names, name)
    }
    
    return names, nil
}
```

---

## 实现要点

| 要点 | 说明 |
|-----|------|
| SQLite 连接 | sql.Open("sqlite3", "hanzi.db") |
| 查询五行 | WHERE wuxing IN (?, ?) |
| 资源管理 | defer db.Close(), defer rows.Close() |

---

## 总结

数据库实现使用 SQLite，查询五行匹配的汉字。

**核心实现**：queryDatabase()
**查询方式**：SQLite + WHERE 条件