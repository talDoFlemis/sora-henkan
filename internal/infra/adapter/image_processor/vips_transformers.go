package image

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/go-playground/validator/v10"
	"github.com/go-viper/mapstructure/v2"
	"github.com/taldoflemis/sora-henkan/internal/core/domain/images"
	"github.com/taldoflemis/sora-henkan/internal/core/ports"
)

var validate = validator.New()

// VipsResizeTransformer implements the resize transformation using VIPS
type VipsResizeTransformer struct{}

var _ ports.ImageTransformer = (*VipsResizeTransformer)(nil)

func NewVipsResizeTransformer() *VipsResizeTransformer {
	return &VipsResizeTransformer{}
}

func (t *VipsResizeTransformer) Name() string {
	return "resize"
}

func (t *VipsResizeTransformer) ValidateConfig(ctx context.Context, config map[string]any) error {
	var cfg images.ResizeConfig
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return fmt.Errorf("failed to decode resize config: %w", err)
	}
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("invalid resize config: %w", err)
	}
	return nil
}

func (t *VipsResizeTransformer) Transform(ctx context.Context, image []byte, config map[string]any) ([]byte, error) {
	// Decode config
	var cfg images.ResizeConfig
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode resize config: %w", err)
	}

	slog.DebugContext(ctx, "Applying resize transformation",
		slog.Int("width", cfg.Width),
		slog.Int("height", cfg.Height))

	// Load image to check original dimensions
	imageRef, err := vips.NewImageFromBuffer(image)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load image", slog.Any("err", err))
		return nil, fmt.Errorf("failed to load image: %w", err)
	}
	defer imageRef.Close()

	originalWidth := imageRef.Width()
	originalHeight := imageRef.Height()

	// Determine if we're upscaling
	isUpscaling := cfg.Width > originalWidth || cfg.Height > originalHeight

	if isUpscaling {
		slog.DebugContext(ctx, "Upscaling detected, using Lanczos3 interpolation",
			slog.Int("originalWidth", originalWidth),
			slog.Int("originalHeight", originalHeight))

		// Calculate scale factor
		scaleX := float64(cfg.Width) / float64(originalWidth)
		scaleY := float64(cfg.Height) / float64(originalHeight)

		// Resize using Lanczos3 interpolation for better upscaling quality
		err = imageRef.ResizeWithVScale(scaleX, scaleY, vips.KernelLanczos3)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to resize image with Lanczos3", slog.Any("err", err))
			return nil, fmt.Errorf("failed to resize image: %w", err)
		}
	} else {
		// For downscaling, use the thumbnail method which is optimized
		imageRef.Close() // Close the previous reference
		imageRef, err = vips.NewThumbnailFromBuffer(image, cfg.Width, cfg.Height, vips.InterestingCentre)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to create thumbnail", slog.Any("err", err))
			return nil, fmt.Errorf("failed to create thumbnail: %w", err)
		}
		defer imageRef.Close()
	}

	output, _, err := imageRef.ExportNative()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to export image", slog.Any("err", err))
		return nil, fmt.Errorf("failed to export image: %w", err)
	}

	return output, nil
}

// VipsGrayscaleTransformer implements grayscale transformation using VIPS
type VipsGrayscaleTransformer struct{}

var _ ports.ImageTransformer = (*VipsGrayscaleTransformer)(nil)

func NewVipsGrayscaleTransformer() *VipsGrayscaleTransformer {
	return &VipsGrayscaleTransformer{}
}

func (t *VipsGrayscaleTransformer) Name() string {
	return "grayscale"
}

func (t *VipsGrayscaleTransformer) ValidateConfig(ctx context.Context, config map[string]any) error {
	// Grayscale has no required config, always valid
	return nil
}

func (t *VipsGrayscaleTransformer) Transform(ctx context.Context, image []byte, config map[string]any) ([]byte, error) {
	slog.DebugContext(ctx, "Applying grayscale transformation")

	imageRef, err := vips.NewImageFromBuffer(image)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load image", slog.Any("err", err))
		return nil, fmt.Errorf("failed to load image: %w", err)
	}
	defer imageRef.Close()

	// Convert to grayscale
	err = imageRef.ToColorSpace(vips.InterpretationBW)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to convert to grayscale", slog.Any("err", err))
		return nil, fmt.Errorf("failed to convert to grayscale: %w", err)
	}

	output, _, err := imageRef.ExportNative()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to export image", slog.Any("err", err))
		return nil, fmt.Errorf("failed to export image: %w", err)
	}

	return output, nil
}

// VipsTrimTransformer implements edge trimming using VIPS
type VipsTrimTransformer struct{}

var _ ports.ImageTransformer = (*VipsTrimTransformer)(nil)

func NewVipsTrimTransformer() *VipsTrimTransformer {
	return &VipsTrimTransformer{}
}

func (t *VipsTrimTransformer) Name() string {
	return "trim"
}

func (t *VipsTrimTransformer) ValidateConfig(ctx context.Context, config map[string]any) error {
	var cfg images.TrimConfig
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return fmt.Errorf("failed to decode trim config: %w", err)
	}
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("invalid trim config: %w", err)
	}
	return nil
}

