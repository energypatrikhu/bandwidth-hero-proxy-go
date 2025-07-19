package utils

func ToMap[T any](m map[string]T) map[string]T {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sortedMap := make(map[string]T, len(m))
	for _, k := range keys {
		sortedMap[k] = m[k]
	}
	return sortedMap
}
