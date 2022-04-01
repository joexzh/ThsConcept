FROM golang:1.18.0-alpine3.15 as builder
RUN apk update && apk upgrade
WORKDIR /app
ENV GO111MODULE=on GOPROXY="https://goproxy.cn"
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN go build -ldflags '-s -w' -buildvcs=false

FROM alpine:3.15
RUN apk add --no-cache iptables
COPY --from=builder /app /bin
EXPOSE 8088