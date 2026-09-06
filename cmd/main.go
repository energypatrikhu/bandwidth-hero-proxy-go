package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/energypatrikhu/bandwidth-hero-proxy-go/internal/utils"
	"github.com/energypatrikhu/bandwidth-hero-proxy-go/third_party/vips"
)

var validFormats = map[string]bool{
	"webp": true,
	"jpeg": true,
	"jxl":  true,
}

func main() {
	log.Println("Starting Bandwidth Hero Proxy...")

	log.Println("> Config:")

	t := reflect.TypeFor[utils.Config]()
	v := reflect.ValueOf(utils.ConfigInstance)

	for _, field := range reflect.VisibleFields(t) {
		name := field.Tag.Get("env")
		if name == "" {
			name = field.Name
		}

		value := v.FieldByName(field.Name).Interface()
		if field.Type.Kind() == reflect.String && value == "" {
			value = "not set"
		}

		log.Printf(" > %s: %v", name, value)
	}

	if utils.ConfigInstance.TargetFormat != "" && utils.ConfigInstance.AutoSelectFormat {
		log.Panicln("Error: BHP_TARGET_FORMAT and BHP_AUTO_SELECT_FORMAT cannot be both enabled at the same time.")
	}

	if utils.ConfigInstance.AutoSelectFormat && utils.ConfigInstance.AutoReduceQuality {
		log.Panicln("Error: BHP_AUTO_SELECT_FORMAT and BHP_AUTO_REDUCE_QUALITY cannot be both enabled at the same time.")
	}

	if utils.ConfigInstance.IgnoreSizeCheck && utils.ConfigInstance.AutoSelectFormat {
		log.Panicln("Error: BHP_IGNORE_SIZE_CHECK and BHP_AUTO_SELECT_FORMAT cannot be both enabled at the same time.")
	}

	if utils.ConfigInstance.TargetFormat != "" {
		if !validFormats[utils.ConfigInstance.TargetFormat] {
			log.Panicf("Error: BHP_TARGET_FORMAT must be one of the following: webp, jpeg, jxl. Got: %s", utils.ConfigInstance.TargetFormat)
		}
	}

	for _, format := range utils.ConfigInstance.AutoSelectFormatList {
		if !validFormats[format] {
			log.Panicf("Error: BHP_AUTO_SELECT_FORMAT_LIST must contain only the following formats: webp, jpeg, jxl. Got: %s", format)
		}
	}

	// WEBP options validation
	if utils.ConfigInstance.WebpEffort < 1 || utils.ConfigInstance.WebpEffort > 6 {
		log.Panicln("Error: BHP_WEBP_EFFORT must be between 1 and 6.")
	}
	if utils.ConfigInstance.WebpPasses < 1 || utils.ConfigInstance.WebpPasses > 10 {
		log.Panicln("Error: BHP_WEBP_PASSES must be between 1 and 10.")
	}

	// JPEG options validation
	if utils.ConfigInstance.JpegQuantTable < 0 || utils.ConfigInstance.JpegQuantTable > 8 {
		log.Panicln("Error: BHP_JPEG_QUANT_TABLE must be between 0 and 8.")
	}

	// JXL options validation
	if utils.ConfigInstance.JxlTier < 0 || utils.ConfigInstance.JxlTier > 4 {
		log.Panicln("Error: BHP_JXL_TIER must be between 0 and 4.")
	}
	if utils.ConfigInstance.JxlEffort < 1 || utils.ConfigInstance.JxlEffort > 9 {
		log.Panicln("Error: BHP_JXL_EFFORT must be between 1 and 9.")
	}

	vips.SetLogging(nil, 0) // Suppress vips logs
	vips.Startup(&vips.Config{
		ConcurrencyLevel: utils.ConfigInstance.VipsMaxConcurrency, // Set concurrency level to BHP_MAX_CONCURRENCY
		MaxCacheFiles:    0,                                       // Set max cache files to 0 (disable)
		MaxCacheMem:      0,                                       // Set max cache memory to 0 (disable)
		MaxCacheSize:     0,                                       // Set max cache size to 0 (disable)
		ReportLeaks:      false,                                   // Disable leak reporting
		CacheTrace:       false,                                   // Disable cache tracing
		VectorEnabled:    true,                                    // Enable vector support
	})
	defer vips.Shutdown()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /favicon.ico", utils.FaviconHandler)
	mux.HandleFunc("GET /", utils.ProxyHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", utils.ConfigInstance.Port),
		Handler: mux,
	}

	log.Println("Server is running on port", utils.ConfigInstance.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Panicln("Error starting server:", err)
		return
	}
	log.Println("Server stopped")
}
