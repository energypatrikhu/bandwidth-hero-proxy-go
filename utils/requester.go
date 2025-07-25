package utils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func RequestImage(url string, headers http.Header) (*ImageResponse, error) {
	requestHeaders := map[string]string{}

reqHeaderLoop:
	for headerKey, headerValue := range headers {
		headerKey = strings.ToLower(headerKey)

		if headerKey == "host" {
			continue // Skip Host header
		}

		for _, omittedHeader := range omittedHeadersRegexes {
			if omittedHeader.MatchString(headerKey) {
				continue reqHeaderLoop
			}
		}

		requestHeaders[headerKey] = headerValue[0]
	}

	duration, err := time.ParseDuration(BHP_EXTERNAL_REQUEST_TIMEOUT)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout duration: %v", err)
	}

	var resp *http.Response
	var imgData []byte
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

		imgData, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
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
	if resp == nil || imgData == nil {
		return nil, fmt.Errorf("no valid response received after %d attempts", BHP_EXTERNAL_REQUEST_RETRIES+1)
	}

	imageResponse := &ImageResponse{
		Data:            imgData,
		RequestHeaders:  requestHeaders,
		ResponseHeaders: resp.Header,
	}
	return imageResponse, nil
}
