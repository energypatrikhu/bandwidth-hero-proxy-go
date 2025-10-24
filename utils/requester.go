package utils

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"
)

func RequestImage(url string, headers http.Header) (*ImageResponse, error) {
	requestHeaders := map[string]string{}

reqHeaderLoop:
	for headerKey, headerValue := range headers {
		headerKey = strings.ToLower(headerKey)

		if slices.Contains([]string{"host", "accept-encoding"}, headerKey) {
			continue
		}

		for _, omittedHeader := range omittedHeadersRegexes {
			if omittedHeader.MatchString(headerKey) {
				continue reqHeaderLoop
			}
		}

		requestHeaders[headerKey] = headerValue[0]
	}

	// Set Accept-Encoding header to handle all compression types we support
	requestHeaders["accept-encoding"] = "br, zstd, gzip, deflate, lz4, xz, identity"

	duration, err := time.ParseDuration(BHP_EXTERNAL_REQUEST_TIMEOUT)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout duration: %v", err)
	}

	var resp *http.Response
	var data []byte
	var lastErr error

	client := &http.Client{
		Timeout: duration,
		// Follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= BHP_EXTERNAL_REQUEST_REDIRECTS {
				return fmt.Errorf("stopped after %d redirects", BHP_EXTERNAL_REQUEST_REDIRECTS)
			}
			return nil
		},
	}

	for attempt := 0; attempt < BHP_EXTERNAL_REQUEST_RETRIES+1; attempt++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			lastErr = err
			continue
		}

		for k, v := range requestHeaders {
			req.Header.Set(k, v)
		}

		resp, err = client.Do(req)
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
		Bytes:           &data,
		RequestHeaders:  requestHeaders,
		ResponseHeaders: resp.Header,
	}
	return imageResponse, nil
}
