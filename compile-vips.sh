#!/bin/bash

set -e

# Install build dependencies and all vips optional dependencies
apt-get update
apt-get install -y \
  build-essential \
  pkg-config \
  git \
  cmake \
  meson \
  ninja-build \
  nasm \
  autoconf \
  automake \
  libtool \
  libgirepository1.0-dev \
  gi-docgen \
  doxygen \
  libglib2.0-dev \
  libexpat1-dev \
  libtiff5-dev \
  libpng-dev \
  zlib1g-dev \
  libwebp-dev \
  libgif-dev \
  librsvg2-dev \
  libpoppler-glib-dev \
  libcairo2-dev \
  libpango1.0-dev \
  libfftw3-dev \
  liborc-0.4-dev \
  liblcms2-dev \
  libexif-dev \
  libgsf-1-dev \
  libopenexr-dev \
  libraw-dev \
  libheif-dev \
  libimagequant-dev \
  libopenjp2-7-dev \
  libarchive-dev \
  libcfitsio-dev \
  libspng-dev \
  libjxl-dev \
  libcgif-dev \
  libhighwayhash-dev \
  libjpeg62-turbo-dev \
  libmagick++-dev \
  libopenslide-dev \
  libmatio-dev

# Build and install the latest version of libvips with JPEG and JPEG XL support
TMP_DIR=$(mktemp -d)
VIPS_VERSION=$(wget -qO- "https://api.github.com/repos/libvips/libvips/releases/latest" | grep -o '"tag_name": "v[^"]*"' | cut -d'"' -f4 | sed 's/^v//')

wget -O $TMP_DIR/vips.tar.xz "https://github.com/libvips/libvips/releases/latest/download/vips-${VIPS_VERSION}.tar.xz"
cd $TMP_DIR

tar -xf vips.tar.xz
cd vips-${VIPS_VERSION}

sed -i "s/value: 'auto',/value: 'enabled',/g" meson_options.txt

PKG_CONFIG_PATH="/usr/lib/pkgconfig:/usr/local/lib/pkgconfig" \
meson setup build \
  --buildtype=release \
  --prefix=/usr/local \
  -Ddocs=true \
  -Dcpp-docs=true \
  -Dintrospection=enabled \
  -Dmodules=enabled \
  -Dcplusplus=true \
  -Dpdfium=disabled \
  -Dnifti=disabled

ninja -C build
ninja -C build install

ldconfig /usr/local/lib

rm -rf $TMP_DIR
