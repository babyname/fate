# 数据表设计

## constants.go 数据表

```go
package chronos

// 天干五行对照表
var ganWuxingMap = map[string]string{
    "甲": "木", "乙": "木",
    "丙": "火", "丁": "火",
    "戊": "土", "己": "土",
    "庚": "金", "辛": "金",
    "壬": "水", "癸": "水",
}

// 地支五行对照表
var zhiWuxingMap = map[string]string{
    "子": "水", "丑": "土",
    "寅": "木", "卯": "木",
    "辰": "土", "巳": "火",
    "午": "火", "未": "土",
    "申": "金", "酉": "金",
    "戌": "土", "亥": "水",
}

// 藏干对照表
var cangganMap = map[string][]string{
    "子": {"癸"}, "丑": {"己", "辛", "癸"},
    "寅": {"甲", "丙", "戊"}, "卯": {"乙"},
    "辰": {"戊", "乙", "癸"}, "巳": {"丙", "戊", "庚"},
    "午": {"丁", "己"}, "未": {"己", "丁", "乙"},
    "申": {"庚", "壬", "戊"}, "酉": {"辛"},
    "戌": {"戊", "辛", "丁"}, "亥": {"壬", "甲"},
}

// 五行生克关系表
var wuxingShengkeMap = map[string]map[string]string{
    "木": {"shengwo": "水", "kewo": "金", "xiewo": "火", "haowo": "土"},
    "火": {"shengwo": "木", "kewo": "水", "xiewo": "土", "haowo": "金"},
    "土": {"shengwo": "火", "kewo": "木", "xiewo": "金", "haowo": "水"},
    "金": {"shengwo": "土", "kewo": "火", "xiewo": "水", "haowo": "木"},
    "水": {"shengwo": "金", "kewo": "土", "xiewo": "木", "haowo": "火"},
}
```

---

## 数据表要点

| 数据表 | 说明 |
|-----|------|
| ganWuxingMap | 天干五行对照 |
| zhiWuxingMap | 地支五行对照 |
| cangganMap | 藏干对照 |
| wuxingShengkeMap | 五行生克关系 |

---

## 总结

数据表设计包括天干五行、地支五行、藏干、五行生克关系对照表，预加载在 constants.go 中，提高性能。

**核心数据表**：ganWuxingMap、zhiWuxingMap、cangganMap、wuxingShengkeMap
**使用方式**：预加载常量，快速查询