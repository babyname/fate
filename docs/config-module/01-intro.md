# config 模块介绍

## 模块定位

**config** 是 fate 项目的**配置管理模块**，负责：
- 加载和解析 YAML 配置文件
- 管理运行时参数
- 提供配置验证

---

## 核心职责

| 职责 | 说明 |
|-----|------|
| **配置加载** | 加载 YAML 配置文件 |
| **配置验证** | 验证配置参数合法性 |
| **运行时配置** | 管理运行时参数 |

---

## 核心接口

```go
package config

func Load(path string) (*Config, error)
func Validate(cfg *Config) error
func Get() *Config
```

---

## 总结

config 模块负责加载和解析 YAML 配置文件，管理运行时参数。

**核心职责**：配置加载、验证、运行时管理