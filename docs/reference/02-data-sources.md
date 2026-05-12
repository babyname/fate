# 数据源

## 依赖库

| 库 | 版本 | 用途 |
|-----|------|------|
| lunar-go | v1.3.0 | 农历日历计算 |
| go-sqlite3 | latest | SQLite 驱动 |
| yaml.v3 | latest | YAML 解析 |

---

## 数据文件

| 文件 | 说明 |
|-----|------|
| hanzi.db | 汉字库（SQLite） |
| config.yaml | 配置文件 |

---

## 总结

数据源包括 lunar-go 日历库、SQLite 汉字库、YAML 配置文件。

**核心依赖**：lunar-go、go-sqlite3、yaml.v3