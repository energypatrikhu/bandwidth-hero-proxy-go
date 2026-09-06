FROM golang:alpine AS builder

WORKDIR /src

RUN apk add --no-cache \
  build-base \
  pkgconfig \
  vips-dev \
  && rm -rf /var/lib/apt/lists/*

RUN mkdir -p third_party && \
  go install github.com/cshum/vipsgen/cmd/vipsgen@latest && \
  vipsgen -out ./third_party/vips

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal

RUN CGO_ENABLED=1 GOOS=linux \
  go build \
  -v \
  -ldflags="-s -w" \
  -trimpath \
  -o /bandwidth-hero-proxy \
  ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache \
  vips \
  && rm -rf /var/lib/apt/lists/*

COPY --from=builder /bandwidth-hero-proxy /bandwidth-hero-proxy

ENTRYPOINT ["/bandwidth-hero-proxy"]
