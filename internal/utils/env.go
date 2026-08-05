package utils

import (
	"runtime"
)

type Config struct {
	Port            int    `env:"BHP_PORT"`
	FlareSolverrURL string `env:"BHP_FLARESOLVERR_URL"`
	VerboseLogging  bool   `env:"BHP_VERBOSE_LOGGING"`

	VipsMaxConcurrency int `env:"BHP_VIPS_MAX_CONCURRENCY"`

	ForceFormat           bool `env:"BHP_FORCE_FORMAT"`
	AutoDecrementQuality  bool `env:"BHP_AUTO_DECREMENT_QUALITY"`
	UseBestCompression    bool `env:"BHP_USE_BEST_COMPRESSION_FORMAT"`
	DisableAnimatedImages bool `env:"BHP_DISABLE_ANIMATED_IMAGES"`

	ExternalRequestTimeout     string   `env:"BHP_EXTERNAL_REQUEST_TIMEOUT"`
	ExternalRequestRetries     int      `env:"BHP_EXTERNAL_REQUEST_RETRIES"`
	ExternalRequestRedirects   int      `env:"BHP_EXTERNAL_REQUEST_REDIRECTS"`
	ExternalRequestOmitHeaders []string `env:"BHP_EXTERNAL_REQUEST_OMIT_HEADERS"`

	WebPLossless       bool `env:"BHP_WEBP_LOSSLESS"`
	WebPEffort         int  `env:"BHP_WEBP_EFFORT"`
	WebPSmartSubsample bool `env:"BHP_WEBP_SMART_SUBSAMPLE"`
	WebPSmartDeblock   bool `env:"BHP_WEBP_SMART_DEBLOCK"`
	WebPPasses         int  `env:"BHP_WEBP_PASSES"`

	JPEGOptimizeCoding     bool `env:"BHP_JPEG_OPTIMIZE_CODING"`
	JPEGOptimizeScans      bool `env:"BHP_JPEG_OPTIMIZE_SCANS"`
	JPEGInterlace          bool `env:"BHP_JPEG_INTERLACE"`
	JPEGTrellisQuant       bool `env:"BHP_JPEG_TRELLIS_QUANT"`
	JPEGOvershootDeringing bool `env:"BHP_JPEG_OVERSHOOT_DERINGING"`
	JPEGQuantTable         int  `env:"BHP_JPEG_QUANT_TABLE"`
}

var ConfigInstance = Config{
	Port:            GetEnv("BHP_PORT", 80),
	FlareSolverrURL: GetEnv("BHP_FLARESOLVERR_URL", ""),

	VipsMaxConcurrency: GetEnv("BHP_VIPS_MAX_CONCURRENCY", runtime.NumCPU()),

	ForceFormat:           GetEnv("BHP_FORCE_FORMAT", false),
	AutoDecrementQuality:  GetEnv("BHP_AUTO_DECREMENT_QUALITY", false),
	UseBestCompression:    GetEnv("BHP_USE_BEST_COMPRESSION_FORMAT", false),
	DisableAnimatedImages: GetEnv("BHP_DISABLE_ANIMATED_IMAGES", false),
	VerboseLogging:        GetEnv("BHP_VERBOSE_LOGGING", false),

	ExternalRequestTimeout:     GetEnv("BHP_EXTERNAL_REQUEST_TIMEOUT", "60s"),
	ExternalRequestRetries:     GetEnv("BHP_EXTERNAL_REQUEST_RETRIES", 5),
	ExternalRequestRedirects:   GetEnv("BHP_EXTERNAL_REQUEST_REDIRECTS", 10),
	ExternalRequestOmitHeaders: GetEnv("BHP_EXTERNAL_REQUEST_OMIT_HEADERS", []string{}),

	WebPLossless:       GetEnv("BHP_WEBP_LOSSLESS", false),
	WebPEffort:         GetEnv("BHP_WEBP_EFFORT", 4),
	WebPSmartSubsample: GetEnv("BHP_WEBP_SMART_SUBSAMPLE", true),
	WebPSmartDeblock:   GetEnv("BHP_WEBP_SMART_DEBLOCK", true),
	WebPPasses:         GetEnv("BHP_WEBP_PASSES", 2),

	JPEGOptimizeCoding:     GetEnv("BHP_JPEG_OPTIMIZE_CODING", true),
	JPEGOptimizeScans:      GetEnv("BHP_JPEG_OPTIMIZE_SCANS", true),
	JPEGInterlace:          GetEnv("BHP_JPEG_INTERLACE", false),
	JPEGTrellisQuant:       GetEnv("BHP_JPEG_TRELLIS_QUANT", true),
	JPEGOvershootDeringing: GetEnv("BHP_JPEG_OVERSHOOT_DERINGING", true),
	JPEGQuantTable:         GetEnv("BHP_JPEG_QUANT_TABLE", 3),
}
