# naming 模块介绍

## 模块定位

**naming** 是 fate 项目的**起名推荐模块**，负责：
- 根据五行喜忌筛选名字
- 综合评分名字（五行、笔画、音韵）
- 推荐最佳名字列表

---

## 核心职责

| 职责 | 说明 | 重要性 |
|-----|------|-------|
| **名字筛选** | 根据五行喜忌筛选名字 | 核心 |
| **名字评分** | 综合评分名字 | 核心 |
| **名字推荐** | 推荐最佳名字列表 | 核心 |

---

## 核心接口

```go
package naming

func FilterNames(xiWuxing []string) ([]NameInfo, error)
func RateNames(names []NameInfo, fateData *FateData) (*RatedNames, error)
func RecommendNames(fateData *FateData) (*RatedNames, error)
```

---

## 模块文件结构

```
naming/
├── naming.go        # 核心入口
├── filter.go        # 筛选名字
├── rate.go          # 评分名字
├── database.go      # 数据库查询
├── types.go         # 类型定义
└── errors.go        # 错误定义
```

---

## 总结

naming 模块负责名字筛选、评分、推荐。核心接口包括 FilterNames()、RateNames()、RecommendNames()。

**核心职责**：筛选、评分、推荐
**核心接口**：FilterNames()、RateNames()、RecommendNames()