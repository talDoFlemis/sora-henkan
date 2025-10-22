package images

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID                     uuid.UUID
	OriginalImageURL       string
	ObjectStorageImageKey  string
	MimeType               string
	Status                 string
	TransformedImageKey    string
	Checksum               string
	TransformationsApplied []string
	UpdatedAt              time.Time
	CreatedAt              time.Time
}
