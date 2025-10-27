package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/taldoflemis/sora-henkan/internal/core/domain/images"
	"github.com/taldoflemis/sora-henkan/internal/core/ports"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("")

// imageModel represents the database persistence model for images
type imageModel struct {
	ID                    uuid.UUID       `db:"id"`
	OriginalImageURL      string          `db:"original_image_url"`
	ObjectStorageImageKey string          `db:"object_storage_image_key"`
	MimeType              string          `db:"mime_type"`
	Status                string          `db:"status"`
	TransformedImageKey   string          `db:"transformed_image_key"`
	Checksum              string          `db:"checksum"`
	Transformations       json.RawMessage `db:"transformations"`
	UpdatedAt             time.Time       `db:"updated_at"`
	CreatedAt             time.Time       `db:"created_at"`
}

// toDomain converts a persistence model to domain model
func (m *imageModel) toDomain() (*images.Image, error) {
	var transformations images.TransformationList
	if len(m.Transformations) > 0 && string(m.Transformations) != "null" {
		err := json.Unmarshal(m.Transformations, &transformations)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal transformations: %w", err)
		}
	}

	return &images.Image{
		ID:                    m.ID,
		OriginalImageURL:      m.OriginalImageURL,
		ObjectStorageImageKey: m.ObjectStorageImageKey,
		MimeType:              m.MimeType,
		Status:                m.Status,
		TransformedImageKey:   m.TransformedImageKey,
		Checksum:              m.Checksum,
		Transformations:       transformations,
		UpdatedAt:             m.UpdatedAt,
		CreatedAt:             m.CreatedAt,
	}, nil
}

// fromDomain converts a domain model to persistence model
func fromDomain(img *images.Image) (*imageModel, error) {
	transformationsJSON, err := json.Marshal(img.Transformations)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformations: %w", err)
	}

	return &imageModel{
		ID:                    img.ID,
		OriginalImageURL:      img.OriginalImageURL,
		ObjectStorageImageKey: img.ObjectStorageImageKey,
		MimeType:              img.MimeType,
		Status:                img.Status,
		TransformedImageKey:   img.TransformedImageKey,
		Checksum:              img.Checksum,
		Transformations:       transformationsJSON,
		UpdatedAt:             img.UpdatedAt,
		CreatedAt:             img.CreatedAt,
	}, nil
}

type PostgresImageRepository struct {
	pool *pgxpool.Pool
}

var _ ports.ImageRepository = (*PostgresImageRepository)(nil)

// NewPostgresImageRepository creates a new PostgreSQL image repository
func NewPostgresImageRepository(pool *pgxpool.Pool) *PostgresImageRepository {
	return &PostgresImageRepository{
		pool: pool,
	}
}

// CreateNewImage implements ports.ImageRepository.
func (p *PostgresImageRepository) CreateNewImage(ctx context.Context, image *images.Image) error {
	ctx, span := tracer.Start(ctx, "PostgresImageRepository.CreateNewImage")
	defer span.End()

	model, err := fromDomain(image)
	if err != nil {
		return fmt.Errorf("failed to convert to persistence model: %w", err)
	}

	query := `
		INSERT INTO images (
			id, original_image_url, object_storage_image_key, mime_type, status,
			transformed_image_key, checksum, transformations, updated_at, created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
	`

	_, err = p.pool.Exec(ctx, query,
		model.ID,
		model.OriginalImageURL,
		model.ObjectStorageImageKey,
		model.MimeType,
		model.Status,
		model.TransformedImageKey,
		model.Checksum,
		model.Transformations,
		model.UpdatedAt,
		model.CreatedAt,
	)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to create image: %w", err)
	}

	span.SetAttributes(attribute.String("image.id", image.ID.String()))
	return nil
}

// DeleteImage implements ports.ImageRepository.
func (p *PostgresImageRepository) DeleteImage(ctx context.Context, id string) error {
	ctx, span := tracer.Start(ctx, "PostgresImageRepository.DeleteImage")
	defer span.End()

	imageID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid UUID: %w", err)
	}

	query := `DELETE FROM images WHERE id = $1`

	result, err := p.pool.Exec(ctx, query, imageID)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to delete image: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ports.ErrImageNotFound
	}

	span.SetAttributes(attribute.String("image.id", id))
	return nil
}

