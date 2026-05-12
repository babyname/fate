# CI/CD 配置

## CI/CD 方案

**使用 GitHub Actions** 作为 CI/CD 工具。

---

## ci.yml

```yaml
name: CI
on: [push, pull_request]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: "1.21"
      - run: go build ./...
      - run: go test ./...
      - run: go vet ./...
```

---

## 总结

CI/CD 使用 GitHub Actions，包括构建、测试和静态分析。

**CI/CD 工具**：GitHub Actions
**检查步骤**：build、test、vet