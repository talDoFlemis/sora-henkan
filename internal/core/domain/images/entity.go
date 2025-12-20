package images

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TransformationList is a wrapper for storing transformations in the database
type TransformationList []TransformationRequest

// Value implements the driver.Valuer interface for database storage
func (t TransformationList) Value() (driver.Value, error) {
	if t == nil {
		return nil, nil
	}
	return json.Marshal(t)
}

// Scan implements the sql.Scanner interface for database retrieval
func (t *TransformationList) Scan(value interface{}) error {
	if value == nil {
		*t = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal transformations: expected []byte, got %T", value)
	}

	return json.Unmarshal(bytes, t)
}

type Image struct {
	ID                    uuid.UUID          `json:"id"`
	OriginalImageURL      string             `json:"original_image_url"`
	ObjectStorageImageKey string             `json:"object_storage_image_key"`
	MimeType              string             `json:"mime_type"`
	Status                string             `json:"status"`
	TransformedImageKey   string             `json:"transformed_image_key"`
	Checksum              string             `json:"checksum"`
	ErrorMessage          string             `json:"error_message,omitempty"`
	Transformations       TransformationList `json:"transformations"`
	UpdatedAt             time.Time          `json:"updated_at"`
	CreatedAt             time.Time          `json:"created_at"`
}

// ImageMetadata represents metadata stored in DynamoDB for fast querying
type ImageMetadata struct {
	ID                    string    `json:"id"`
	OriginalImageURL      string    `json:"original_image_url"`
	ObjectStorageImageKey string    `json:"object_storage_image_key"`
	TransformedImageKey   string    `json:"transformed_image_key"`
	MimeType              string    `json:"mime_type"`
	Status                string    `json:"status"`
	Checksum              string    `json:"checksum"`
	ErrorMessage          string    `json:"error_message,omitempty"`
	TransformationCount   int       `json:"transformation_count"`
	UpdatedAt             time.Time `json:"updated_at"`
	CreatedAt             time.Time `json:"created_at"`
}

// ToMetadata converts an Image to ImageMetadata
func (i *Image) ToMetadata() *ImageMetadata {
	return &ImageMetadata{
		ID:                    i.ID.String(),
		OriginalImageURL:      i.OriginalImageURL,
		ObjectStorageImageKey: i.ObjectStorageImageKey,
		TransformedImageKey:   i.TransformedImageKey,
		MimeType:              i.MimeType,
		Status:                i.Status,
		Checksum:              i.Checksum,
		ErrorMessage:          i.ErrorMessage,
		TransformationCount:   len(i.Transformations),
		UpdatedAt:             i.UpdatedAt,
		CreatedAt:             i.CreatedAt,
	}
}
