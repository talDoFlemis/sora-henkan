package images

import (
	"time"

	"github.com/google/uuid"
)

type ScaleTransformation struct {
	Enabled bool `json:"enabled"`
	Height  int  `json:"height" validate:"required_if=Enabled true,gte=1"`
	Width   int  `json:"width" validate:"required_if=Enabled true,gte=1"`
}

type Image struct {
	ID                    uuid.UUID           `json:"id"`
	OriginalImageURL      string              `json:"original_image_url"`
	ObjectStorageImageKey string              `json:"object_storage_image_key"`
	MimeType              string              `json:"mime_type"`
	Status                string              `json:"status"`
	TransformedImageKey   string              `json:"transformed_image_key"`
	Checksum              string              `json:"checksum"`
	ScaleTransformation   ScaleTransformation `json:"scale_transformation"`
	UpdatedAt             time.Time           `json:"updated_at"`
	CreatedAt             time.Time           `json:"created_at"`
}
