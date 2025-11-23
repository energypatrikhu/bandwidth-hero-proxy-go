package utils

var AnimatedImageFormats = []string{
	"image/gif",
	"image/apng",
	"image/webp",
	"image/heif",
	"image/jxl",
	"image/tiff",
	"application/pdf",
}

var FormatsSupportingVipsUnlimited = []string{
	"image/heif",
	"image/jpeg",
	"image/png",
	"image/svg+xml",
	"image/tiff",
}

// Pre-computed maps for faster lookups
var (
	animatedFormatsMap = map[string]bool{
		"image/gif":        true,
		"image/apng":       true,
		"image/webp":       true,
		"image/heif":       true,
		"image/jxl":        true,
		"image/tiff":       true,
		"application/pdf": true,
	}
	unlimitedFormatsMap = map[string]bool{
		"image/heif":    true,
		"image/jpeg":    true,
		"image/png":     true,
		"image/svg+xml": true,
		"image/tiff":    true,
	}
)

func IsAnimatedFormat(format string) bool {
	return animatedFormatsMap[format]
}

func SupportsUnlimited(format string) bool {
	return unlimitedFormatsMap[format]
}
