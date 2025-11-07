# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o miniblog ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 复制编译好的二进制文件
COPY --from=builder /app/miniblog .

# 复制配置文件
COPY --from=builder /app/config ./config

# 创建日志目录
RUN mkdir -p /root/logs

EXPOSE 8080 9090

CMD ["./miniblog"]

