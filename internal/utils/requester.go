package utils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	httpClient     *http.Client
	httpClientOnce sync.Once
	skipHeadersMap = map[string]bool{
		"host":            true,
		"accept-encoding": true,
	}
)

func RequestImage(url string, headers http.Header) (*ImageResponse, error) {
	requestHeaders := map[string]string{}

reqHeaderLoop:
	for headerKey, headerValue := range headers {
		headerKeyLower := strings.ToLower(headerKey)

		if skipHeadersMap[headerKeyLower] {
			continue
		}

		for _, omittedHeader := range omittedHeadersRegexes {
			if omittedHeader.MatchString(headerKeyLower) {
				continue reqHeaderLoop
			}
		}

		requestHeaders[headerKeyLower] = headerValue[0]
	}

	// Set Accept-Encoding header to handle all compression types we support
	requestHeaders["accept-encoding"] = "br, zstd, gzip, deflate, lz4, xz, identity"

	duration, err := time.ParseDuration(BHP_EXTERNAL_REQUEST_TIMEOUT)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout duration: %v", err)
	}

	// Initialize shared HTTP client once
	httpClientOnce.Do(func() {
		httpClient = &http.Client{
			Timeout: duration,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= BHP_EXTERNAL_REQUEST_REDIRECTS {
					return fmt.Errorf("stopped after %d redirects", BHP_EXTERNAL_REQUEST_REDIRECTS)
				}
				return nil
			},
		}
	})

	// If a FlareSolverr instance is configured, use it to solve any
	// anti-bot/Cloudflare challenge for this host and reuse the resulting
	// cookies + User-Agent for the actual fetch below.
	if strings.TrimSpace(BHP_FLARESOLVERR_URL) != "" {
		solution, err := SolveWithFlareSolverr(url, duration)
		if err != nil {
			return nil, fmt.Errorf("flaresolverr failed to solve challenge for %s: %v", url, err)
		}

		if solution.UserAgent != "" {
			requestHeaders["user-agent"] = solution.UserAgent
		}

		if len(solution.Cookies) > 0 {
			cookiePairs := make([]string, 0, len(solution.Cookies))
			for _, cookie := range solution.Cookies {
				cookiePairs = append(cookiePairs, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
			}
			requestHeaders["cookie"] = strings.Join(cookiePairs, "; ")
		}
	}

	var resp *http.Response
	var data []byte
	var lastErr error

	for attempt := 0; attempt < BHP_EXTERNAL_REQUEST_RETRIES+1; attempt++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			lastErr = err
			continue
		}

		for k, v := range requestHeaders {
			req.Header.Set(k, v)
		}

		resp, err = httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			statusCode := resp.StatusCode // Save status code before closing
			resp.Body.Close()
			resp = nil // Reset resp to nil for failed attempts
			lastErr = fmt.Errorf("failed to fetch image: status %d", statusCode)
			continue
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		// Decompress response data based on Content-Encoding header or magic bytes
		contentEncoding := resp.Header.Get("Content-Encoding")
		data, err = DecompressResponse(respBody, contentEncoding)
		if err != nil {
			lastErr = fmt.Errorf("failed to decompress response data (encoding: %s): %v", contentEncoding, err)
			continue
		}

		// Success
		lastErr = nil
		break
	}

	if lastErr != nil {
		return nil, lastErr
	}
	// Additional safety check - ensure we have valid response data
	if resp == nil || data == nil {
		return nil, fmt.Errorf("no valid response received after %d attempts", BHP_EXTERNAL_REQUEST_RETRIES+1)
	}

	imageResponse := &ImageResponse{
		Bytes:           data,
		RequestHeaders:  requestHeaders,
		ResponseHeaders: resp.Header,
	}
	return imageResponse, nil
}
