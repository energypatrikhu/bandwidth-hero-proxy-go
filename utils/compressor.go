package utils

import (
	"fmt"

	"github.com/energypatrikhu/bandwidth-hero-proxy-go/vips"
)

func CompressImage(imageBytes []byte, options CompressImageOptions) (*CompressImageResult, error) {
	loadOptions := &vips.LoadOptions{
		FailOnError: false,
	}
	if options.IsAnimated {
		loadOptions.N = -1 // Load all frames for animated images
	}
	if SupportsUnlimited(options.InputFormat) {
		loadOptions.Unlimited = true // Allow unlimited image size for supported formats
	}

	vipsImage, vipsError := vips.NewImageFromBuffer(imageBytes, loadOptions)
	if vipsError != nil {
		return nil, fmt.Errorf("failed to create image from buffer: %w", vipsError)
	}
	defer vipsImage.Close()

	vipsImage.RemoveICCProfile()

	if options.Greyscale {
		vipsImage.Colourspace(vips.InterpretationBW, nil)
	}

	var compressedImageBytes []byte

	switch options.Format {
	case "webp":
		compressedImageBytes, vipsError = vipsImage.WebpsaveBuffer(&vips.WebpsaveBufferOptions{
			Q:        options.Quality,
			Lossless: false,
			Keep:     vips.KeepNone,
			Effort:   6,
		})
	case "jpeg":
		compressedImageBytes, vipsError = vipsImage.JpegsaveBuffer(&vips.JpegsaveBufferOptions{
			Q:                  options.Quality,
			OptimizeCoding:     true,
			OptimizeScans:      true,
			Keep:               vips.KeepNone,
			Interlace:          false,
			SubsampleMode:      vips.SubsampleAuto,
			TrellisQuant:       true,
			OvershootDeringing: true,
			QuantTable:         3,
		})
	}

	if vipsError != nil {
		return nil, fmt.Errorf("failed to export image buffer: %w", vipsError)
	}

	return &CompressImageResult{Bytes: compressedImageBytes, Format: options.Format}, nil
}

func CompressImageWithAutoQualityDecrement(imageBytes []byte, options CompressImageWithAutoQualityDecrementOptions) (*CompressImageResult, int, error) {
	currentQuality := options.InitialQuality
	var compressedImage *CompressImageResult
	var err error

	// Reuse options struct to reduce allocations
	compressOpts := CompressImageOptions{
		InputFormat: options.InputFormat,
		Format:      options.Format,
		Greyscale:   options.Greyscale,
		IsAnimated:  false,
	}

	// Try compressing the image, decreasing quality by 5 each time until we find a smaller size or reach quality - 10
	for {
		compressOpts.Quality = currentQuality
		compressedImage, err = CompressImage(imageBytes, compressOpts)
		if err != nil {
			return nil, currentQuality, fmt.Errorf("failed to compress image: %w", err)
		}

		if len(compressedImage.Bytes) < options.OriginalImageSize {
			return compressedImage, currentQuality, nil // Return the first compressed image that is smaller than the original
		}

		if currentQuality < options.InitialQuality-10 { // Stop if we've decreased quality by 10
			// If no compression was better, return the original image
			return nil, currentQuality, fmt.Errorf("could not compress image into smaller size than original")
		}

		currentQuality -= 5 // Decrease quality by 5 and try again
	}
}

// Compress to webp and jpeg concurrently using goroutines
func CompressImageToBestFormat(imageBytes []byte, options CompressImageToBestFormatOptions) (*CompressImageResult, error) {
	type result struct {
		resp *CompressImageResult
		err  error
	}

	webpCh := make(chan result, 1)
	jpegCh := make(chan result, 1)

	go func() {
		webpImageBytes, errWebp := CompressImage(imageBytes, CompressImageOptions{
			Format:      "webp",
			InputFormat: options.InputFormat,
			Greyscale:   options.Greyscale,
			Quality:     options.Quality,
			IsAnimated:  false,
		})
		webpCh <- result{resp: webpImageBytes, err: errWebp}
	}()

	go func() {
		jpegImageBytes, errJpeg := CompressImage(imageBytes, CompressImageOptions{
			Format:      "jpeg",
			InputFormat: options.InputFormat,
			Greyscale:   options.Greyscale,
			Quality:     options.Quality,
			IsAnimated:  false,
		})
		jpegCh <- result{resp: jpegImageBytes, err: errJpeg}
	}()

	var webpResp, jpegResp *CompressImageResult
	var errWebp, errJpeg error

	for range 2 {
		select {
		case res := <-webpCh:
			webpResp = res.resp
			errWebp = res.err
		case res := <-jpegCh:
			jpegResp = res.resp
			errJpeg = res.err
		}
	}

	if errWebp != nil || errJpeg != nil {
		return nil, fmt.Errorf("failed to compress image:\n\t%w\n\t%w", errWebp, errJpeg)
	}

	originalSize := len(imageBytes)

	webpSize := 0
	if webpResp != nil {
		webpSize = len(webpResp.Bytes)
	}

	jpegSize := 0
	if jpegResp != nil {
		jpegSize = len(jpegResp.Bytes)
	}

	if (webpResp != nil && webpSize < originalSize) || (jpegResp != nil && jpegSize < originalSize) {
		if webpResp != nil && (jpegResp == nil || webpSize < jpegSize) {
			return webpResp, nil
		}
		return jpegResp, nil
	}
	return nil, fmt.Errorf("could not compress image into smaller size than original")
}
