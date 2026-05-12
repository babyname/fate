# 性能优化

## 性能目标

| 性能指标 | 目标值 |
|---------|-------|
| 单次计算响应时间 | < 1秒 |
| 批量计算（100个） | < 30秒 |

---

## 优化策略

| 策略 | 说明 |
|-----|------|
| 缓存 lunar-go 数据 | 缓存节气时刻、常用计算结果 |
| 预加载常量数据 | 预加载天干五行、地支五行对照表 |
| 避免重复计算 | 一次性获取八字数据 |

---

## 缓存示例

```go
var jieqiCache = make(map[int]time.Time)

func GetJieqi(year int) time.Time {
    if cached, ok := jieqiCache[year]; ok {
        return cached
    }
    jieqi := lunar.GetJieqi(year)
    jieqiCache[year] = jieqi
    return jieqi
}
```

---

## 总结

性能优化策略包括缓存 lunar-go 数据、预加载常量数据、避免重复计算。单次计算响应时间 < 1秒，批量计算（100个） < 30秒。

**优化策略**：缓存、预加载、避免重复计算
**性能目标**：单次 < 1秒，批量 < 30秒