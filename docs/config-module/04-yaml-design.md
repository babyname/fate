# YAML 配置设计

## config.yaml 示例

```yaml
chronos:
  dataRangeStart: 1900
  dataRangeEnd: 2100

naming:
  dbPath: "data/hanzi.db"
  maxNames: 10

analysis:
  outputFormat: "json"
```

---

## 配置项说明

| 配置项 | 说明 | 默认值 |
|-------|------|-------|
| chronos.dataRangeStart | 数据起始年份 | 1900 |
| naming.dbPath | 数据库路径 | data/hanzi.db |
| analysis.outputFormat | 输出格式 | json |

---

## 总结

YAML 配置设计包含 chronos、naming、analysis 配置项。

**配置文件**：config.yaml
**配置项**：日期范围、数据库路径、输出格式