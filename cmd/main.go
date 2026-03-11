package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/energypatrikhu/bandwidth-hero-proxy-go/internal/utils"
	"github.com/energypatrikhu/bandwidth-hero-proxy-go/third_party/vips"
)

func main() {
	log.Println("Starting Bandwidth Hero Proxy...")

	log.Println("> Config:")
	log.Println("  > BHP_PORT:", utils.BHP_PORT)
	log.Println("  > BHP_MAX_CONCURRENCY:", utils.BHP_MAX_CONCURRENCY)
	log.Println("  > BHP_FORCE_FORMAT:", utils.BHP_FORCE_FORMAT)
	log.Println("  > BHP_AUTO_DECREMENT_QUALITY:", utils.BHP_AUTO_DECREMENT_QUALITY)
	log.Println("  > BHP_USE_BEST_COMPRESSION_FORMAT:", utils.BHP_USE_BEST_COMPRESSION_FORMAT)
	log.Println("  > BHP_EXTERNAL_REQUEST_TIMEOUT:", utils.BHP_EXTERNAL_REQUEST_TIMEOUT)
	log.Println("  > BHP_EXTERNAL_REQUEST_RETRIES:", utils.BHP_EXTERNAL_REQUEST_RETRIES)
	log.Println("  > BHP_EXTERNAL_REQUEST_REDIRECTS:", utils.BHP_EXTERNAL_REQUEST_REDIRECTS)
	log.Println("  > BHP_EXTERNAL_REQUEST_OMIT_HEADERS:", utils.BHP_EXTERNAL_REQUEST_OMIT_HEADERS)

	if utils.BHP_FORCE_FORMAT && utils.BHP_USE_BEST_COMPRESSION_FORMAT {
		log.Println("Error: BHP_FORCE_FORMAT and BHP_USE_BEST_COMPRESSION_FORMAT cannot be both enabled at the same time.")
		return
	}

	if utils.BHP_USE_BEST_COMPRESSION_FORMAT && utils.BHP_AUTO_DECREMENT_QUALITY {
		log.Println("Error: BHP_USE_BEST_COMPRESSION_FORMAT and BHP_AUTO_DECREMENT_QUALITY cannot be both enabled at the same time.")
		return
	}

	vips.SetLogging(nil, 0) // Suppress vips logs
	vips.Startup(&vips.Config{
		ConcurrencyLevel: utils.BHP_MAX_CONCURRENCY, // Set concurrency level to BHP_MAX_CONCURRENCY
		MaxCacheFiles:    0,                         // Set max cache files to 0 (no limit)
		MaxCacheMem:      0,                         // Set max cache memory to 0 (no limit)
		MaxCacheSize:     0,                         // Set max cache size to 0 (no limit)
		ReportLeaks:      false,                     // Disable leak reporting
		CacheTrace:       false,                     // Disable cache tracing
		VectorEnabled:    true,                      // Enable vector support
	})
	defer vips.Shutdown()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /favicon.ico", utils.FaviconHandler)
	mux.HandleFunc("GET /", utils.ProxyHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", utils.BHP_PORT),
		Handler: mux,
	}

	log.Println("Server is running on port", utils.BHP_PORT)
	if err := server.ListenAndServe(); err != nil {
		log.Println("Error starting server:", err)
		return
	}
	log.Println("Server stopped")
}
