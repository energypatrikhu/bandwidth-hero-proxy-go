FROM golang:trixie AS builder

WORKDIR /src

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
  ca-certificates \
  libvips-dev \
  && rm -rf /var/lib/apt/lists/*

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
  -trimpath \
  -o /bandwidth-hero-proxy \
  ./cmd/main.go

FROM debian:trixie-slim

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
  ca-certificates \
  libjemalloc2 \
  libvips \
  && rm -rf /var/lib/apt/lists/*

ENV LD_PRELOAD=/usr/lib/x86_64-linux-gnu/libjemalloc.so.2
ENV GOMEMLIMIT=512MiB

COPY --from=builder /bandwidth-hero-proxy /bandwidth-hero-proxy

ENTRYPOINT ["/bandwidth-hero-proxy"]
