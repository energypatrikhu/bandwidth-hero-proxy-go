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
)
