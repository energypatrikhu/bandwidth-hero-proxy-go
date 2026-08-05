package utils

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func GetEnv[T any](key string, defaultValue T) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	var result T

	switch any(defaultValue).(type) {
	case string:
		return any(value).(T)

	case int:
		if v, err := strconv.Atoi(value); err == nil {
			return any(v).(T)
		}

	case int64:
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			return any(v).(T)
		}

	case bool:
		if v, err := strconv.ParseBool(value); err == nil {
			return any(v).(T)
		}

	case float64:
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			return any(v).(T)
		}

	case []string:
		value = os.ExpandEnv(value)

		parts := strings.FieldsFunc(value, func(r rune) bool {
			return r == '\n' || r == ';' || r == ','
		})

		return any(parts).(T)

	case time.Duration:
		if v, err := time.ParseDuration(value); err == nil {
			return any(v).(T)
		}
	}

	// fallback
	result = defaultValue
	return result
}
