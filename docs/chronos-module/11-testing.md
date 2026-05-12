# 测试设计

## 测试范围

| 测试类型 | 测试内容 |
|---------|---------|
| 单元测试 | GetFateData()、CalculateBazi()、CalculateWuxingXiji() |
| 集成测试 | chronos + naming、chronos + analysis |
| 性能测试 | 单次计算响应时间 < 1秒 |

---

## 单元测试示例

```go
func TestGetFateData(t *testing.T) {
    input := &FateInput{
        BirthDate: time.Date(2020, 1, 23, 11, 31, 0, 0, time.UTC),
        Gender:    1,
    }
    
    fateData, err := GetFateData(input)
    assert.NoError(t, err)
    assert.NotNil(t, fateData)
    assert.Equal(t, "甲子", fateData.Bazi.Sizhu[0])
}
```

---

## 测试要点

| 要点 | 说明 |
|-----|------|
| 已知案例 | 使用已知八字案例验证正确性 |
| 边界测试 | 测试1900年、2100年边界日期 |
| 错误测试 | 测试无效输入、错误日期 |

---

## 总结

测试设计包括单元测试、集成测试、性能测试。使用已知八字案例验证正确性，测试边界日期和错误输入。

**测试范围**：GetFateData()、CalculateBazi()、CalculateWuxingXiji()
**测试要点**：已知案例、边界测试、错误测试