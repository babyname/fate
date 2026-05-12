# Docker 部署

## Dockerfile

```dockerfile
FROM golang:1.21 AS build
WORKDIR /app
COPY . .
RUN go build -o fate ./cmd/fate

FROM alpine:3.19
COPY --from=build /app/fate /usr/local/bin/fate
COPY config.yaml /etc/fate/config.yaml
CMD ["fate", "server", "--config", "/etc/fate/config.yaml"]
```

---

## docker-compose.yml

```yaml
version: "3"
services:
  fate:
    build: .
    volumes:
      - ./data:/data
    ports:
      - "8080:8080"
```

---

## 总结

Docker 部署使用多阶段构建减小镜像大小，支持 docker-compose。

**镜像方式**：多阶段构建
**部署工具**：docker-compose