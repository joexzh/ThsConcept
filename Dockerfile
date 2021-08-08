FROM golang:1.16.7-alpine3.14 as builder
ENV GO111MODULE=on
RUN apk update && \
    apk upgrade && \
    apk add gcc libc-dev linux-headers
WORKDIR /app
COPY . .
RUN go build && go install

FROM alpine:3.14
RUN apk add --no-cache iptables
COPY --from=builder /go/bin /bin
EXPOSE 8088