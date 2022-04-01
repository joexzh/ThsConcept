FROM golang:1.18.0-alpine3.15 as builder
RUN apk update && apk upgrade
WORKDIR /app
ENV GO111MODULE=on GOPROXY="https://goproxy.cn"
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN go install -ldflags '-s -w'

FROM alpine:3.15
RUN apk add --no-cache iptables
COPY --from=builder /go/bin /bin
EXPOSE 8088