FROM golang:alpine as builder
ENV GO111MODULE=on
RUN apk update && \
    apk upgrade && \
    apk add gcc libc-dev linux-headers
WORKDIR /app \
COPY . .
RUN go build && go install

FROM alpine
RUN apk add --no-cache iptables
COPY --from=builder /go/bin /bin
EXPOSE 8088