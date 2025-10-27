package ports

import (
	"context"
	"errors"
	"io"

	"github.com/taldoflemis/sora-henkan/internal/core/domain/images"
)

// ImageTransformer is the common interface for all image transformation steps.
type ImageTransformer interface {
	// Transform applies the transformation to the image data.
	// Returns the transformed image bytes and an error if the transformation fails.
	Transform(ctx context.Context, image []byte, config map[string]any) ([]byte, error)

	// ValidateConfig validates the configuration data for the transformer.
	// Returns an error if the config is invalid.
	ValidateConfig(ctx context.Context, config map[string]any) error

	// Name returns the name of the transformation
	Name() string
}

// ImageTransformerFactory creates image transformers from transformation requests
type ImageTransformerFactory interface {
	CreateTransformer(ctx context.Context, req images.TransformationRequest) (ImageTransformer, error)
}

// ImagePipelineProcessor processes images through a pipeline of transformations
type ImagePipelineProcessor interface {
	ProcessPipeline(ctx context.Context, image io.Reader, transformations []images.TransformationRequest) ([]byte, error)
	ValidateTransformations(ctx context.Context, transformations []images.TransformationRequest) error
}

var ErrUnknownImageTransformer = errors.New("unknown image transformer")
