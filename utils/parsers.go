package utils

import (
	"fmt"
	"net/http"
	"strconv"
)

func ParseParams(r *http.Request) (*BhpParams, error) {
	query := r.URL.Query()

	urlParam := query.Get("url")
	if urlParam == "" {
		return nil, fmt.Errorf("missing required parameter: url")
	}
	url := inputUrlRegex.ReplaceAllString(urlParam, "http://")

	format := "webp" // Set webp as default format
	if query.Get("jpg") == "1" {
		format = "jpeg"
	}

	greyscale := query.Get("bw") == "1"

	quality := 80 // Set default quality to 80
	if qualityStr := query.Get("l"); qualityStr != "" {
		if parsedQuality, err := strconv.Atoi(qualityStr); err == nil && parsedQuality > 0 && parsedQuality <= 100 {
			quality = parsedQuality
		}
	}

	return &BhpParams{
		Url:       url,
		Format:    format,
		Greyscale: greyscale,
		Quality:   quality,
	}, nil
}
