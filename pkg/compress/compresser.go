package compress

import (
	"compress/gzip"
	"io"
	"os"
)

type compressObject struct {
	os.File
}

func (c compressObject) Close() error {

	if err := c.File.Close(); err != nil {
		return err
	}

	return os.Remove(c.Name())
}

func Compress(compressType Type, body io.Reader, path string) (io.ReadCloser, error) {

	switch compressType {
	case GzipType:
		return gzipCompress(body)
	}

	return nil, nil
}

func gzipCompress(body io.Reader) (io.ReadCloser, error) {

	object, err := gzip.NewReader(body)

	if err != nil {
		return nil, err
	}

	var compressedObject compressObject

	if _, err := io.Copy(&compressedObject, object); err != nil {
		return nil, err
	}

	return &compressedObject, nil
}
