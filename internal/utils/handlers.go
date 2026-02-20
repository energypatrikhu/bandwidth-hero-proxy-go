package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method != http.MethodGet { // Only allow GET requests
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		fmt.Println("Error: Method Not Allowed")
		return
	}

	if r.URL.Path == "/favicon.ico" { // Handle favicon requests
		w.Header().Set("Content-Type", "image/x-icon")
		w.WriteHeader(http.StatusOK)

		// Serve a blank favicon
		if _, err := w.Write([]byte{}); err != nil {
			fmt.Println("Error writing favicon response:", err)
		}

		return
	}

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
	imageFormat := imageResponse.ResponseHeaders.Get("Content-Type")
	isAnimated := IsAnimatedFormat(imageFormat)
	originalImageSize := len(imageResponse.Bytes)

	currentQuality := bhpParams.Quality
	var compressedImage *CompressImageResult
	if BHP_USE_BEST_COMPRESSION_FORMAT && !isAnimated {
		compressedImage, err = CompressImageToBestFormat(imageResponse.Bytes, CompressImageToBestFormatOptions{
			InputFormat: imageFormat,
			Greyscale:   bhpParams.Greyscale,
			Quality:     bhpParams.Quality,
		})
	} else if BHP_AUTO_DECREMENT_QUALITY && !isAnimated {
		compressedImage, currentQuality, err = CompressImageWithAutoQualityDecrement(imageResponse.Bytes, CompressImageWithAutoQualityDecrementOptions{
			InputFormat:       imageFormat,
			Format:            bhpParams.Format,
			Greyscale:         bhpParams.Greyscale,
			InitialQuality:    bhpParams.Quality,
			OriginalImageSize: originalImageSize,
		})
	} else {
		compressedImage, err = CompressImage(imageResponse.Bytes, CompressImageOptions{
			InputFormat: imageFormat,
			IsAnimated:  isAnimated,
			Format:      bhpParams.Format,
			Greyscale:   bhpParams.Greyscale,
			Quality:     bhpParams.Quality,
		})
	}
	if err != nil {
		w.Header().Set("Location", bhpParams.Url)
		w.WriteHeader(http.StatusFound)

		fmt.Printf("\n> Params:\n > URL: %s\n > Format: %s\n > Quality: %d (%d)\n > Greyscale: %t\n> Info:\n > Error: %s\n > Action: Redirecting to original URL\n",
			bhpParams.Url, bhpParams.Format, bhpParams.Quality, currentQuality, bhpParams.Greyscale, err.Error())
		return
	}

	if !BHP_FORCE_FORMAT && compressedImage.Format == "" {
		w.Header().Set("Location", bhpParams.Url)
		w.WriteHeader(http.StatusFound)

		fmt.Printf("\n> Params:\n > URL: %s\n > Format: %s\n > Quality: %d (%d)\n > Greyscale: %t\n> Info:\n > Error: Could not compress image into smaller size than original\n > Action: Redirecting to original URL\n",
			bhpParams.Url, bhpParams.Format, bhpParams.Quality, currentQuality, bhpParams.Greyscale)
		return
	}

	compressedImageSize := len(compressedImage.Bytes)
	savedSize := originalImageSize - compressedImageSize

	if !BHP_FORCE_FORMAT && savedSize <= 0 {
		w.Header().Set("Location", bhpParams.Url)
		w.WriteHeader(http.StatusFound)

		fmt.Printf("\n> Params:\n > URL: %s\n > Format: %s\n > Quality: %d (%d)\n > Greyscale: %t\n> Info:\n > Error: Compressed image is not smaller than original\n > Action: Redirecting to original URL\n",
			bhpParams.Url, bhpParams.Format, bhpParams.Quality, currentQuality, bhpParams.Greyscale)
		return
	}

	skipHeaders := map[string]bool{"transfer-encoding": true, "content-encoding": true, "vary": true}
	for headerKey, headerValue := range imageResponse.ResponseHeaders {
		headerKeyLower := strings.ToLower(headerKey)

		if skipHeaders[headerKeyLower] {
			continue
		}

		w.Header().Set(headerKeyLower, headerValue[0]) // Set other headers from the original response
	}
	w.Header().Set("Content-Type", "image/"+compressedImage.Format)
	w.Header().Set("Content-Length", strconv.Itoa(compressedImageSize))
	w.Header().Set("X-Original-Size", strconv.Itoa(originalImageSize))
	w.Header().Set("X-Compressed-Size", strconv.Itoa(compressedImageSize))
	w.Header().Set("X-Size-Saved", strconv.Itoa(savedSize))

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(compressedImage.Bytes); err != nil {
		http.Error(w, "Failed to write image response: "+err.Error(), http.StatusInternalServerError)
		fmt.Println("Error writing image response:", err)
		return
	}

	var reqHeaders strings.Builder
	sortedRequestHeaders := GetSortedKeys(imageResponse.RequestHeaders)
	reqHeaders.Grow(len(sortedRequestHeaders) * 40) // Pre-allocate approximate size
	for _, k := range sortedRequestHeaders {
		v := imageResponse.RequestHeaders[k]
		kLower := strings.ToLower(k)
		reqHeaders.WriteString(" > ")
		reqHeaders.WriteString(kLower)
		reqHeaders.WriteString(": ")
		reqHeaders.WriteString(v)
		reqHeaders.WriteString("\n")
	}

	var resHeaders strings.Builder
	responseHeader := w.Header()
	sortedResponseHeaders := GetSortedKeys(responseHeader)
	resHeaders.Grow(len(sortedResponseHeaders) * 40) // Pre-allocate approximate size
	for _, k := range sortedResponseHeaders {
		v := responseHeader.Get(k)
		kLower := strings.ToLower(k)
		resHeaders.WriteString(" > ")
		resHeaders.WriteString(kLower)
		resHeaders.WriteString(": ")
		resHeaders.WriteString(v)
		resHeaders.WriteString("\n")
	}

	compressedImageSizeStr := FormatSize(int64(compressedImageSize))
	originalImageSizeStr := FormatSize(int64(originalImageSize))
	savedSizeStr := FormatSize(int64(savedSize))

	compressedImageSizePerc := CalcPercentage(int64(compressedImageSize), int64(originalImageSize))
	savedSizePerc := CalcPercentage(int64(savedSize), int64(originalImageSize))

	formatModifiers := make([]string, 0, 3)
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
		formatInfo = " (" + strings.Join(formatModifiers, ", ") + ")"
	}

	reqHeadersStr := reqHeaders.String()
	resHeadersStr := resHeaders.String()
	fmt.Printf("\n> Params:\n > URL: %s\n > Format: %s\n > Quality: %d (%d)\n > Greyscale: %t\n> Request headers:\n%s> Response headers:\n%s> Info:\n > Using format: %s%s\n > Original size: %s\n > Compressed size: %s ( %.2f%% )\n > Saved size: %s ( %.2f%% )\n",
		bhpParams.Url, bhpParams.Format, bhpParams.Quality, currentQuality, bhpParams.Greyscale,
		reqHeadersStr, resHeadersStr, compressedImage.Format, formatInfo,
		originalImageSizeStr,
		compressedImageSizeStr, compressedImageSizePerc,
		savedSizeStr, savedSizePerc)
}
