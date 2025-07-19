package utils

import "net/http"

type BhpParams struct {
	Url       string `json:"url"`
	Format    string `json:"format"`
	Greyscale bool   `json:"greyscale"`
	Quality   int    `json:"quality"`
}

type ImageResponse struct {
	Data            []byte            `json:"data"`
	RequestHeaders  map[string]string `json:"requestHeaders"`
	ResponseHeaders http.Header       `json:"responseHeaders"`
}

type CompressedImageResponse struct {
	Data   []byte `json:"data"`
	Format string `json:"format"`
}
