FROM golang:alpine AS builder

WORKDIR /tmp/app-source

# Install build dependencies and all vips optional dependencies
RUN apk add --no-cache \
  build-base \
  pkgconfig \
  git \
  cmake \
  meson \
  ninja \
  nasm \
  autoconf \
  automake \
  libtool \
  gettext-dev \
  glib-dev \
  expat-dev \
  tiff-dev \
  libpng-dev \
  zlib-dev \
  libwebp-dev \
  giflib-dev \
  librsvg-dev \
  poppler-dev \
  cairo-dev \
  pango-dev \
  fftw-dev \
  orc-dev \
  lcms2-dev \
  libexif-dev \
  libgsf-dev \
  openexr-dev \
  libraw-dev \
  libheif-dev \
  libimagequant-dev \
  openjpeg-dev \
  libarchive-dev \
  cfitsio-dev \
  libspng-dev \
  libjxl-dev \
  cgif-dev \
  highway-dev \
  libjpeg-turbo-dev

# Build vips from source
RUN VIPS_VERSION=$(wget -qO- "https://api.github.com/repos/libvips/libvips/releases/latest" | grep -o '"tag_name": "v[^"]*"' | cut -d'"' -f4 | sed 's/^v//') && \
  mkdir -p /tmp/vips-source && \
  wget -O /tmp/vips-source/vips.tar.xz "https://github.com/libvips/libvips/releases/latest/download/vips-${VIPS_VERSION}.tar.xz" && \
  cd /tmp/vips-source && \
  tar -xf vips.tar.xz && \
  cd vips-${VIPS_VERSION} && \
  echo "Checking JPEG libraries before vips build:" && \
  pkg-config --exists libjpeg && echo "libjpeg found" || echo "libjpeg not found" && \
  pkg-config --exists libturbojpeg && echo "libturbojpeg found" || echo "libturbojpeg not found" && \
  PKG_CONFIG_PATH="/usr/lib/pkgconfig:/usr/local/lib/pkgconfig" \
  meson setup build \
    --buildtype=release \
    --prefix=/usr/local \
    -Dintrospection=disabled \
    -Dmodules=enabled \
    -Dcplusplus=true \
    -Djpeg=enabled \
    -Djpeg-xl=enabled && \
  ninja -C build && \
  ninja -C build install && \
  ldconfig /usr/local/lib && \
  rm -rf /tmp/vips-source

COPY go.mod go.sum ./
RUN go mod download

ENV PKG_CONFIG_PATH="/usr/lib/pkgconfig:/usr/local/lib/pkgconfig"
ENV LD_LIBRARY_PATH="/usr/lib:/usr/local/lib"

# Check JPEG library versions for debugging
RUN ldconfig && \
  pkg-config --modversion libjpeg || echo "No libjpeg pkg-config" && \
  pkg-config --modversion libturbojpeg || echo "No libturbojpeg pkg-config" && \
  ls -la /usr/lib/lib*jpeg* || echo "No jpeg libs in /usr/lib" && \
  ls -la /usr/local/lib/lib*jpeg* || echo "No jpeg libs in /usr/local/lib"

RUN go install github.com/cshum/vipsgen/cmd/vipsgen@latest
RUN vipsgen -out ./vips

COPY . .

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -x -v -a -tags vips \
  -ldflags="-s -w -linkmode external" \
  -o /bandwidth-hero-proxy ./main.go

# Runtime stage
FROM alpine:latest AS runtime

# Install only runtime dependencies
RUN apk add --no-cache \
  glib \
  expat \
  tiff \
  libpng \
  zlib \
  libwebp \
  libwebp-tools \
  giflib \
  librsvg \
  cairo \
  pango \
  fftw \
  orc \
  lcms2 \
  libexif \
  libgsf \
  openexr \
  libraw \
  libheif \
  libimagequant \
  openjpeg \
  libarchive \
  cfitsio \
  libspng \
  libjxl \
  cgif \
  highway \
  libjpeg-turbo \
  ca-certificates

# Create directories and copy the compiled binary and libraries from builder stage
RUN mkdir -p /usr/local/lib /usr/local/lib/pkgconfig /usr/lib
COPY --from=builder /bandwidth-hero-proxy /bandwidth-hero-proxy
COPY --from=builder /usr/local/lib/libvips* /usr/local/lib/
COPY --from=builder /usr/local/lib/pkgconfig/vips* /usr/local/lib/pkgconfig/

# Set up library path for runtime
ENV LD_LIBRARY_PATH="/usr/lib:/usr/local/lib"
RUN ldconfig /usr/local/lib

# Test that the binary can find its dependencies
RUN ldd /bandwidth-hero-proxy | grep vips || echo "Warning: vips library not found in dependencies"

ENTRYPOINT [ "/bandwidth-hero-proxy" ]
