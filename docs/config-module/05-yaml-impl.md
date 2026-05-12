# YAML 配置实现

## yaml.go 实现

```go
package config

func LoadYAML(path string) (*Config, error) {
    buf, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("读取配置: %w", err)
    }
    
    cfg := &Config{}
    if err := yaml.Unmarshal(buf, cfg); err != nil {
        return nil, fmt.Errorf("解析配置: %w", err)
    }
    
    return cfg, Validate(cfg)
}
```

---

## 实现要点

| 要点 | 说明 |
|-----|------|
| 读取文件 | os.ReadFile |
| YAML解析 | yaml.Unmarshal |
| 配置验证 | Validate() |

---

## 总结

YAML 配置实现使用 os.ReadFile 读取文件，yaml.Unmarshal 解析配置。

**核心实现**：LoadYAML()
**解析流程**：读取 → 解析 → 验证