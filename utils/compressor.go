package utils

import (
	"fmt"

	"github.com/energypatrikhu/bandwidth-hero-proxy-go/vips"
)

func CompressImage(imageBytes []byte, imageFormat string, format string, greyscale bool, quality int) (*CompressedImageResponse, error) {
	loadOptions := &vips.LoadOptions{}
	loadOptions.N = -1              // Load all pages/frames
	loadOptions.FailOnError = false // Do not fail on error

	vipsImage, err := vips.NewImageFromBuffer(imageBytes, loadOptions)

	// Fallback to specific loaders if generic loader fails
	if err != nil {
		switch imageFormat {
		case "image/jpeg", "image/jpg":
			vipsImage, err = vips.NewJpegloadBuffer(imageBytes, &vips.JpegloadBufferOptions{
				FailOn: vips.FailOnNone,
			})
		case "image/webp":
			vipsImage, err = vips.NewWebploadBuffer(imageBytes, &vips.WebploadBufferOptions{
				FailOn: vips.FailOnNone,
			})
		case "image/png":
			vipsImage, err = vips.NewPngloadBuffer(imageBytes, &vips.PngloadBufferOptions{
				FailOn: vips.FailOnNone,
			})
		}

		if err != nil {
			return nil, fmt.Errorf("failed to create image from buffer: %w", err)
		}
	}
	defer vipsImage.Close()

	vipsImage.RemoveICCProfile()

	if greyscale {
		vipsImage.Colourspace(vips.InterpretationBW, nil)
	}

	var compressedData []byte

	switch format {
	case "webp":
		compressedData, err = vipsImage.WebpsaveBuffer(&vips.WebpsaveBufferOptions{
			Q:        quality,
			Lossless: false,
			Keep:     vips.KeepNone,
			Effort:   6,
		})
	case "jpeg":
		compressedData, err = vipsImage.JpegsaveBuffer(&vips.JpegsaveBufferOptions{
			Q:                  quality,
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
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to export image buffer: %w", err)
	}

	return &CompressedImageResponse{Data: compressedData, Format: format}, nil
}

func CompressImageWithAutoQualityDecrement(imageBytes []byte, imageFormat string, format string, greyscale bool, quality int, originalSize int) (*CompressedImageResponse, int, error) {
	currentQuality := quality
	var compressedImg *CompressedImageResponse
	var err error

	// Try compressing the image, decreasing quality by 5 each time until we find a smaller size or reach quality - 10
	for {
		compressedImg, err = CompressImage(imageBytes, imageFormat, format, greyscale, currentQuality)
		if err != nil {
			return nil, currentQuality, fmt.Errorf("failed to compress image: %w", err)
		}

		if len(compressedImg.Data) < originalSize {
			return compressedImg, currentQuality, nil // Return the first compressed image that is smaller than the original
		}

		if currentQuality < quality-10 { // Stop if we've decreased quality by 10
			// If no compression was better, return the original image
			return nil, currentQuality, fmt.Errorf("could not compress image into smaller size than original")
		}

		currentQuality -= 5 // Decrease quality by 5 and try again
	}
}

func CompressImageToBestFormat(imageBytes []byte, imageFormat string, greyscale bool, quality int) (*CompressedImageResponse, error) {
	// Compress to webp and jpeg concurrently using goroutines
	type result struct {
		resp *CompressedImageResponse
		err  error
	}

	webpCh := make(chan result)
	jpegCh := make(chan result)

	go func() {
		webpImg, errWebp := CompressImage(imageBytes, imageFormat, "webp", greyscale, quality)
		if errWebp == nil {
			webpCh <- result{resp: webpImg, err: nil}
		} else {
			webpCh <- result{resp: nil, err: errWebp}
		}
	}()

	go func() {
		jpegImg, errJpeg := CompressImage(imageBytes, imageFormat, "jpeg", greyscale, quality)
		if errJpeg == nil {
			jpegCh <- result{resp: jpegImg, err: nil}
		} else {
			jpegCh <- result{resp: nil, err: errJpeg}
		}
	}()

	var webpResp, jpegResp *CompressedImageResponse
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
		webpSize = len(webpResp.Data)
	}

	jpegSize := 0
	if jpegResp != nil {
		jpegSize = len(jpegResp.Data)
	}

	if (webpResp != nil && webpSize < originalSize) || (jpegResp != nil && jpegSize < originalSize) {
		if webpResp != nil && (jpegResp == nil || webpSize < jpegSize) {
			return webpResp, nil
		}
		return jpegResp, nil
	}
	return nil, fmt.Errorf("could not compress image into smaller size than original")
}
