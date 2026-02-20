package utils

import "fmt"

func FormatSize(bytes int64) string {
	const (
		_          = iota
		KB float64 = 1 << (10 * iota)
		MB
		GB
		TB
		PB
	)

	b := float64(bytes)

	switch {
	case b >= PB:
		return fmt.Sprintf("%.2f PB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2f TB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2f GB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2f MB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2f KB", b/KB)
	case b >= 1:
		return fmt.Sprintf("%.0f B", b)
	default:
		return fmt.Sprintf("%.0f b", b*8)
	}
}

func CalcPercentage(part, total int64) float64 {
	if total == 0 {
		return 0
	}
	return (float64(part) / float64(total)) * 100
}
