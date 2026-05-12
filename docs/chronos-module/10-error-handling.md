# 错误处理

## errors.go 定义

```go
package chronos

// FateError 统一错误类型
type FateError struct {
    Code    int
    Message string
    Module  string
}

func (e *FateError) Error() string {
    return fmt.Sprintf("[%s] Error %d: %s", e.Module, e.Code, e.Message)
}

// 错误代码
const (
    ErrCodeInputInvalid    = 1001
    ErrCodeDateRange       = 1002
    ErrCodeGenderInvalid   = 1003
    ErrCodeCalculateBazi   = 2001
    ErrCodeCalculateWuxing = 2002
)
```

---

## 错误处理示例

```go
// 输入验证错误
if input.BirthDate.Year() < 1900 {
    return nil, &FateError{
        Code:    ErrCodeDateRange,
        Message: "日期范围错误：1900-2100年",
        Module:  "chronos",
    }
}

// 降级处理
wuxing, err := CalculateWuxingXiji(bazi)
if err != nil {
    wuxing = defaultWuxingXiji(bazi)  // 使用默认值
}
```

---

## 错误处理要点

| 要点 | 说明 |
|-----|------|
| 统一错误类型 | FateError 统一错误格式 |
| 错误代码 | 模块化错误代码（1000-4000） |
| 降级处理 | 计算失败使用默认值 |

---

## 总结

错误处理使用统一错误类型 FateError，错误代码按模块划分（1000-输入、2000-计算、3000-推荐、4000-输出），降级处理使用默认值。

**核心类型**：FateError
**错误代码**：ErrCodeInputInvalid、ErrCodeDateRange等
**降级处理**：defaultWuxingXiji()