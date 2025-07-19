package utils

import (
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

func CompressImage(imgData []byte, format string, greyscale bool, quality int) (*CompressedImageResponse, error) {
	importParams := &vips.ImportParams{}
	importParams.NumPages.Set(-1)       // Set to -1 to ensure all pages are processed (animated webp support)
	importParams.FailOnError.Set(false) // Do not fail on error, allow processing to continue

	vipsImage, err := vips.LoadImageFromBuffer(imgData, importParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create image from buffer: %w", err)
	}
	defer vipsImage.Close()

	vipsImage.RemoveICCProfile()

	if greyscale {
		vipsImage.ToColorSpace(vips.InterpretationBW)
	}

	var compressedData []byte

	switch format {
	case "webp":
		compressedData, _, err = vipsImage.ExportWebp(&vips.WebpExportParams{
			Quality:         quality,
			Lossless:        false,
			StripMetadata:   true,
			ReductionEffort: 6,
		})
	case "jpeg":
		compressedData, _, err = vipsImage.ExportJpeg(&vips.JpegExportParams{
			Quality:            quality,
			OptimizeCoding:     true,
			OptimizeScans:      true,
			StripMetadata:      true,
			Interlace:          false,
			SubsampleMode:      vips.VipsForeignSubsampleAuto,
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

func CompressImageToBestFormat(imgData []byte, greyscale bool, quality int) (*CompressedImageResponse, error) {
	// Compress to webp and jpeg concurrently using goroutines
	type result struct {
		resp *CompressedImageResponse
		err  error
	}

	webpCh := make(chan result)
	jpegCh := make(chan result)

	go func() {
		webpImg, errWebp := CompressImage(imgData, "webp", greyscale, quality)
		if errWebp == nil {
			webpCh <- result{resp: webpImg, err: nil}
		} else {
			webpCh <- result{resp: nil, err: errWebp}
		}
	}()

	go func() {
		jpegImg, errJpeg := CompressImage(imgData, "jpeg", greyscale, quality)
		if errJpeg == nil {
			jpegCh <- result{resp: jpegImg, err: nil}
		} else {
			jpegCh <- result{resp: nil, err: errJpeg}
		}
	}()

	var webpResp, jpegResp *CompressedImageResponse
	var errWebp, errJpeg error

	for i := 0; i < 2; i++ {
		select {
		case res := <-webpCh:
			webpResp = res.resp
			errWebp = res.err
		case res := <-jpegCh:
			jpegResp = res.resp
			errJpeg = res.err
		}
	}

	if errWebp != nil {
		return nil, fmt.Errorf("failed to compress image to webp: %w", errWebp)
	}
	if errJpeg != nil {
		return nil, fmt.Errorf("failed to compress image to jpeg: %w", errJpeg)
	}

	if (webpResp != nil && len(webpResp.Data) < len(imgData)) || (jpegResp != nil && len(jpegResp.Data) < len(imgData)) {
		if webpResp != nil && (jpegResp == nil || len(webpResp.Data) < len(jpegResp.Data)) {
			return webpResp, nil
		}
		return jpegResp, nil
	}
	return &CompressedImageResponse{Data: imgData, Format: ""}, nil // Return original image if no compression is better
}
