FROM golang:alpine AS builder

WORKDIR /app

RUN apk add gcc g++ git

ADD . /app/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -o build/docker/go-proxy main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/build/docker/go-proxy /app/

RUN chmod +x /app/go-proxy

ENTRYPOINT [ "/app/go-proxy"]
