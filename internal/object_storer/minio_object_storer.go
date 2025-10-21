package objectstorer

import (
	"context"
	"io"

	"github.com/taldoflemis/sora-henkan/internal/ports"
)

type MinioObjectStorer struct{}

var _ ports.ObjectStorer = (*MinioObjectStorer)(nil)

// Delete implements ports.ObjectStorer.
func (m *MinioObjectStorer) Delete(ctx context.Context, key string, bucket string) error {
	panic("unimplemented")
}

// Get implements ports.ObjectStorer.
func (m *MinioObjectStorer) Get(ctx context.Context, key string, bucket string) (io.ReadCloser, error) {
	panic("unimplemented")
}

// ListKeys implements ports.ObjectStorer.
func (m *MinioObjectStorer) ListKeys(ctx context.Context) ([]string, error) {
	panic("unimplemented")
}

// PresignedUpload implements ports.ObjectStorer.
func (m *MinioObjectStorer) PresignedUpload(ctx context.Context, bucket string) (string, error) {
	panic("unimplemented")
}

// Store implements ports.ObjectStorer.
func (m *MinioObjectStorer) Store(ctx context.Context, key string, bucket string, obj io.Reader) error {
	panic("unimplemented")
}
