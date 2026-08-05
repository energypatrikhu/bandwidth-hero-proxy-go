package utils

import (
	"os"
	"runtime"
	"strconv"
	"strings"
)

func GetEnv[T any](key string, defaultValue T) T {
	switch any(defaultValue).(type) {
	case int:
		if value, exists := os.LookupEnv(key); exists {
			if intValue, err := strconv.Atoi(value); err == nil {
				return any(intValue).(T)
			}
		}
	case bool:
		if value, exists := os.LookupEnv(key); exists {
			if boolValue, err := strconv.ParseBool(value); err == nil {
				return any(boolValue).(T)
			}
		}
	case []string:
		if value, exists := os.LookupEnv(key); exists {
			splitFunc := func(r rune) bool {
				return r == '\n' || r == ';'
			}
			expandedValue := os.ExpandEnv(value)
			parts := strings.FieldsFunc(expandedValue, splitFunc)
			return any(parts).(T)
		}
	default:
		if value, exists := os.LookupEnv(key); exists {
			return any(value).(T) // Assuming the type matches
		}
	}

	return defaultValue
}

var (
	BHP_PORT                          = GetEnv("BHP_PORT", 80)
	BHP_MAX_CONCURRENCY               = GetEnv("BHP_MAX_CONCURRENCY", runtime.NumCPU())
	BHP_FORCE_FORMAT                  = GetEnv("BHP_FORCE_FORMAT", false)
	BHP_AUTO_DECREMENT_QUALITY        = GetEnv("BHP_AUTO_DECREMENT_QUALITY", false)
	BHP_USE_BEST_COMPRESSION_FORMAT   = GetEnv("BHP_USE_BEST_COMPRESSION_FORMAT", false)
	BHP_EXTERNAL_REQUEST_TIMEOUT      = GetEnv("BHP_EXTERNAL_REQUEST_TIMEOUT", "60s")
	BHP_EXTERNAL_REQUEST_RETRIES      = GetEnv("BHP_EXTERNAL_REQUEST_RETRIES", 5)
	BHP_EXTERNAL_REQUEST_REDIRECTS    = GetEnv("BHP_EXTERNAL_REQUEST_REDIRECTS", 10)
	BHP_EXTERNAL_REQUEST_OMIT_HEADERS = GetEnv("BHP_EXTERNAL_REQUEST_OMIT_HEADERS", []string{})
	BHP_FLARESOLVERR_URL              = GetEnv("BHP_FLARESOLVERR_URL", "")
	BHP_DISABLE_ANIMATED_IMAGES       = GetEnv("BHP_DISABLE_ANIMATED_IMAGES", false)

	BHP_WEBP_LOSSLESS        = GetEnv("BHP_WEBP_LOSSLESS", false)
	BHP_WEBP_EFFORT          = GetEnv("BHP_WEBP_EFFORT", 6)
	BHP_WEBP_SMART_SUBSAMPLE = GetEnv("BHP_WEBP_SMART_SUBSAMPLE", true)
	BHP_WEBP_SMART_DEBLOCK   = GetEnv("BHP_WEBP_SMART_DEBLOCK", true)
	BHP_WEBP_PASSES          = GetEnv("BHP_WEBP_PASSES", 10)

	BHP_JPEG_OPTIMIZE_CODING     = GetEnv("BHP_JPEG_OPTIMIZE_CODING", true)
	BHP_JPEG_OPTIMIZE_SCANS      = GetEnv("BHP_JPEG_OPTIMIZE_SCANS", true)
	BHP_JPEG_INTERLACE           = GetEnv("BHP_JPEG_INTERLACE", false)
	BHP_JPEG_TRELLIS_QUANT       = GetEnv("BHP_JPEG_TRELLIS_QUANT", true)
	BHP_JPEG_OVERSHOOT_DERINGING = GetEnv("BHP_JPEG_OVERSHOOT_DERINGING", true)
	BHP_JPEG_QUANT_TABLE         = GetEnv("BHP_JPEG_QUANT_TABLE", 3)
)
