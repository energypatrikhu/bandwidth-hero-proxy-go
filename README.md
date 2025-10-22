# Bandwidth Hero Proxy (Go)

Image compression proxy server that reduces bandwidth usage by compressing images on-the-fly.

> Typescript version: [bandwidth-hero-proxy](https://github.com/energypatrikhu/bandwidth-hero-proxy) (no longer maintained)

## Features

- Supports WebP and JPEG compression
- Automatic format selection for best compression
- Optional greyscale conversion
- Configurable quality levels
- Animated GIF support
- Request retry logic and redirect handling

## Quick Start

- **Option 1: Use prebuilt package**
    ```bash
    docker run --publish 80:80 ghcr.io/energypatrikhu/bandwidth-hero-proxy-go:latest
    ```

- **Option 2: Build it yourself**
    1. Clone the repository:
    ```bash
    git clone energypatrikhu/bandwidth-hero-proxy-go
    cd bandwidth-hero-proxy-go
    ```

    2. Build and compose up:
    ```bash
    docker-compose up --build
    ```

## Development
  1. Clone the repository, then navigate into the directory
  ```bash
  git clone energypatrikhu/bandwidth-hero-proxy-go
  cd bandwidth-hero-proxy-go
  ```
  
  2. Install `vipsgen` and generate bindings
  ```bash
  # Install libvips (Ubuntu/Debian: libvips-dev, macOS: brew install vips)
  go install github.com/cshum/vipsgen/cmd/vipsgen@latest

  # Generate vips bindings
  vipsgen -out ./vips
  ```

  3. Download dependencies, build and run
  ```bash
  go mod download
  go build -o bandwidth-hero-proxy main.go
  ./bandwidth-hero-proxy
  ```

## Usage
> Note: It is recommended to **place** the `url` query to the **end of the request** and **url encode** it, to prevent the query strings to mix up and get placed into the wrong request.

```
http://your-proxy-server/?quality=<QUALITY>&jpg=<0|1>&greyscale=<0|1>&url=<IMAGE_URL>
```

**Parameters:**
- `quality`: Compression quality 1-100 (default: 80)
- `jpg`: Use JPEG instead of WebP (default: 0)
- `greyscale`: Convert to greyscale (default: 0)
- `url` (required): Image URL to compress

**Examples:**
```bash
# Default WebP compression
http://localhost/?url=https://example.com/image.jpg

# JPEG with 60% quality
http://localhost/?jpg=1&quality=60&url=https://example.com/image.png
```

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `BHP_PORT` | `80` | Server port |
| `BHP_MAX_CONCURRENCY` | Number of CPU cores | Max concurrent tasks |
| `BHP_FORCE_FORMAT` | `false` | Force selected format, even if the output is bigger |
| `BHP_AUTO_DECREMENT_QUALITY` | `false` | Auto decrement quality if output is larger than input |
| `BHP_USE_BEST_COMPRESSION_FORMAT` | `false` | Automatically choose WebP or JPEG based on compression ratio |
| `BHP_EXTERNAL_REQUEST_TIMEOUT` | `60s` | External request timeout |
| `BHP_EXTERNAL_REQUEST_RETRIES` | `5` | Number of retries for external requests |
| `BHP_EXTERNAL_REQUEST_REDIRECTS` | `10` | Maximum redirects for external requests |
| `BHP_EXTERNAL_REQUEST_OMIT_HEADERS` | `[]` | Headers to omit from external requests |

Example:
```bash
export BHP_PORT=8080
export BHP_USE_BEST_COMPRESSION_FORMAT=true
./bandwidth-hero-proxy
```

## Response Headers

- `X-Original-Size`: Original image size in bytes
- `X-Compressed-Size`: Compressed image size in bytes
- `X-Size-Saved`: Bytes saved through compression

## Behavior

- Defaults to WebP format, use `jpg=1` for JPEG
- Redirects to original URL if compression fails or doesn't reduce size
- Preserves animation in GIFs meanwhile it compresses each frame
- Automatically retries failed requests

## Troubleshooting

- **Build issues**: Install libvips dev headers and ensure `CGO_ENABLED=1`
- **Images not compressing**: Check source URL accessibility and image format support
- **URL not provided**: Ensure `url` query is included in the request, if still gives an error, try URL encoding the URL
- **High memory usage**: Reduce `BHP_MAX_CONCURRENCY`
- **Timeouts**: Increase `BHP_EXTERNAL_REQUEST_TIMEOUT`
