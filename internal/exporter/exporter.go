package exporter

import (
	"fmt"
	"github.com/gelleson/packup/pkg/database"
	"github.com/gelleson/packup/pkg/upload"
	"github.com/pkg/errors"
	"io"
	"strings"
	"sync"
)

type ExportService struct {
	db         *database.Database
	uploader   map[Type]upload.Uploader
	needToSkip bool
}

func (s ExportService) findByTag(tag string) ([]Export, error) {

	exporters := make([]Export, 0)

	if tx := s.db.Conn().Find(&exporters); tx.Error != nil {
		return nil, tx.Error
	}

	return exporters, nil
}

func (s ExportService) Export(snapshotId uint, namespace, tag, name string, size uint, body io.Reader) error {

	exports, err := s.findByTag(tag)

	if err != nil {
		return err
	}

	errs := make([]error, 0)

	var wg sync.WaitGroup

	for _, export := range exports {
		go func(e Export) {

			wg.Add(1)

			defer wg.Done()

			uploader, isOk := s.uploader[e.Type]

			if !isOk && s.needToSkip {
				return
			}

			if !isOk && !s.needToSkip {
				errs = append(errs, errors.New("exporter doesn't support"))
				return
			}

			id, err := uploader.Put(namespace, name, body)

			if err != nil {
				errs = append(errs, errors.New("exporter doesn't support"))
				return
			}

			snapshot := Snapshot{
				Size:       size,
				Filename:   name,
				Namespace:  namespace,
				UploadID:   id,
				SnapshotID: snapshotId,
			}

			if tx := s.db.Conn().Create(&snapshot); tx.Error != nil {
				errs = append(errs, tx.Error)
				return
			}
		}(export)
	}

	wg.Wait()

	if errs != nil {
		return mergeErrors(errs)
	}

	return nil
}

func mergeErrors(array []error) error {

	errorString := make([]string, len(array), cap(array))

	for index, err := range array {
		errorString[index] = err.Error()
	}

	return fmt.Errorf(strings.Join(errorString, "\n"))
}
