package image

import (
	"context"
	"io"
	"log/slog"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/taldoflemis/sora-henkan/internal/core/ports"
)

type VipsImageProcessor struct{}

var _ ports.ImageScaler = (*VipsImageProcessor)(nil)

func NewVipsImageProcessor() *VipsImageProcessor {
	return &VipsImageProcessor{}
}

// Scale implements ports.ImageScaler.
func (v *VipsImageProcessor) Scale(ctx context.Context, image io.Reader, targetHeight int, targetWidth int) ([]byte, error) {
	slog.DebugContext(ctx, "Scaling image", slog.Int("targetHeight", targetHeight), slog.Int("targetWidth", targetWidth))

	buf, err := io.ReadAll(image)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read image", slog.Any("err", err))
		return nil, err
	}

	// Load image to check original dimensions
	imageRef, err := vips.NewImageFromBuffer(buf)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load image", slog.Any("err", err))
		return nil, err
	}
	defer imageRef.Close()

	originalWidth := imageRef.Width()
	originalHeight := imageRef.Height()

	// Determine if we're upscaling
	isUpscaling := targetWidth > originalWidth || targetHeight > originalHeight

	if isUpscaling {
		slog.DebugContext(ctx, "Upscaling detected, using Lanczos3 interpolation",
			slog.Int("originalWidth", originalWidth),
			slog.Int("originalHeight", originalHeight))

		// Calculate scale factor
		scaleX := float64(targetWidth) / float64(originalWidth)
		scaleY := float64(targetHeight) / float64(originalHeight)

		// Use the smaller scale to maintain aspect ratio
		scale := scaleX
		if scaleY < scaleX {
			scale = scaleY
		}

		// Resize using Lanczos3 interpolation for better upscaling quality
		err = imageRef.ResizeWithVScale(scale, scale, vips.KernelLanczos3)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to resize image with Lanczos3", slog.Any("err", err))
			return nil, err
		}
	} else {
		// For downscaling, use the thumbnail method which is optimized
		imageRef.Close() // Close the previous reference
		imageRef, err = vips.NewThumbnailFromBuffer(buf, targetWidth, targetHeight, vips.InterestingCentre)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to create thumbnail", slog.Any("err", err))
			return nil, err
		}
	}

	output, _, err := imageRef.ExportNative()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to export image", slog.Any("err", err))
		return nil, err
	}

	return output, nil
}
