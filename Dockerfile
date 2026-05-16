FROM golang:1.23-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN GOPROXY=https://goproxy.cn,direct go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o free_proxy_pool .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/free_proxy_pool .
COPY --from=builder /build/config.yaml .

EXPOSE 5555 10888

ENTRYPOINT ["./free_proxy_pool", "--config", "./config.yaml"]
