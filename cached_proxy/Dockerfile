FROM golang AS builder
# 设置工作目录
WORKDIR /app
# 复制 go.mod 和 go.sum 文件
COPY ./ ./

# build the apps
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s"  -o cached_proxy ./


# 使用轻量级镜像
FROM scratch

# 设置工作目录
WORKDIR /app

# 将编译好的静态文件复制到镜像中
COPY --from=builder /app/cached_proxy /app/cached_proxy

# 暴露应用的端口（根据你的应用而定，例如 8080）
EXPOSE 8000

# 设置可执行文件作为启动命令
ENTRYPOINT ["./cached_proxy"]
