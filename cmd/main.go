package main

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/energypatrikhu/bandwidth-hero-proxy-go/internal/utils"
	"github.com/energypatrikhu/bandwidth-hero-proxy-go/third_party/vips"
)

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

	if utils.ConfigInstance.ForceFormat && utils.ConfigInstance.UseBestCompression {
		log.Panicln("Error: BHP_FORCE_FORMAT and BHP_USE_BEST_COMPRESSION_FORMAT cannot be both enabled at the same time.")
	}

	if utils.ConfigInstance.UseBestCompression && utils.ConfigInstance.AutoDecrementQuality {
		log.Panicln("Error: BHP_USE_BEST_COMPRESSION_FORMAT and BHP_AUTO_DECREMENT_QUALITY cannot be both enabled at the same time.")
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
