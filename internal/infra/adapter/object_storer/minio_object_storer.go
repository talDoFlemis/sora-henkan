package objectstorer

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/taldoflemis/sora-henkan/internal/core/ports"
	"github.com/taldoflemis/sora-henkan/internal/infra/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("")

type MinioObjectStorer struct {
	minioClient *minio.Client
}

var _ ports.ObjectStorer = (*MinioObjectStorer)(nil)

// NewMinioObjectStorer creates a new MinioObjectStorer instance
func NewMinioObjectStorer(minioClient *minio.Client) *MinioObjectStorer {
	return &MinioObjectStorer{
		minioClient: minioClient,
	}
}

// Delete implements ports.ObjectStorer.
func (m *MinioObjectStorer) Delete(ctx context.Context, key string, bucket string) error {
	ctx, span := tracer.Start(ctx, "MinioObjectStorer.Delete", trace.WithAttributes(
		attribute.String("object.key", key),
		attribute.String("object.bucket", bucket),
	))
	defer span.End()

	err := m.minioClient.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete object", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	slog.InfoContext(ctx, "successfully deleted object", slog.String("key", key))

	return nil
}

// Get implements ports.ObjectStorer.
func (m *MinioObjectStorer) Get(ctx context.Context, key string, bucket string) (io.ReadCloser, error) {
	ctx, span := tracer.Start(ctx, "MinioObjectStorer.Get", trace.WithAttributes(
		attribute.String("object.key", key),
		attribute.String("object.bucket", bucket),
	))
	defer span.End()

	object, err := m.minioClient.GetObject(ctx, bucket, key, minio.GetObjectOptions{Checksum: true})
	if err != nil {
		slog.ErrorContext(ctx, "failed to get object", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	slog.InfoContext(ctx, "successfully retrieved object", slog.String("key", key))

	return object, nil
}

// ListKeys implements ports.ObjectStorer.
func (m *MinioObjectStorer) ListKeys(ctx context.Context) ([]string, error) {
	ctx, span := tracer.Start(ctx, "MinioObjectStorer.ListKeys")
	defer span.End()

	buckets, err := m.minioClient.ListBuckets(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list buckets", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return nil, err
	}

	var keys []string
	for _, bucket := range buckets {
		objectCh := m.minioClient.ListObjects(ctx, bucket.Name, minio.ListObjectsOptions{
			Recursive: true,
		})

		for object := range objectCh {
			if object.Err != nil {
				slog.ErrorContext(ctx, "error listing object", slog.Any("err", object.Err))
				telemetry.RegisterSpanError(span, object.Err)
				return nil, object.Err
			}
			keys = append(keys, object.Key)
		}
	}

	span.SetAttributes(attribute.Int("objects.count", len(keys)))
	slog.InfoContext(ctx, "successfully listed keys", slog.Int("count", len(keys)))

	return keys, nil
}

// PresignedUpload implements ports.ObjectStorer.
func (m *MinioObjectStorer) PresignedUpload(ctx context.Context, bucket string) (string, error) {
	ctx, span := tracer.Start(ctx, "MinioObjectStorer.PresignedUpload", trace.WithAttributes(
		attribute.String("object.bucket", bucket),
	))
	defer span.End()

	// Generate a unique key for the upload
	// Note: The interface might need to be updated to accept a key parameter
	// For now, generating a placeholder key
	key := "upload-placeholder"
	
	// Generate presigned URL valid for 15 minutes
	presignedURL, err := m.minioClient.PresignedPutObject(ctx, bucket, key, 15*time.Minute)
	if err != nil {
		slog.ErrorContext(ctx, "failed to generate presigned upload URL", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return "", err
	}

	span.SetAttributes(attribute.String("presigned.url", presignedURL.String()))
	slog.InfoContext(ctx, "successfully generated presigned upload URL", slog.String("url", presignedURL.String()))
	
	return presignedURL.String(), nil
}

// Store implements ports.ObjectStorer.
func (m *MinioObjectStorer) Store(ctx context.Context, key string, bucket string, mineType string, obj io.Reader) error {
	ctx, span := tracer.Start(ctx, "MinioObjectStorer.Store", trace.WithAttributes(
		attribute.String("object.key", key),
		attribute.String("object.bucket", bucket),
		attribute.String("object.mime_type", mineType),
	))
	defer span.End()

	uploadInfo, err := m.minioClient.PutObject(ctx, bucket, key, obj, -1, minio.PutObjectOptions{
		ContentType: mineType,
	})
	if err != nil {
		slog.ErrorContext(ctx, "failed to upload object", slog.Any("err", err))
		telemetry.RegisterSpanError(span, err)
		return err
	}

	span.SetAttributes(attribute.String("object.etag", uploadInfo.ChecksumCRC32C))

	slog.InfoContext(ctx, "successfully uploaded object", slog.String("etag", uploadInfo.ETag))

	return nil
}
