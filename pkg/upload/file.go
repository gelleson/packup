package upload

import (
	"context"
	"io"
	"os"
	"path"
)

type FileProvider struct {
	BaseFolder string
}

func (f FileProvider) Put(ctx context.Context, namespace string, filename string, body io.Reader) (string, error) {

	if _, err := os.Stat(f.BaseFolder); os.IsNotExist(err) {
		return "", err
	}

	folder := path.Join(f.BaseFolder, namespace)

	if err := f.existPath(folder); err != nil {
		if err := os.MkdirAll(folder, os.ModePerm); err != nil {
			return "", err
		}
		return "", err
	}

	fname := path.Join(folder, filename)

	file, err := os.Create(fname)

	if err != nil {
		return "", err
	}

	if _, err = io.Copy(file, body); err != nil {
		return "", err
	}

	return filename, nil
}

func (f FileProvider) Get(ctx context.Context, namespace string, id string) (io.ReadCloser, error) {

	if _, err := os.Stat(f.BaseFolder); os.IsNotExist(err) {
		return nil, err
	}

	fPath := path.Join(f.BaseFolder, namespace, id)

	if err := f.existPath(fPath); err != nil {
		return nil, err
	}

	file, err := os.Open(fPath)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f FileProvider) existPath(p string) error {

	if _, err := os.Stat(p); os.IsNotExist(err) {
		return err
	}

	return nil
}
