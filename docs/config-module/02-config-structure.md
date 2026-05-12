# 配置结构设计

## Config 结构

```go
package config

type Config struct {
    Chonos   ChronosConfig   `yaml:"chronos"`
    Naming   NamingConfig    `yaml:"naming"`
    Analysis AnalysisConfig  `yaml:"analysis"`
}

type ChronosConfig struct {
    DataRangeStart int    `yaml:"dataRangeStart"`
    DataRangeEnd   int    `yaml:"dataRangeEnd"`
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

配置结构包括 ChronosConfig、NamingConfig、AnalysisConfig，使用 YAML 标签。

**核心类型**：Config、各模块 Config