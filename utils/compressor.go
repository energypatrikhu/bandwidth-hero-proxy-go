package utils

import (
	"fmt"
	"slices"

	"github.com/energypatrikhu/bandwidth-hero-proxy-go/vips"
)

func CompressImage(imageBytes []byte, options *CompressImageOptions) (*CompressImageResult, error) {
	loadOptions := &vips.LoadOptions{
		FailOnError: false,
	}
	if options.IsAnimated {
		loadOptions.N = -1 // Load all frames for animated images
	}
	if slices.Contains(FormatsSupportingVipsUnlimited, options.InputFormat) {
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

	var compressedData []byte

	switch options.Format {
	case "webp":
		compressedData, vipsError = vipsImage.WebpsaveBuffer(&vips.WebpsaveBufferOptions{
			Q:        options.Quality,
			Lossless: false,
			Keep:     vips.KeepNone,
			Effort:   6,
		})
	case "jpeg":
		compressedData, vipsError = vipsImage.JpegsaveBuffer(&vips.JpegsaveBufferOptions{
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
	default:
		return nil, fmt.Errorf("unsupported format: %s", options.Format)
	}

	if vipsError != nil {
		return nil, fmt.Errorf("failed to export image buffer: %w", vipsError)
	}

	return &CompressImageResult{Bytes: compressedData, Format: options.Format}, nil
}

func CompressImageWithAutoQualityDecrement(imageBytes []byte, options *CompressImageWithAutoQualityDecrementOptions) (*CompressImageResult, int, error) {
	currentQuality := options.InitialQuality
	var compressedImage *CompressImageResult
	var err error

	// Try compressing the image, decreasing quality by 5 each time until we find a smaller size or reach quality - 10
	for {
		compressedImage, err = CompressImage(imageBytes, &CompressImageOptions{
			InputFormat: options.InputFormat,
			Format:      options.Format,
			Greyscale:   options.Greyscale,
			Quality:     currentQuality,
		})
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

func CompressImageToBestFormat(imageBytes []byte, options *CompressImageToBestFormatOptions) (*CompressImageResult, error) {
	// Compress to webp and jpeg concurrently using goroutines
	type result struct {
		resp *CompressImageResult
		err  error
	}

	webpCh := make(chan result)
	jpegCh := make(chan result)

	go func() {
		// webpImg, errWebp := CompressImage(imageBytes, false, "webp", greyscale, quality)
		webpImg, errWebp := CompressImage(imageBytes, &CompressImageOptions{
			Format:      "webp",
			InputFormat: options.InputFormat,
			Greyscale:   options.Greyscale,
			Quality:     options.Quality,
		})
		if errWebp == nil {
			webpCh <- result{resp: webpImg, err: nil}
		} else {
			webpCh <- result{resp: nil, err: errWebp}
		}
	}()

	go func() {
		jpegImg, errJpeg := CompressImage(imageBytes, &CompressImageOptions{
			Format:      "jpeg",
			InputFormat: options.InputFormat,
			Greyscale:   options.Greyscale,
			Quality:     options.Quality,
		})
		if errJpeg == nil {
			jpegCh <- result{resp: jpegImg, err: nil}
		} else {
			jpegCh <- result{resp: nil, err: errJpeg}
		}
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
