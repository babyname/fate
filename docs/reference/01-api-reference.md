# API 参考

## 核心 API

| API | 说明 |
|-----|------|
| `GetFateData(input)` | 获取八字数据 |
| `FilterNames(xiWuxing, surname)` | 筛选名字 |
| `RateNames(names, fateData)` | 评分名字 |
| `RecommendNames(fateData)` | 推荐名字 |
| `FormatReport(fateData, names)` | 格式化报告 |
| `Load(path)` | 加载配置 |

---

## 数据结构

| 类型 | 说明 |
|-----|------|
| `FateInput` | 输入参数（日期、姓氏） |
| `FateData` | 八字数据（四柱、喜忌） |
| `NameInfo` | 名字信息（字、五行、笔画） |
| `RatedNames` | 评分名字列表 |
| `Config` | 配置结构 |

---

## 总结

API 参考包括 GetFateData()、RecommendNames()、FormatReport() 等核心 API。

**核心 API**：6个
**数据结构**：5个