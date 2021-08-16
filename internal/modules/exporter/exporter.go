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
		wg.Add(1)

		go func(e Export) {

			defer wg.Done()

			if err := s.export(e.Type, snapshotId, namespace, name, size, body); err != nil {
				errs = append(errs, err)
			}

		}(export)
	}

	wg.Wait()

	if errs != nil {
		return mergeErrors(errs)
	}

	return nil
}

func (s ExportService) export(ext Type, snapshotId uint, namespace, name string, size uint, body io.Reader) error {

	uploader, isOk := s.uploader[ext]

	if !isOk && !s.needToSkip {
		return errors.New("doesn't support uploader")
	} else if !isOk && s.needToSkip {
		return nil
	}

	id, err := uploader.Put(namespace, name, body)

	if err != nil {
		return err
	}

	snapshot := createSnapshot(id, snapshotId, namespace, name, size)

	if tx := s.db.Conn().Create(&snapshot); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func createSnapshot(id string, snapshotId uint, namespace, name string, size uint) Snapshot {

	snapshot := Snapshot{
		Size:       size,
		Filename:   name,
		Namespace:  namespace,
		UploadID:   id,
		SnapshotID: snapshotId,
	}

	return snapshot
}

func mergeErrors(array []error) error {

	errorString := make([]string, len(array), cap(array))

	for index, err := range array {
		errorString[index] = err.Error()
	}

	return fmt.Errorf(strings.Join(errorString, "\n"))
}
