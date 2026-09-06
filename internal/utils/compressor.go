package utils

import (
	"fmt"
	"strings"

	"github.com/energypatrikhu/bandwidth-hero-proxy-go/third_party/vips"
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

	if ConfigInstance.PassthroughAnimated && vipsImage.Pages() > 1 {
		return nil, fmt.Errorf("failed to use image buffer, animated images are disabled")
	}

	vipsImage.RemoveICCProfile()

	if options.Grayscale {
		if grayErr := vipsImage.Colourspace(vips.InterpretationBW, nil); grayErr != nil {
			return nil, fmt.Errorf("failed to convert image to grayscale: %w", grayErr)
		}
	}

	var compressedImageBytes []byte

	if ConfigInstance.TargetFormat != "" {
		options.Format = ConfigInstance.TargetFormat
	}

	switch options.Format {
	case "webp":
		compressedImageBytes, vipsError = vipsImage.WebpsaveBuffer(&vips.WebpsaveBufferOptions{
			Q:    options.Quality,
			Keep: vips.KeepNone,

			Lossless:       ConfigInstance.WebpLossless,
			Effort:         ConfigInstance.WebpEffort,
			SmartSubsample: ConfigInstance.WebpSmartSubsample,
			SmartDeblock:   ConfigInstance.WebpSmartDeblock,
			Passes:         ConfigInstance.WebpPasses,
		})
	case "jpeg":
		compressedImageBytes, vipsError = vipsImage.JpegsaveBuffer(&vips.JpegsaveBufferOptions{
			Q:             options.Quality,
			Keep:          vips.KeepNone,
			SubsampleMode: vips.SubsampleOn,

			OptimizeCoding:     ConfigInstance.JpegOptimizeCoding,
			OptimizeScans:      ConfigInstance.JpegOptimizeScans,
			Interlace:          ConfigInstance.JpegInterlace,
			TrellisQuant:       ConfigInstance.JpegTrellisQuant,
			OvershootDeringing: ConfigInstance.JpegOvershootDeringing,
			QuantTable:         ConfigInstance.JpegQuantTable,
		})
	case "jxl":
		compressedImageBytes, vipsError = vipsImage.JxlsaveBuffer(&vips.JxlsaveBufferOptions{
			Q:    options.Quality,
			Keep: vips.KeepNone,

			Tier:     ConfigInstance.JxlTier,
			Effort:   ConfigInstance.JxlEffort,
			Lossless: ConfigInstance.JxlLossless,
		})
	default:
		return nil, fmt.Errorf("unsupported output format: %s", options.Format)
	}

	if vipsError != nil {
		return nil, fmt.Errorf("failed to export image buffer: %w", vipsError)
	}

	return &CompressImageResult{Bytes: compressedImageBytes, Format: options.Format}, nil
}

func CompressImageWithAutoQualityReducer(imageBytes []byte, options CompressImageWithAutoQualityReducerOptions) (*CompressImageResult, int, error) {
	currentQuality := options.InitialQuality
	var compressedImage *CompressImageResult
	var err error

	// Reuse options struct to reduce allocations
	compressOpts := CompressImageOptions{
		InputFormat: options.InputFormat,
		Format:      options.Format,
		Grayscale:   options.Grayscale,
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

		if currentQuality < options.InitialQuality-10 || currentQuality <= 5 {
			// Stop if we've decreased quality by 10, or we're already at a floor
			// where further decrements would hit invalid quality values (0 or negative)
			// If no compression was better, return the original image
			return nil, currentQuality, fmt.Errorf("could not compress image into smaller size than original")
		}

		currentQuality -= 5 // Decrease quality by 5 and try again
	}
}

// Compress to all requested formats concurrently and return the smallest result
func CompressImageToSmallestFormat(imageBytes []byte, options CompressImageToSmallestFormatOptions) (*CompressImageResult, error) {
	formats := ConfigInstance.AutoSelectFormatList
	if len(formats) == 0 {
		formats = []string{"jxl", "webp", "jpeg"} // defaults
	}

	type result struct {
		resp *CompressImageResult
		err  error
	}

	resultCh := make(chan result, len(formats))

	for _, format := range formats {
		go func() {
			resp, err := CompressImage(imageBytes, CompressImageOptions{
				Format:      format,
				InputFormat: options.InputFormat,
				Grayscale:   options.Grayscale,
				Quality:     options.Quality,
				IsAnimated:  false,
			})
			resultCh <- result{resp: resp, err: err}
		}()
	}

	var errs []error
	var best *CompressImageResult
	originalSize := len(imageBytes)

	for range formats {
		res := <-resultCh
		if res.err != nil {
			errs = append(errs, res.err)
			continue
		}
		if res.resp == nil {
			continue
		}

		if len(res.resp.Bytes) >= originalSize {
			continue
		}
		if best == nil || len(res.resp.Bytes) < len(best.Bytes) {
			best = res.resp
		}
	}

	if best != nil {
		return best, nil
	}

	if len(errs) > 0 {
		var errStr strings.Builder
		errStr.WriteString("failed to compress image:")
		for _, e := range errs {
			errStr.WriteString("\n\t" + e.Error())
		}
		return nil, fmt.Errorf("%s", errStr.String())
	}
	return nil, fmt.Errorf("could not compress image into smaller size than original")
}
