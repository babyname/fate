# 部署方案

## 部署配置

| 参数 | 值 |
|-----|-----|
| 语言 | Go 1.21+ |
| 构建 | go build |
| 平台 | Windows/Linux/macOS |

---

## 部署步骤

```
# 1. 克隆代码
git clone https://github.com/godcong/fate.git

# 2. 构建
go build -o fate ./cmd/fate

# 3. 配置
编辑 config.yaml

# 4. 运行
./fate server --config config.yaml
```

---

## 总结

部署方案使用 go build 构建，支持多平台。

**部署方式**：go build
**支持平台**：Windows/Linux/macOS