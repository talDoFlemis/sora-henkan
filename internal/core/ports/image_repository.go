package ports

import (
	"context"
	"errors"

	"github.com/taldoflemis/sora-henkan/internal/core/domain/images"
)

type ImageRepository interface {
	CreateNewImage(ctx context.Context, image *images.Image) error
	UpdateImage(ctx context.Context, image *images.Image) error
	DeleteImage(ctx context.Context, id string) error
	FindImageByID(ctx context.Context, id string) (*images.Image, error)
	FindAllImages(ctx context.Context, req *images.ListImagesRequest) (*images.ListImagesResponse, error)
}

// ImageMetadataRepository is responsible for storing image metadata for fast querying
// This is separate from the main ImageRepository to allow different storage backends
type ImageMetadataRepository interface {
	SaveMetadata(ctx context.Context, image *images.Image) error
	UpdateMetadata(ctx context.Context, image *images.Image) error
	DeleteMetadata(ctx context.Context, id string) error
	GetMetadata(ctx context.Context, id string) (*images.ImageMetadata, error)
}

var ErrImageNotFound = errors.New("image not found")
var ErrMetadataNotFound = errors.New("metadata not found")
