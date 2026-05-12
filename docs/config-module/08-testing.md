# 配置模块测试

## 测试范围

| 测试类型 | 测试内容 |
|---------|---------|
| 单元测试 | Load()、Validate() |
| 集成测试 | 加载 YAML + 验证 |

---

## 测试示例

```go
func TestLoad(t *testing.T) {
    cfg, err := Load("testdata/config.yaml")
    assert.NoError(t, err)
    assert.Equal(t, 1900, cfg.Chonos.DataRangeStart)
}
```

---

## 总结

测试设计包括配置加载测试和验证测试。

**测试要点**：Load()、Validate()、异常路径测试