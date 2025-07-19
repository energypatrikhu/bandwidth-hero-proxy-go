package utils

import (
	"fmt"
	"net/http"
	"strings"
)

func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.WriteHeader(http.StatusOK)

	// Serve a blank favicon
	if _, err := w.Write([]byte{}); err != nil {
		fmt.Println("Error writing favicon response:", err)
	}
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	bhpParams, err := ParseParams(r)
	if err != nil {
		fmt.Fprint(w, "bandwidth-hero-proxy")
		fmt.Println(err)
		return
	}

	imageResponse, err := RequestImage(bhpParams.Url, r.Header)
	if err != nil {
		w.Header().Set("Location", bhpParams.Url)
		w.WriteHeader(http.StatusFound)

		fmt.Printf("\n> Params:\n > URL: %s\n > Format: %s\n > Quality: %d\n > Greyscale: %t\n> Info:\n > Error: %s\n > Action: Redirecting to original URL\n",
			bhpParams.Url, bhpParams.Format, bhpParams.Quality, bhpParams.Greyscale, err.Error())
		return
	}

	isAnimated := strings.Contains(imageResponse.ResponseHeaders.Get("Content-Type"), "image/gif")

	var compressedImg *CompressedImageResponse
	if BHP_USE_BEST_COMPRESSION_FORMAT && !isAnimated {
		compressedImg, err = CompressImageToBestFormat(imageResponse.Data, bhpParams.Greyscale, bhpParams.Quality)
	} else {
		compressedImg, err = CompressImage(imageResponse.Data, bhpParams.Format, bhpParams.Greyscale, bhpParams.Quality)
	}
	if err != nil {
		w.Header().Set("Location", bhpParams.Url)
		w.WriteHeader(http.StatusFound)

		fmt.Printf("\n> Params:\n > URL: %s\n > Format: %s\n > Quality: %d\n > Greyscale: %t\n> Info:\n > Error: %s\n > Action: Redirecting to original URL\n",
			bhpParams.Url, bhpParams.Format, bhpParams.Quality, bhpParams.Greyscale, err.Error())
		return
	}

	if !BHP_FORCE_FORMAT && compressedImg.Format == "" {
		w.Header().Set("Location", bhpParams.Url)
		w.WriteHeader(http.StatusFound)

		fmt.Printf("\n> Params:\n > URL: %s\n > Format: %s\n > Quality: %d\n > Greyscale: %t\n> Info:\n > Error: Could not compress image into smaller size than original\n > Action: Redirecting to original URL\n",
			bhpParams.Url, bhpParams.Format, bhpParams.Quality, bhpParams.Greyscale)
		return
	}

	compressedImageSize := len(compressedImg.Data)
	originalImageSize := len(imageResponse.Data)
	savedSize := originalImageSize - compressedImageSize

	if !BHP_FORCE_FORMAT && savedSize <= 0 {
		w.Header().Set("Location", bhpParams.Url)
		w.WriteHeader(http.StatusFound)

		fmt.Printf("\n> Params:\n > URL: %s\n > Format: %s\n > Quality: %d\n > Greyscale: %t\n> Info:\n > Error: Compressed image is not smaller than original\n > Action: Redirecting to original URL\n",
			bhpParams.Url, bhpParams.Format, bhpParams.Quality, bhpParams.Greyscale)
		return
	}

	for headerKey, headerValue := range imageResponse.ResponseHeaders {
		headerKey = strings.ToLower(headerKey)

		if headerKey == "Transfer-Encoding" {
			continue // Skip Content-Type and Content-Length headers
		}

		w.Header().Set(headerKey, headerValue[0]) // Set other headers from the original response
	}
	w.Header().Set("Content-Type", "image/"+compressedImg.Format)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", compressedImageSize))
	w.Header().Set("X-Original-Size", fmt.Sprintf("%d", originalImageSize))
	w.Header().Set("X-Compressed-Size", fmt.Sprintf("%d", compressedImageSize))
	w.Header().Set("X-Size-Saved", fmt.Sprintf("%d", savedSize))

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(compressedImg.Data); err != nil {
		http.Error(w, "Failed to write image response: "+err.Error(), http.StatusInternalServerError)
		fmt.Println("Error writing image response:", err)
		return
	}

	var reqHeaders string
	sortedRequestHeaders := GetSortedKeys(imageResponse.RequestHeaders)
	for _, k := range sortedRequestHeaders {
		v := imageResponse.RequestHeaders[k]
		k = strings.ToLower(k)
		reqHeaders += fmt.Sprintf(" > %s: %s\n", k, v)
	}

	var resHeaders string
	sortedResponseHeaders := GetSortedKeys(w.Header())
	for _, k := range sortedResponseHeaders {
		v := w.Header().Get(k)
		k = strings.ToLower(k)
		resHeaders += fmt.Sprintf(" > %s: %s\n", k, v)
	}

	compressedImageSizeStr := FormatSize(int64(compressedImageSize))
	originalImageSizeStr := FormatSize(int64(originalImageSize))
	savedSizeStr := FormatSize(int64(savedSize))

	compressedImageSizePerc := CalcPercentage(int64(compressedImageSize), int64(originalImageSize))
	savedSizePerc := CalcPercentage(int64(savedSize), int64(originalImageSize))

	formatModifiers := []string{}
	if BHP_FORCE_FORMAT {
		formatModifiers = append(formatModifiers, "forced")
	}
	if BHP_USE_BEST_COMPRESSION_FORMAT {
		formatModifiers = append(formatModifiers, "auto")
	}
	if isAnimated {
		formatModifiers = append(formatModifiers, "animated")
	}

	formatInfo := ""
	if len(formatModifiers) > 0 {
		formatInfo = fmt.Sprintf(" (%s)", strings.Join(formatModifiers, ", "))
	}

	fmt.Printf("\n> Params:\n > URL: %s\n > Format: %s\n > Quality: %d\n > Greyscale: %t\n> Request headers:\n%s> Response headers:\n%s> Info:\n > Using format: %s%s\n > Original size: %s\n > Compressed size: %s ( %s )\n > Saved size: %s ( %s )\n",
		bhpParams.Url, bhpParams.Format, bhpParams.Quality, bhpParams.Greyscale,
		reqHeaders, resHeaders, compressedImg.Format, formatInfo,
		originalImageSizeStr,
		compressedImageSizeStr, fmt.Sprintf("%.2f%%", compressedImageSizePerc),
		savedSizeStr, fmt.Sprintf("%.2f%%", savedSizePerc))
}
