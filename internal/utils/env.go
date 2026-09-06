package utils

import (
	"runtime"
)

type Config struct {
	Port            int    `env:"BHP_PORT"`
	FlareSolverrURL string `env:"BHP_FLARESOLVERR_URL"`
	VerboseLogging  bool   `env:"BHP_VERBOSE_LOGGING"`

	VipsMaxConcurrency int `env:"BHP_VIPS_MAX_CONCURRENCY"`

	ForceFormat           bool     `env:"BHP_FORCE_FORMAT"`
	CustomFormat          string   `env:"BHP_CUSTOM_FORMAT"`
	TestFormats           []string `env:"BHP_TEST_FORMATS"`
	AutoDecrementQuality  bool     `env:"BHP_AUTO_DECREMENT_QUALITY"`
	UseBestCompression    bool     `env:"BHP_USE_BEST_COMPRESSION_FORMAT"`
	DisableAnimatedImages bool     `env:"BHP_DISABLE_ANIMATED_IMAGES"`

	ExternalRequestTimeout     string   `env:"BHP_EXTERNAL_REQUEST_TIMEOUT"`
	ExternalRequestRetries     int      `env:"BHP_EXTERNAL_REQUEST_RETRIES"`
	ExternalRequestRedirects   int      `env:"BHP_EXTERNAL_REQUEST_REDIRECTS"`
	ExternalRequestOmitHeaders []string `env:"BHP_EXTERNAL_REQUEST_OMIT_HEADERS"`

	WebpLossless       bool `env:"BHP_WEBP_LOSSLESS"`
	WebpEffort         int  `env:"BHP_WEBP_EFFORT"`
	WebpSmartSubsample bool `env:"BHP_WEBP_SMART_SUBSAMPLE"`
	WebpSmartDeblock   bool `env:"BHP_WEBP_SMART_DEBLOCK"`
	WebpPasses         int  `env:"BHP_WEBP_PASSES"`

	JpegOptimizeCoding     bool `env:"BHP_JPEG_OPTIMIZE_CODING"`
	JpegOptimizeScans      bool `env:"BHP_JPEG_OPTIMIZE_SCANS"`
	JpegInterlace          bool `env:"BHP_JPEG_INTERLACE"`
	JpegTrellisQuant       bool `env:"BHP_JPEG_TRELLIS_QUANT"`
	JpegOvershootDeringing bool `env:"BHP_JPEG_OVERSHOOT_DERINGING"`
	JpegQuantTable         int  `env:"BHP_JPEG_QUANT_TABLE"`

	JxlTier     int  `env:"BHP_JXL_TIER"`
	JxlEffort   int  `env:"BHP_JXL_EFFORT"`
	JxlLossless bool `env:"BHP_JXL_LOSSLESS"`
}

var ConfigInstance = Config{
	Port:            GetEnv("BHP_PORT", 80),
	FlareSolverrURL: GetEnv("BHP_FLARESOLVERR_URL", ""),

	VipsMaxConcurrency: GetEnv("BHP_VIPS_MAX_CONCURRENCY", runtime.NumCPU()),

	ForceFormat:           GetEnv("BHP_FORCE_FORMAT", false),
	CustomFormat:          GetEnv("BHP_CUSTOM_FORMAT", ""),
	TestFormats:           GetEnv("BHP_TEST_FORMATS", []string{"jxl", "webp", "jpeg"}),
	AutoDecrementQuality:  GetEnv("BHP_AUTO_DECREMENT_QUALITY", false),
	UseBestCompression:    GetEnv("BHP_USE_BEST_COMPRESSION_FORMAT", false),
	DisableAnimatedImages: GetEnv("BHP_DISABLE_ANIMATED_IMAGES", false),
	VerboseLogging:        GetEnv("BHP_VERBOSE_LOGGING", false),

	ExternalRequestTimeout:     GetEnv("BHP_EXTERNAL_REQUEST_TIMEOUT", "60s"),
	ExternalRequestRetries:     GetEnv("BHP_EXTERNAL_REQUEST_RETRIES", 5),
	ExternalRequestRedirects:   GetEnv("BHP_EXTERNAL_REQUEST_REDIRECTS", 10),
	ExternalRequestOmitHeaders: GetEnv("BHP_EXTERNAL_REQUEST_OMIT_HEADERS", []string{}),

	WebpLossless:       GetEnv("BHP_WEBP_LOSSLESS", false),
	WebpEffort:         GetEnv("BHP_WEBP_EFFORT", 4),
	WebpSmartSubsample: GetEnv("BHP_WEBP_SMART_SUBSAMPLE", true),
	WebpSmartDeblock:   GetEnv("BHP_WEBP_SMART_DEBLOCK", true),
	WebpPasses:         GetEnv("BHP_WEBP_PASSES", 1),

	JpegOptimizeCoding:     GetEnv("BHP_JPEG_OPTIMIZE_CODING", true),
	JpegOptimizeScans:      GetEnv("BHP_JPEG_OPTIMIZE_SCANS", true),
	JpegInterlace:          GetEnv("BHP_JPEG_INTERLACE", true),
	JpegTrellisQuant:       GetEnv("BHP_JPEG_TRELLIS_QUANT", true),
	JpegOvershootDeringing: GetEnv("BHP_JPEG_OVERSHOOT_DERINGING", true),
	JpegQuantTable:         GetEnv("BHP_JPEG_QUANT_TABLE", 3),

	JxlTier:     GetEnv("BHP_JXL_TIER", 0),
	JxlEffort:   GetEnv("BHP_JXL_EFFORT", 7),
	JxlLossless: GetEnv("BHP_JXL_LOSSLESS", false),
}
