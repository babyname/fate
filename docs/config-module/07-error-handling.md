# 配置模块错误处理

## errors.go

```go
package config

type ConfigError struct {
    Code    int
    Message string
    Path    string
}

func (e *ConfigError) Error() string {
    return fmt.Sprintf("ConfigError[%s] %d: %s", e.Path, e.Code, e.Message)
}

const (
    ErrCodeFileNotFound = 6001
    ErrCodeParseFailed  = 6002
    ErrCodeValidation   = 6003
)
```

---

## 总结

错误处理使用 ConfigError，包含路径信息，错误代码6000系列。

**核心类型**：ConfigError
**错误代码**：ErrCodeFileNotFound、ErrCodeParseFailed、ErrCodeValidation