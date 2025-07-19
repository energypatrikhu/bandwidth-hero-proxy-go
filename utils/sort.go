package utils

import "sort"

func GetSortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func SortMapByKey[T any](m map[string]T) map[string]T {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortedMap := make(map[string]T, len(m))
	for _, k := range keys {
		sortedMap[k] = m[k]
	}
	return sortedMap
}
