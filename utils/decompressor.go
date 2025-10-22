package utils

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4/v4"
	"github.com/ulikunitz/xz"
)

// DecompressResponse automatically detects and decompresses HTTP response data
// based on Content-Encoding header and magic bytes
func DecompressResponse(data []byte, contentEncoding string) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	// Normalize content encoding
	encoding := strings.ToLower(strings.TrimSpace(contentEncoding))

	// If no encoding specified, try to detect by magic bytes
	if encoding == "" || encoding == "identity" {
		return DecompressByMagicBytes(data)
	}

	// Handle multiple encodings (e.g., "gzip, deflate")
	// Encodings are applied in reverse order (last applied first to decompress)
	encodings := strings.Split(encoding, ",")
	result := data

	for i := len(encodings) - 1; i >= 0; i-- {
		enc := strings.TrimSpace(encodings[i])
		if enc == "" || enc == "identity" {
			continue // Skip empty or identity encodings
		}

		var err error
		result, err = decompressSingle(result, enc)
		if err != nil {
			// If decompression fails, try magic bytes detection as fallback
			fallback, fallbackErr := DecompressByMagicBytes(result)
			if fallbackErr != nil {
				return nil, fmt.Errorf("failed to decompress with %s: %v (fallback also failed: %v)", enc, err, fallbackErr)
			}
			result = fallback
		}
	}

	return result, nil
}

// decompressSingle handles decompression for a single encoding type
func decompressSingle(data []byte, encoding string) ([]byte, error) {
	switch encoding {
	case "gzip":
		return DecompressGzip(data)
	case "deflate":
		return DecompressDeflate(data)
	case "br", "brotli":
		return DecompressBrotli(data)
	case "zstd":
		return DecompressZstd(data)
	case "lz4":
		return DecompressLZ4(data)
	case "xz":
		return DecompressXZ(data)
	case "identity", "":
		return data, nil // No compression
	default:
		// Try to detect by magic bytes if encoding is unknown
		return DecompressByMagicBytes(data)
	}
}

// DecompressGzip decompresses gzip-encoded data
func DecompressGzip(data []byte) ([]byte, error) {
	if len(data) < 2 {
		return data, nil
	}

	// Check gzip magic bytes
	if data[0] != 0x1f || data[1] != 0x8b {
		return data, nil // Not gzip data
	}

	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer reader.Close()

	result, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read gzip data: %v", err)
	}

	return result, nil
}

// DecompressDeflate decompresses deflate-encoded data
func DecompressDeflate(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	// Try raw deflate first
	reader := flate.NewReader(bytes.NewReader(data))
	defer reader.Close()

	result, err := io.ReadAll(reader)
	if err == nil {
		return result, nil
	}

	// Try zlib format (deflate with header)
	if len(data) >= 2 {
		zlibReader := flate.NewReader(bytes.NewReader(data))
		defer zlibReader.Close()

		zlibResult, zlibErr := io.ReadAll(zlibReader)
		if zlibErr == nil {
			return zlibResult, nil
		}
	}

	return nil, fmt.Errorf("failed to decompress deflate data: %v", err)
}

// DecompressBrotli decompresses brotli-encoded data
func DecompressBrotli(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	reader := brotli.NewReader(bytes.NewReader(data))
	result, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress brotli data: %v", err)
	}

	return result, nil
}

// DecompressZstd decompresses zstandard-encoded data
func DecompressZstd(data []byte) ([]byte, error) {
	decoder, err := zstd.NewReader(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create zstd decoder: %v", err)
	}
	defer decoder.Close()

	result, err := decoder.DecodeAll(data, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress zstd data: %v", err)
	}

	return result, nil
}

// DecompressLZ4 decompresses LZ4-encoded data
func DecompressLZ4(data []byte) ([]byte, error) {
	var out bytes.Buffer
	reader := lz4.NewReader(bytes.NewReader(data))

	_, err := io.Copy(&out, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress LZ4 data: %v", err)
	}

	return out.Bytes(), nil
}

// DecompressXZ decompresses XZ-encoded data
func DecompressXZ(data []byte) ([]byte, error) {
	reader, err := xz.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create xz reader: %v", err)
	}

	result, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress xz data: %v", err)
	}

	return result, nil
}

// DecompressByMagicBytes attempts to detect compression format by magic bytes
func DecompressByMagicBytes(data []byte) ([]byte, error) {
	if len(data) < 2 {
		return data, nil
	}

	// Try each format and return the first successful decompression
	// Gzip magic bytes: 0x1f 0x8b
	if data[0] == 0x1f && data[1] == 0x8b {
		if result, err := DecompressGzip(data); err == nil {
			return result, nil
		}
	}

	// Zlib magic bytes: 0x78 followed by various values
	if data[0] == 0x78 && (data[1] == 0x01 || data[1] == 0x5e || data[1] == 0x9c || data[1] == 0xda) {
		if result, err := DecompressDeflate(data); err == nil {
			return result, nil
		}
	}

	if len(data) >= 4 {
		// Zstandard magic bytes: 0x28 0xb5 0x2f 0xfd
		if data[0] == 0x28 && data[1] == 0xb5 && data[2] == 0x2f && data[3] == 0xfd {
			if result, err := DecompressZstd(data); err == nil {
				return result, nil
			}
		}

		// LZ4 magic bytes: 0x04 0x22 0x4d 0x18
		if data[0] == 0x04 && data[1] == 0x22 && data[2] == 0x4d && data[3] == 0x18 {
			if result, err := DecompressLZ4(data); err == nil {
				return result, nil
			}
		}
	}

	if len(data) >= 6 {
		// XZ magic bytes: 0xfd 0x37 0x7a 0x58 0x5a 0x00
		if data[0] == 0xfd && data[1] == 0x37 && data[2] == 0x7a &&
			data[3] == 0x58 && data[4] == 0x5a && data[5] == 0x00 {
			if result, err := DecompressXZ(data); err == nil {
				return result, nil
			}
		}
	}

	// If no magic bytes match or decompression fails, return original data
	return data, nil
}

// GetSupportedEncodings returns a list of all supported compression encodings
func GetSupportedEncodings() []string {
	return []string{
		"gzip",
		"deflate",
		"br", "brotli",
		"zstd",
		"lz4",
		"xz",
		"identity",
	}
}

// IsEncodingSupported checks if a given encoding is supported
func IsEncodingSupported(encoding string) bool {
	encoding = strings.ToLower(strings.TrimSpace(encoding))
	supported := GetSupportedEncodings()

	for _, enc := range supported {
		if encoding == enc {
			return true
		}
	}
	return false
}