// FindAllImages implements ports.ImageRepository.
func (p *PostgresImageRepository) FindAllImages(ctx context.Context, req *images.ListImagesRequest) (*images.ListImagesResponse, error) {
	ctx, span := tracer.Start(ctx, "PostgresImageRepository.FindAllImages")
	defer span.End()

	offset := (req.Page - 1) * req.Limit

	query := `
		SELECT id, original_image_url, object_storage_image_key, mime_type, status,
		       transformed_image_key, checksum, transformations, updated_at, created_at
		FROM images
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := p.pool.Query(ctx, query, req.Limit, offset)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to query images: %w", err)
	}
	defer rows.Close()

	var imagesList []images.Image
	for rows.Next() {
		var model imageModel
		err := rows.Scan(
			&model.ID,
			&model.OriginalImageURL,
			&model.ObjectStorageImageKey,
			&model.MimeType,
			&model.Status,
			&model.TransformedImageKey,
			&model.Checksum,
			&model.Transformations,
			&model.UpdatedAt,
			&model.CreatedAt,
		)
		if err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("failed to scan image row: %w", err)
		}

		domainImage, err := model.toDomain()
		if err != nil {
			span.RecordError(err)
			return nil, fmt.Errorf("failed to convert to domain: %w", err)
		}

		imagesList = append(imagesList, *domainImage)
	}

	if err := rows.Err(); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Get total count
	var count int
	countQuery := `SELECT COUNT(*) FROM images`
	err = p.pool.QueryRow(ctx, countQuery).Scan(&count)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get count: %w", err)
	}

	span.SetAttributes(
		attribute.Int("page", req.Page),
		attribute.Int("limit", req.Limit),
		attribute.Int("count", count),
	)

	return &images.ListImagesResponse{
		Page:  req.Page,
		Limit: req.Limit,
		Count: count,
		Data:  imagesList,
	}, nil
}

// FindImageByID implements ports.ImageRepository.
func (p *PostgresImageRepository) FindImageByID(ctx context.Context, id string) (*images.Image, error) {
	ctx, span := tracer.Start(ctx, "PostgresImageRepository.FindImageByID")
	defer span.End()

	imageID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %w", err)
	}

	query := `
		SELECT id, original_image_url, object_storage_image_key, mime_type, status,
		       transformed_image_key, checksum, transformations, updated_at, created_at
		FROM images
		WHERE id = $1
	`

	var model imageModel
	err = p.pool.QueryRow(ctx, query, imageID).Scan(
		&model.ID,
		&model.OriginalImageURL,
		&model.ObjectStorageImageKey,
		&model.MimeType,
		&model.Status,
		&model.TransformedImageKey,
		&model.Checksum,
		&model.Transformations,
		&model.UpdatedAt,
		&model.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ports.ErrImageNotFound
		}
		span.RecordError(err)
		return nil, fmt.Errorf("failed to query image: %w", err)
	}

	domainImage, err := model.toDomain()
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to convert to domain: %w", err)
	}

	span.SetAttributes(attribute.String("image.id", id))
	return domainImage, nil
}

// UpdateImage implements ports.ImageRepository.
func (p *PostgresImageRepository) UpdateImage(ctx context.Context, image *images.Image) error {
	ctx, span := tracer.Start(ctx, "PostgresImageRepository.UpdateImage")
	defer span.End()

	model, err := fromDomain(image)
	if err != nil {
		return fmt.Errorf("failed to convert to persistence model: %w", err)
	}

	query := `
		UPDATE images
		SET original_image_url = $2,
		    object_storage_image_key = $3,
		    mime_type = $4,
		    status = $5,
		    transformed_image_key = $6,
		    checksum = $7,
		    transformations = $8,
		    updated_at = $9
		WHERE id = $1
	`

	result, err := p.pool.Exec(ctx, query,
		model.ID,
		model.OriginalImageURL,
		model.ObjectStorageImageKey,
		model.MimeType,
		model.Status,
		model.TransformedImageKey,
		model.Checksum,
		model.Transformations,
		time.Now(),
	)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to update image: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ports.ErrImageNotFound
	}

	span.SetAttributes(attribute.String("image.id", image.ID.String()))
	return nil
}
