# 五行喜忌格式化

## FormatWuxingXiji() 实现

```go
package analysis

func FormatWuxingXiji(fateData *FateData) string {
    xiWuxing := strings.Join(fateData.WuxingXiji.XiWuxing, ", ")
    jiWuxing := strings.Join(fateData.WuxingXiji.JiWuxing, ", ")
    
    return fmt.Sprintf("喜用五行: %s\n忌神五行: %s", xiWuxing, jiWuxing)
}
```

---

## 格式说明

| 内容 | 说明 |
|-----|------|
| 喜用五行 | 五行列表（用逗号分隔） |
| 忌神五行 | 五行列表（用逗号分隔） |

---

## 总结

五行喜忌格式化输出喜用五行和忌神五行，用逗号分隔。

**核心接口**：FormatWuxingXiji()
**输出格式**：喜用五行: xxx, 忌神五行: xxx