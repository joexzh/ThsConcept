FROM golang:1.17.7-alpine3.15 as builder
RUN apk update && apk upgrade
WORKDIR /app
COPY . .
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"
RUN go install

FROM alpine:3.15
RUN apk add --no-cache iptables
COPY --from=builder /go/bin /bin
COPY --from=builder /app/tmpl/ /bin/tmpl/
EXPOSE 8088