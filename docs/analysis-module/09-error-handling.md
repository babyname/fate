# 分析模块错误处理

## errors.go

```go
package analysis

type FormatError struct {
    Code    int
    Message string
}

func (e *FormatError) Error() string {
    return fmt.Sprintf("FormatError %d: %s", e.Code, e.Message)
}

const (
    ErrCodeFormatFail = 5001
    ErrCodeEmptyData  = 5002
)
```

---

## 总结

错误处理使用 FormatError 统一格式，错误代码5000系列。

**核心类型**：FormatError
**错误代码**：ErrCodeFormatFail、ErrCodeEmptyData