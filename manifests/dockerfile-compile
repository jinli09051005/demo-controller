FROM golang:1.20-alpine AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o main .



###################
# 接下来创建一个小镜像
###################
FROM scratch
COPY --from=builder /build/main /root/
ENTRYPOINT ["/root/main"]