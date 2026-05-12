# 配置实现

## config.go 实现

```go
package config

var globalConfig *Config

func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    
    cfg := &Config{}
    err = yaml.Unmarshal(data, cfg)
    if err != nil {
        return nil, err
    }
    
    globalConfig = cfg
    return cfg, nil
}

func Get() *Config {
    return globalConfig
}
```

---

## 总结

配置实现使用 YAML 加载配置，存储到全局变量中。

**核心实现**：Load()、Get()