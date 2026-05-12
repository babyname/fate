# 数据库设计

## 数据库选择

**SQLite**：轻量级、无服务器、适合单机应用

---

## 数据表设计

### hanzi 表

```sql
CREATE TABLE hanzi (
    id INTEGER PRIMARY KEY,
    char TEXT NOT NULL,
    wuxing TEXT NOT NULL,
    bihua INTEGER NOT NULL,
    pinyin TEXT NOT NULL
);
```

---

## 数据表要点

| 表名 | 说明 |
|-----|------|
| hanzi | 汉字库（字、五行、笔画、拼音） |

---

## 总结

数据库使用 SQLite，数据表 hanzi 存储汉字库（字、五行、笔画、拼音）。

**数据库选择**：SQLite
**数据表**：hanzi（汉字库）