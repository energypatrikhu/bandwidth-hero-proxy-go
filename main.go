package main

import (
	"fmt"
	"net/http"

	"github.com/davidbyttow/govips/v2/vips"

	"github.com/energypatrikhu/bandwidth-hero-proxy-go/utils"
)

func main() {
	fmt.Println("Starting Bandwidth Hero Proxy...")

	fmt.Println("> Config:")
	fmt.Println(" > BHP_PORT:", utils.BHP_PORT)
	fmt.Println(" > BHP_MAX_CONCURRENCY:", utils.BHP_MAX_CONCURRENCY)
	fmt.Println(" > BHP_FORCE_FORMAT:", utils.BHP_FORCE_FORMAT)
	fmt.Println(" > BHP_AUTO_DECREMENT_QUALITY:", utils.BHP_AUTO_DECREMENT_QUALITY)
	fmt.Println(" > BHP_USE_BEST_COMPRESSION_FORMAT:", utils.BHP_USE_BEST_COMPRESSION_FORMAT)
	fmt.Println(" > BHP_EXTERNAL_REQUEST_TIMEOUT:", utils.BHP_EXTERNAL_REQUEST_TIMEOUT)
	fmt.Println(" > BHP_EXTERNAL_REQUEST_RETRIES:", utils.BHP_EXTERNAL_REQUEST_RETRIES)
	fmt.Println(" > BHP_EXTERNAL_REQUEST_REDIRECTS:", utils.BHP_EXTERNAL_REQUEST_REDIRECTS)
	fmt.Println(" > BHP_EXTERNAL_REQUEST_OMIT_HEADERS:", utils.BHP_EXTERNAL_REQUEST_OMIT_HEADERS)

	if utils.BHP_FORCE_FORMAT && utils.BHP_USE_BEST_COMPRESSION_FORMAT {
		fmt.Println("Error: BHP_FORCE_FORMAT and BHP_USE_BEST_COMPRESSION_FORMAT cannot be both enabled at the same time.")
		return
	}

	if utils.BHP_USE_BEST_COMPRESSION_FORMAT && utils.BHP_AUTO_DECREMENT_QUALITY {
		fmt.Println("Error: BHP_USE_BEST_COMPRESSION_FORMAT and BHP_AUTO_DECREMENT_QUALITY cannot be both enabled at the same time.")
		return
	}

	vips.LoggingSettings(nil, 0) // Suppress vips logs
	vips.Startup(&vips.Config{
		ConcurrencyLevel: utils.BHP_MAX_CONCURRENCY, // Set concurrency level to BHP_MAX_CONCURRENCY
		MaxCacheFiles:    0,                         // Set max cache files to 0 (no limit)
		MaxCacheMem:      0,                         // Set max cache memory to 0 (no limit)
		MaxCacheSize:     0,                         // Set max cache size to 0 (no limit)
		CacheTrace:       false,                     // Disable cache tracing
		ReportLeaks:      false,                     // Disable leak reporting
		CollectStats:     false,                     // Disable stats collection
	})
	defer vips.Shutdown()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", utils.BHP_PORT),
		Handler: http.HandlerFunc(utils.ProxyHandler),
	}

	fmt.Println("Server is running on port", utils.BHP_PORT)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	fmt.Println("Server stopped")
}
