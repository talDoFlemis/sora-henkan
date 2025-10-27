package image

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/taldoflemis/sora-henkan/internal/core/domain/images"
	"github.com/taldoflemis/sora-henkan/internal/core/ports"
)

// TransformerFactory implements the ImageTransformerFactory interface with injected transformers
type TransformerFactory struct {
	resizeTransformer    ports.ImageTransformer
	grayscaleTransformer ports.ImageTransformer
	trimTransformer      ports.ImageTransformer
	blurTransformer      ports.ImageTransformer
	rotateTransformer    ports.ImageTransformer
}

var _ ports.ImageTransformerFactory = (*TransformerFactory)(nil)

func NewTransformerFactory(
	resizeTransformer ports.ImageTransformer,
	grayscaleTransformer ports.ImageTransformer,
	trimTransformer ports.ImageTransformer,
	blurTransformer ports.ImageTransformer,
	rotateTransformer ports.ImageTransformer,
) *TransformerFactory {
	return &TransformerFactory{
		resizeTransformer:    resizeTransformer,
		grayscaleTransformer: grayscaleTransformer,
		trimTransformer:      trimTransformer,
		blurTransformer:      blurTransformer,
		rotateTransformer:    rotateTransformer,
	}
}

// CreateTransformer creates a transformer from a transformation request
func (f *TransformerFactory) CreateTransformer(ctx context.Context, req images.TransformationRequest) (ports.ImageTransformer, error) {
	transformName := strings.ToLower(strings.TrimSpace(req.Name))

	switch transformName {
	case "resize":
		return f.resizeTransformer, nil

	case "grayscale":
		return f.grayscaleTransformer, nil

	case "trim":
		return f.trimTransformer, nil

	case "blur":
		return f.blurTransformer, nil

	case "rotate":
		return f.rotateTransformer, nil

	default:
		return nil, fmt.Errorf("unknown transformation: %s", req.Name)
	}
}

// Pipeline implements the ImagePipelineProcessor interface with injected transformers
type Pipeline struct {
	factory ports.ImageTransformerFactory
}

var _ ports.ImagePipelineProcessor = (*Pipeline)(nil)

func NewPipeline(factory ports.ImageTransformerFactory) *Pipeline {
	return &Pipeline{
		factory: factory,
	}
}

// ValidateTransformations validates all transformations and their configs without processing the image
func (p *Pipeline) ValidateTransformations(ctx context.Context, transformations []images.TransformationRequest) error {
	slog.InfoContext(ctx, "Validating transformations", slog.Int("count", len(transformations)))

	for i, txReq := range transformations {
		slog.DebugContext(ctx, "Validating transformation step",
			slog.Int("step", i+1),
			slog.String("transformation", txReq.Name))

		// Create the transformer to ensure it exists
		transformer, err := p.factory.CreateTransformer(ctx, txReq)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to create transformer",
				slog.Int("step", i+1),
				slog.String("transformation", txReq.Name),
				slog.Any("err", err))
			return fmt.Errorf("invalid transformation at step %d (%s): %w", i+1, txReq.Name, err)
		}

		configMap := txReq.Config
		if configMap == nil {
			configMap = make(map[string]any)
		}

		// Validate config
		if err := transformer.ValidateConfig(ctx, configMap); err != nil {
			slog.ErrorContext(ctx, "Config validation failed",
				slog.Int("step", i+1),
				slog.String("transformation", transformer.Name()),
				slog.Any("err", err))
			return fmt.Errorf("invalid config at step %d (%s): %w", i+1, transformer.Name(), err)
		}

		slog.DebugContext(ctx, "Transformation validation successful",
			slog.Int("step", i+1),
			slog.String("transformation", transformer.Name()))
	}

	slog.InfoContext(ctx, "All transformations validated successfully", slog.Int("total_steps", len(transformations)))
	return nil
}

// ProcessPipeline processes an image through a pipeline of transformations
func (p *Pipeline) ProcessPipeline(ctx context.Context, image io.Reader, transformations []images.TransformationRequest) ([]byte, error) {
	slog.InfoContext(ctx, "Starting image transformation pipeline", slog.Int("steps", len(transformations)))

	// Read initial image data
	initialData, err := io.ReadAll(image)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to read initial image data", slog.Any("err", err))
		return nil, fmt.Errorf("failed to read initial image data: %w", err)
	}

	// Current image data
	currentData := initialData

	// Build and execute the pipeline
	for i, txReq := range transformations {
		slog.DebugContext(ctx, "Building transformation step",
			slog.Int("step", i+1),
			slog.String("transformation", txReq.Name))

		// Create the transformer
		transformer, err := p.factory.CreateTransformer(ctx, txReq)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to create transformer",
				slog.Int("step", i+1),
				slog.String("transformation", txReq.Name),
				slog.Any("err", err))
			return nil, fmt.Errorf("failed to create transformer at step %d (%s): %w", i+1, txReq.Name, err)
		}

		// Validate config before transformation
		if err := transformer.ValidateConfig(ctx, txReq.Config); err != nil {
			slog.ErrorContext(ctx, "Config validation failed",
				slog.Int("step", i+1),
				slog.String("transformation", transformer.Name()),
				slog.Any("err", err))
			return nil, fmt.Errorf("config validation failed at step %d (%s): %w", i+1, transformer.Name(), err)
		}

		// Apply the transformation
		slog.InfoContext(ctx, "Applying transformation",
			slog.Int("step", i+1),
			slog.String("transformation", transformer.Name()))

		transformedData, err := transformer.Transform(ctx, currentData, txReq.Config)
		if err != nil {
			slog.ErrorContext(ctx, "Transformation failed",
				slog.Int("step", i+1),
				slog.String("transformation", transformer.Name()),
				slog.Any("err", err))
			return nil, fmt.Errorf("transformation failed at step %d (%s): %w", i+1, transformer.Name(), err)
		}

		// Update current data for next step
		currentData = transformedData

		slog.DebugContext(ctx, "Transformation completed successfully",
			slog.Int("step", i+1),
			slog.String("transformation", transformer.Name()),
			slog.Int("output_size", len(transformedData)))
	}

	slog.InfoContext(ctx, "Image transformation pipeline completed successfully",
		slog.Int("total_steps", len(transformations)),
		slog.Int("final_size", len(currentData)))

	return currentData, nil
}
