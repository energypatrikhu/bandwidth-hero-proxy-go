# Bandwidth Hero Proxy (Go)

Image compression proxy server that reduces bandwidth usage by compressing images on-the-fly.

> Typescript version: [bandwidth-hero-proxy](https://github.com/energypatrikhu/bandwidth-hero-proxy) (no longer maintained)

## Features

- Supports WebP and JPEG compression
- Automatic format selection for best compression
- Optional grayscale conversion
- Configurable quality levels
- Animated GIF support
- Request retry logic and redirect handling
- FlareSolverr support for Cloudflare anti-bot challenges

## Quick Start

- **Option 1: Use prebuilt package**
  - a. Pull and run the Docker image:
  ```bash
  docker run --publish 8080:80 docker.io/energypatrikhu/bandwidth-hero-proxy-go:latest
  ```

  - b. Download compose file and run with Docker Compose:
  ```bash
  curl -O https://raw.githubusercontent.com/energypatrikhu/bandwidth-hero-proxy-go/main/docker-compose.yml
  docker compose up -d
  ```

  - c. Copy docker compose manually into your existing compose file:
  ```yaml
  services:
    bandwidth-hero-proxy:
      container_name: bandwidth-hero-proxy
      image: docker.io/energypatrikhu/bandwidth-hero-proxy-go
      network_mode: bridge
      # environment: # optional environment variables
      #   BHP_PORT: 80
      #   BHP_FLARESOLVERR_URL: "http://flaresolverr:8191"
      #   BHP_VIPS_MAX_CONCURRENCY: 4 # default: number of CPU cores
      #   BHP_FORCE_FORMAT: false
      #   BHP_AUTO_DECREMENT_QUALITY: false
      #   BHP_USE_BEST_COMPRESSION_FORMAT: true
      #   BHP_DISABLE_ANIMATED_IMAGES: false
      #   BHP_EXTERNAL_REQUEST_TIMEOUT: "60s"
      #   BHP_EXTERNAL_REQUEST_RETRIES: 5
      #   BHP_EXTERNAL_REQUEST_REDIRECTS: 10
      #   BHP_EXTERNAL_REQUEST_OMIT_HEADERS: ""
      #   # Webp options
      #   BHP_WEBP_LOSSLESS: false
      #   BHP_WEBP_EFFORT: 6
      #   BHP_WEBP_SMART_SUBSAMPLE: true
      #   BHP_WEBP_SMART_DEBLOCK: true
      #   BHP_WEBP_PASSES: 10
      #   # Jpeg options
      #   BHP_JPEG_OPTIMIZE_CODING: true
      #   BHP_JPEG_OPTIMIZE_SCANS: true
      #   BHP_JPEG_INTERLACE: false
      #   BHP_JPEG_TRELLIS_QUANT: true
      #   BHP_JPEG_OVERSHOOT_DERINGING: true
      #   BHP_JPEG_QUANT_TABLE: 3
      ports:
        - 8080:80
  ```

- **Option 2: Build it yourself**
  1. Clone the repository:
  ```bash
  git clone https://github.com/energypatrikhu/bandwidth-hero-proxy-go
  cd bandwidth-hero-proxy-go
  ```

  2. Build and compose up:
  ```bash
  docker compose up --build
  ```

## Development

1. Clone the repository, then navigate into the directory

```bash
git clone https://github.com/energypatrikhu/bandwidth-hero-proxy-go
cd bandwidth-hero-proxy-go
```

2. Install `vipsgen` and generate bindings

```bash
# Install libvips (Ubuntu/Debian: libvips-dev, macOS: brew install vips)
go install github.com/cshum/vipsgen/cmd/vipsgen@latest

# Generate vips bindings
vipsgen -out ./third_party/vips
```

3. Download dependencies, build and run

```bash
go mod download
go build -o bandwidth-hero-proxy main.go
./bandwidth-hero-proxy
```

## Usage

> Note: It is recommended to **place** the `url` query to the **end of the request** and **url encode** it, to prevent the query strings being mixed up and get placed into the wrong request.

```
http://your-proxy-server/?quality=<QUALITY>&jpg=<0|1>&grayscale=<0|1>&url=<IMAGE_URL>
```

**Parameters:**

- `quality`: Compression quality 1-100 (default: 80)
- `jpg`: Use JPEG instead of WebP (default: 0)
- `grayscale`: Convert to grayscale (default: 0)
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

| Variable                            | Default             | Description                                                     |
| ----------------------------------- | ------------------- | --------------------------------------------------------------- |
| `BHP_PORT`                          | `80`                | Server port                                                     |
| `BHP_FLARESOLVERR_URL`              | `""`                | URL of the FlareSolverr instance to use for anti-bot challenges |
| `BHP_VIPS_MAX_CONCURRENCY`          | Number of CPU cores | Vips max concurrent tasks                                       |
| `BHP_FORCE_FORMAT`                  | `false`             | Force selected format, even if the output is bigger             |
| `BHP_AUTO_DECREMENT_QUALITY`        | `false`             | Auto decrement quality if output is larger than input           |
| `BHP_USE_BEST_COMPRESSION_FORMAT`   | `false`             | Automatically choose WebP or JPEG based on compression ratio    |
| `BHP_DISABLE_ANIMATED_IMAGES`       | `false`             | Disable compression of animated images                          |
| `BHP_EXTERNAL_REQUEST_TIMEOUT`      | `60s`               | External request timeout                                        |
| `BHP_EXTERNAL_REQUEST_RETRIES`      | `5`                 | Number of retries for external requests                         |
| `BHP_EXTERNAL_REQUEST_REDIRECTS`    | `10`                | Maximum redirects for external requests                         |
| `BHP_EXTERNAL_REQUEST_OMIT_HEADERS` | `[]`                | Headers to omit from external requests                          |
| `BHP_WEBP_LOSSLESS`                 | `false`             | Enable lossless compression                                     |
| `BHP_WEBP_EFFORT`                   | `4`                 | Level of CPU effort to reduce file size                         |
| `BHP_WEBP_SMART_SUBSAMPLE`          | `true`              | Enable high quality chroma subsampling                          |
| `BHP_WEBP_SMART_DEBLOCK`            | `true`              | Enable auto-adjusting of the deblocking filter                  |
| `BHP_WEBP_PASSES`                   | `2`                 | Number of entropy-analysis passes (in [1..10])                  |
| `BHP_JPEG_OPTIMIZE_CODING`          | `true`              | Compute optimal Huffman coding tables                           |
| `BHP_JPEG_OPTIMIZE_SCANS`           | `true`              | Split spectrum of DCT coefficients into separate scans          |
| `BHP_JPEG_INTERLACE`                | `false`             | Generate an interlaced (progressive) jpeg                       |
| `BHP_JPEG_TRELLIS_QUANT`            | `true`              | Apply trellis quantisation to each 8x8 block                    |
| `BHP_JPEG_OVERSHOOT_DERINGING`      | `true`              | Apply overshooting to samples with extreme values               |
| `BHP_JPEG_QUANT_TABLE`              | `3`                 | Use predefined quantization table with given index              |

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
- Uses FlareSolverr to solve Cloudflare anti-bot challenges, if configured
- Won't compress animated images, if configured

## Troubleshooting

- **Build issues**: Install libvips dev headers and ensure `CGO_ENABLED=1`
- **Images not compressing**: Check source URL accessibility and image format support
- **URL not provided**: Ensure `url` query is included in the request, if still gives an error, try URL encoding the URL
- **High memory usage**: Reduce `BHP_MAX_CONCURRENCY`
- **Timeouts**: Increase `BHP_EXTERNAL_REQUEST_TIMEOUT`
- **Cloudflare anti-bot challenges**: If using FlareSolverr, ensure `BHP_FLARESOLVERR_URL` is set correctly and the FlareSolverr instance is reachable
