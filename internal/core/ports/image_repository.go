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

var ErrImageNotFound = errors.New("image not found")
