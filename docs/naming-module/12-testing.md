# 命名模块测试设计

## 测试范围

| 测试类型 | 测试内容 |
|---------|---------|
| 单元测试 | FilterNames()、RateNames()、RecommendNames() |
| 集成测试 | naming + chronos |
| 性能测试 | 单次推荐响应 < 2秒 |

---

## 单元测试示例

```go
func TestRecommendNames(t *testing.T) {
    fateData := &FateData{
        WuxingXiji: &WuxingXijiInfo{
            XiWuxing: []string{"木", "水"},
        },
        Surname: "张",
    }
    
    rated, err := RecommendNames(fateData)
    assert.NoError(t, err)
    assert.NotNil(t, rated)
    assert.LessOrEqual(t, len(rated.Names), 10)
}
```

---

## 总结

测试设计包括单元测试、集成测试、性能测试。

**测试范围**：FilterNames()、RateNames()、RecommendNames()
**测试要点**：已知五行测试、边界测试、空结果测试