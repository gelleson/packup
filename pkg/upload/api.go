package upload

import "io"

type Uploader interface {
	Put(namespace string, filename string, body io.Reader) (string, error)
	Get(namespace string, id string) (io.ReadCloser, error)
}
