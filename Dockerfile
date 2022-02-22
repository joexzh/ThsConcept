FROM golang:1.17.7-alpine3.15 as builder
RUN apk update && apk upgrade
WORKDIR /app
ENV GO111MODULE=on GOPROXY="https://goproxy.io"
COPY . .
RUN go install

FROM alpine:3.15
RUN apk add --no-cache iptables
COPY --from=builder /go/bin /bin
EXPOSE 8088