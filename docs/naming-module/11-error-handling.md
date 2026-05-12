# 命名模块错误处理

## errors.go

```go
package naming

type NameError struct {
    Code    int
    Message string
}

func (e *NameError) Error() string {
    return fmt.Sprintf("NameError %d: %s", e.Code, e.Message)
}

const (
    ErrCodeNoMatch   = 4001
    ErrCodeDBError   = 4002
    ErrCodeEmptyList = 4003
)
```

---

## 错误处理要点

| 要点 | 说明 |
|-----|------|
| 统一错误类型 | NameError |
| 错误代码 | 4000系列（命名模块） |
| 降级处理 | 无匹配名字时返回空列表 |

---

## 总结

错误处理使用 NameError 统一错误格式，错误代码4000系列。

**核心类型**：NameError
**错误代码**：ErrCodeNoMatch、ErrCodeDBError、ErrCodeEmptyList