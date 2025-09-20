FROM golang:1.23 AS builder

# 设置工作目录
WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./

# 复制源代码
COPY . .

# 下载依赖
RUN go mod tidy

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o glimgate main.go

# 运行阶段
FROM alpine:latest

# 安装ca证书
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/glimgate .
COPY --from=builder /app/config ./config

# 暴露端口
EXPOSE 20401

# 运行应用
CMD ["./glimgate"]