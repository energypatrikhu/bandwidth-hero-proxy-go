FROM golang:latest AS builder

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
  apt-get install -y \
  libvips-dev && \
  rm -rf /var/lib/apt/lists/*

WORKDIR /src

RUN mkdir -p third_party && \
  go install github.com/cshum/vipsgen/cmd/vipsgen@latest && \
  vipsgen -out ./third_party/vips

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal

ENV CGO_ENABLED=1

RUN go build \
  -v \
  -tags=vips \
  -ldflags="-s -w" \
  -o /bandwidth-hero-proxy \
  ./cmd/main.go

FROM debian:stable-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
  apt-get install -y \
  libvips && \
  rm -rf /var/lib/apt/lists/*

COPY --from=builder /bandwidth-hero-proxy /bandwidth-hero-proxy

ENTRYPOINT ["/bandwidth-hero-proxy"]
