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

// Scale implements ports.ImageScaler.
func (v *VipsImageProcessor) Scale(ctx context.Context, image io.Reader, targetHeight int, targetWidth int) ([]byte, error) {
	slog.DebugContext(ctx, "Scaling image", slog.Int("targetHeight", targetHeight), slog.Int("targetWidth", targetWidth))

	buf, err := io.ReadAll(image)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read image", slog.Any("err", err))
		return nil, err
	}

	imageRef, err := vips.NewThumbnailFromBuffer(buf, targetWidth, targetHeight, vips.InterestingCentre)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create thumbnail", slog.Any("err", err))
		return nil, err
	}

	output, _, err := imageRef.ExportNative()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to export thumbnail", slog.Any("err", err))
		return nil, err
	}

	return output, nil
}