func (t *VipsTrimTransformer) Transform(ctx context.Context, image []byte, config map[string]any) ([]byte, error) {
	// Decode config with default value
	var cfg images.TrimConfig
	cfg.Threshold = 10.0 // Default
	if len(config) > 0 {
		if err := mapstructure.Decode(config, &cfg); err != nil {
			return nil, fmt.Errorf("failed to decode trim config: %w", err)
		}
	}

	slog.DebugContext(ctx, "Applying trim transformation", slog.Float64("threshold", cfg.Threshold))

	imageRef, err := vips.NewImageFromBuffer(image)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load image", slog.Any("err", err))
		return nil, fmt.Errorf("failed to load image: %w", err)
	}
	defer imageRef.Close()

	// Use VIPS FindTrim to detect and remove borders
	left, top, width, height, err := imageRef.FindTrim(cfg.Threshold, &vips.Color{R: 255, G: 255, B: 255})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to find trim area", slog.Any("err", err))
		return nil, fmt.Errorf("failed to find trim area: %w", err)
	}

	// Extract the trimmed region
	err = imageRef.ExtractArea(left, top, width, height)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to extract trim area", slog.Any("err", err))
		return nil, fmt.Errorf("failed to extract trim area: %w", err)
	}

	output, _, err := imageRef.ExportNative()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to export image", slog.Any("err", err))
		return nil, fmt.Errorf("failed to export image: %w", err)
	}

	return output, nil
}

// VipsBlurTransformer implements blur transformation using VIPS
type VipsBlurTransformer struct{}

var _ ports.ImageTransformer = (*VipsBlurTransformer)(nil)

func NewVipsBlurTransformer() *VipsBlurTransformer {
	return &VipsBlurTransformer{}
}

func (t *VipsBlurTransformer) Name() string {
	return "blur"
}

func (t *VipsBlurTransformer) ValidateConfig(ctx context.Context, config map[string]any) error {
	var cfg images.BlurConfig
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return fmt.Errorf("failed to decode blur config: %w", err)
	}
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("invalid blur config: %w", err)
	}
	return nil
}

func (t *VipsBlurTransformer) Transform(ctx context.Context, image []byte, config map[string]any) ([]byte, error) {
	// Decode config
	var cfg images.BlurConfig
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode blur config: %w", err)
	}

	slog.DebugContext(ctx, "Applying blur transformation", slog.Float64("sigma", cfg.Sigma))

	imageRef, err := vips.NewImageFromBuffer(image)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load image", slog.Any("err", err))
		return nil, fmt.Errorf("failed to load image: %w", err)
	}
	defer imageRef.Close()

	// Apply Gaussian blur
	err = imageRef.GaussianBlur(cfg.Sigma)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to apply blur", slog.Any("err", err))
		return nil, fmt.Errorf("failed to apply blur: %w", err)
	}

	output, _, err := imageRef.ExportNative()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to export image", slog.Any("err", err))
		return nil, fmt.Errorf("failed to export image: %w", err)
	}

	return output, nil
}

// VipsRotateTransformer implements rotation transformation using VIPS
type VipsRotateTransformer struct{}

var _ ports.ImageTransformer = (*VipsRotateTransformer)(nil)

func NewVipsRotateTransformer() *VipsRotateTransformer {
	return &VipsRotateTransformer{}
}

func (t *VipsRotateTransformer) Name() string {
	return "rotate"
}

func (t *VipsRotateTransformer) ValidateConfig(ctx context.Context, config map[string]any) error {
	var cfg images.RotateConfig
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return fmt.Errorf("failed to decode rotate config: %w", err)
	}
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("invalid rotate config: %w", err)
	}
	return nil
}

func (t *VipsRotateTransformer) Transform(ctx context.Context, image []byte, config map[string]any) ([]byte, error) {
	// Decode config
	var cfg images.RotateConfig
	if err := mapstructure.Decode(config, &cfg); err != nil {
		return nil, fmt.Errorf("failed to decode rotate config: %w", err)
	}

	slog.DebugContext(ctx, "Applying rotate transformation", slog.Int("angle", cfg.Angle))

	imageRef, err := vips.NewImageFromBuffer(image)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to load image", slog.Any("err", err))
		return nil, fmt.Errorf("failed to load image: %w", err)
	}
	defer imageRef.Close()

	// Convert angle to VIPS angle enum
	var vipsAngle vips.Angle
	switch cfg.Angle {
	case 90:
		vipsAngle = vips.Angle90
	case 180:
		vipsAngle = vips.Angle180
	case 270:
		vipsAngle = vips.Angle270
	default:
		return nil, fmt.Errorf("unsupported rotation angle: %d", cfg.Angle)
	}

	// Rotate the image
	err = imageRef.Rotate(vipsAngle)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to rotate image", slog.Any("err", err))
		return nil, fmt.Errorf("failed to rotate image: %w", err)
	}

	output, _, err := imageRef.ExportNative()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to export image", slog.Any("err", err))
		return nil, fmt.Errorf("failed to export image: %w", err)
	}

	return output, nil
}
