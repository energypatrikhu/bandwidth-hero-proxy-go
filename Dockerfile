FROM golang:alpine AS builder

# Install build dependencies and all vips optional dependencies
RUN apk del --no-cache \
  *libjpeg*

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
  gobject-introspection-dev \
  gi-docgen \
  doxygen \
  vala \
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
  libjpeg-turbo-dev \
  imagemagick-dev

# Build vips from source
RUN VIPS_VERSION=$(wget -qO- "https://api.github.com/repos/libvips/libvips/releases/latest" | grep -o '"tag_name": "v[^"]*"' | cut -d'"' -f4 | sed 's/^v//') && \
  mkdir -p /tmp/vips-source && \
  wget -O /tmp/vips-source/vips.tar.xz "https://github.com/libvips/libvips/releases/latest/download/vips-${VIPS_VERSION}.tar.xz" && \
  cd /tmp/vips-source && \
  tar -xf vips.tar.xz && \
  cd vips-${VIPS_VERSION} && \
  sed -i "s/value: 'auto',/value: 'enabled',/g" meson_options.txt && \
  PKG_CONFIG_PATH="/usr/lib/pkgconfig:/usr/local/lib/pkgconfig" \
  meson setup build \
    --buildtype=release \
    --prefix=/usr/local \
    -Ddocs=true \
    -Dcpp-docs=true \
    -Dintrospection=enabled \
    -Dmodules=enabled \
    -Dcplusplus=true \
    -Dvapi=true \
    -Dmatio=disabled \
    -Dpdfium=disabled \
    -Dnifti=disabled \
    -Dopenslide=disabled \
    -Duhdr=disabled && \
  ninja -C build && \
  ninja -C build install && \
  ldconfig /usr/local/lib && \
  rm -rf /tmp/vips-source

ENV PKG_CONFIG_PATH="/usr/lib/pkgconfig:/usr/local/lib/pkgconfig"
ENV LD_LIBRARY_PATH="/usr/lib:/usr/local/lib"

WORKDIR /tmp/app-source

RUN mkdir -p ./third_party && \
  go install github.com/cshum/vipsgen/cmd/vipsgen@latest && \
  vipsgen -out ./third_party/vips

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal

ENV CGO_ENABLED=1 \
  GOOS=linux \
  GOARCH=amd64

RUN go build -x -v -a -tags vips \
  -ldflags="-s -w -linkmode external" \
  -o /bandwidth-hero-proxy ./cmd/main.go

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
  poppler \
  poppler-glib \
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
  imagemagick \
  ca-certificates

# Create directories and copy the compiled binary and all necessary VIPS files
RUN mkdir -p /usr/local/lib /usr/local/lib/pkgconfig /usr/local/include /usr/local/bin

# Copy the main binary
COPY --from=builder /bandwidth-hero-proxy /bandwidth-hero-proxy

# Copy all VIPS-related files from builder
COPY --from=builder /usr/local/lib/ /usr/local/lib/
COPY --from=builder /usr/local/include/vips* /usr/local/include/
COPY --from=builder /usr/local/bin/vips* /usr/local/bin/

# Clean up unnecessary files to reduce image size
RUN find /usr/local/lib -name "*.a" -delete && \
    find /usr/local/lib -name "*.la" -delete

# Set up library path for runtime
ENV PKG_CONFIG_PATH="/usr/lib/pkgconfig:/usr/local/lib/pkgconfig"
ENV LD_LIBRARY_PATH="/usr/lib:/usr/local/lib"

# Update library cache
RUN ldconfig /usr/local/lib

WORKDIR /

ENTRYPOINT [ "/bandwidth-hero-proxy" ]
