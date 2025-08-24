package utils

import (
	"fmt"
	"net/http"
)

func ParseParams(r *http.Request) (*BhpParams, error) {
	query := r.URL.Query()

	url := query.Get("url")
	if url == "" {
		return nil, fmt.Errorf("Missing required parameter: url")
	}
	url = inputUrlRegex.ReplaceAllString(query.Get("url"), "http://")

	format := "webp" // Set webp as default format
	if query.Get("jpg") == "1" {
		format = "jpeg"
	}

	greyscale := false // Disable greyscale by default
	if query.Get("bw") == "1" {
		greyscale = true
	}

	quality := 80 // Set default quality to 80
	if query.Get("l") != "" {
		fmt.Sscanf(query.Get("l"), "%d", &quality)
	}

	params := &BhpParams{
		Url:       url,
		Format:    format,
		Greyscale: greyscale,
		Quality:   quality,
	}

	return params, nil
}
