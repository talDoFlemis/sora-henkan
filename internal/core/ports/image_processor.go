package ports

import (
	"context"
	"io"
)

type ImageScaler interface {
	Scale(ctx context.Context, image io.Reader, targetHeight int, targetWidth int) ([]byte, error)
}
