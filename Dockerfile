FROM golang:alpine AS builder

ENV CGO_ENABLED=1

WORKDIR /app

RUN apk add --no-cache build-base vips-dev pkgconfig

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bandwidth-hero-proxy -x -ldflags="-s -w" main.go

FROM alpine:latest AS runtime

ENV MALLOC_ARENA_MAX=2

RUN apk add --no-cache vips

COPY --from=builder /bandwidth-hero-proxy /bandwidth-hero-proxy

ENTRYPOINT [ "/bandwidth-hero-proxy" ]
