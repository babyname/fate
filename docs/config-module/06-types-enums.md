# 配置类型定义

## types.go

```go
package config

type Config struct {
    Chonos   ChronosConfig  `yaml:"chronos"`
    Naming   NamingConfig   `yaml:"naming"`
    Analysis AnalysisConfig `yaml:"analysis"`
}

type ChronosConfig struct {
    DataRangeStart int `yaml:"dataRangeStart"`
    DataRangeEnd   int `yaml:"dataRangeEnd"`
}

type NamingConfig struct {
    DBPath   string `yaml:"dbPath"`
    MaxNames int    `yaml:"maxNames"`
}

type AnalysisConfig struct {
    OutputFormat string `yaml:"outputFormat"`
}
```

---

## 总结

类型定义包括 Config 及其子配置类型，使用 YAML 标签。

**核心类型**：Config、ChronosConfig、NamingConfig、AnalysisConfig