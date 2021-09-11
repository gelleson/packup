package storage

import "io"

type API interface {
	Put(namespace string, filename string, body io.Reader) (string, error)
	Get(namespace string, id string) (io.ReadCloser, error)
}
