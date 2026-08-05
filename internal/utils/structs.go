package utils

import (
	"net/http"
)

type BhpParams struct {
	Url       string `json:"url"`
	Format    string `json:"format"`
	Grayscale bool   `json:"grayscale"`
	Quality   int    `json:"quality"`
}

type ImageResponse struct {
	Bytes           []byte
	RequestHeaders  map[string]string
	ResponseHeaders http.Header
}

type CompressImageResult struct {
	Bytes  []byte
	Format string
}

type CompressImageOptions struct {
	InputFormat string
	IsAnimated  bool
	Format      string
	Grayscale   bool
	Quality     int
}

type CompressImageWithAutoQualityDecrementOptions struct {
	InputFormat       string
	Format            string
	Grayscale         bool
	InitialQuality    int
	OriginalImageSize int
}

type CompressImageToBestFormatOptions struct {
	InputFormat string
	Grayscale   bool
	Quality     int
}
