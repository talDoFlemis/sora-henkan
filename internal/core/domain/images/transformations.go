package images

import (
	"fmt"
	"strings"
)

// TransformationRequest represents one transformation step from the API payload.
type TransformationRequest struct {
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
}

// ResizeConfig holds configuration for resizing transformation.
type ResizeConfig struct {
	Width  int `json:"width" validate:"required,gte=1"`
	Height int `json:"height" validate:"required,gte=1"`
}

// GrayscaleConfig holds configuration for grayscale transformation.
type GrayscaleConfig struct {
	// No configuration needed for basic grayscale
	// Can be extended in the future
}

// TrimConfig holds configuration for trim/crop transformation.
type TrimConfig struct {
	// Threshold for edge detection (0-255)
	Threshold float64 `json:"threshold" validate:"gte=0,lte=255"`
}

// BlurConfig holds configuration for blur transformation.
type BlurConfig struct {
	Sigma float64 `json:"sigma" validate:"required,gt=0"`
}

// RotateConfig holds configuration for rotation transformation.
type RotateConfig struct {
	Angle int `json:"angle" validate:"required,oneof=90 180 270"`
}

// ParseTransformations parses a slice of transformation requests.
// This is a utility function that can be used to validate transformation names.
func ParseTransformations(requests []TransformationRequest) ([]string, error) {
	names := make([]string, 0, len(requests))
	for _, req := range requests {
		normalizedName := strings.ToLower(strings.TrimSpace(req.Name))
		if normalizedName == "" {
			return nil, fmt.Errorf("transformation name cannot be empty")
		}
		names = append(names, normalizedName)
	}
	return names, nil
}
