package exporter

import (
	"context"
	"github.com/gelleson/packup/pkg/database"
	"github.com/gelleson/packup/pkg/helpers"
	"github.com/gelleson/packup/pkg/upload"
	"github.com/pkg/errors"
	"io"
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

func (s ExportService) Export(ctx context.Context, snapshotId uint, namespace, tag, name string, size uint, body io.Reader) error {

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

			if err := s.export(ctx, e.Type, snapshotId, namespace, name, size, body); err != nil {
				errs = append(errs, err)
			}

		}(export)
	}

	wg.Wait()

	if errs != nil {
		return helpers.MergeErrors(errs)
	}

	return nil
}

func (s ExportService) export(ctx context.Context, ext Type, snapshotId uint, namespace, name string, size uint, body io.Reader) error {

	uploader, isOk := s.uploader[ext]

	if !isOk && !s.needToSkip {
		return errors.New("doesn't support uploader")
	} else if !isOk && s.needToSkip {
		return nil
	}

	id, err := uploader.Put(ctx, namespace, name, body)

	if err != nil {
		return err
	}

	objectDTO := prepareObject(id, snapshotId, namespace, name, size)

	if tx := s.db.Conn().Create(&objectDTO); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func prepareObject(id string, snapshotId uint, namespace, name string, size uint) Object {

	object := Object{
		Size:       size,
		Filename:   name,
		StorageID:  id,
		SnapshotID: snapshotId,
	}

	return object
}
