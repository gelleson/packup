package upload

import (
	"context"
	"io"
)

type Uploader interface {
	Put(ctx context.Context, namespace string, filename string, body io.Reader) (string, error)
	Get(ctx context.Context, namespace string, id string) (io.ReadCloser, error)
}
