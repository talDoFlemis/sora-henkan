package ports

import (
	"context"
	"io"
)

type ObjectStorer interface {
	Store(ctx context.Context, key string, bucket string, mimeType string, obj io.Reader) error
	Get(ctx context.Context, key string, bucket string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string, bucket string) error
	ListKeys(ctx context.Context) ([]string, error)
	PresignedUpload(ctx context.Context, bucket string) (string, error)
}
