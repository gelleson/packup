package exporter

import (
	"gorm.io/gorm"
	"time"
)

type Type string

const (
	S3Type Type = "s3"
	FSType Type = "fs"
)

type Export struct {
	gorm.Model
	Name string `validate:"required"`
	Key  string `validate:"required"`
	Type Type   `validate:"required,oneof=s3 fs"`
	Tag  string `validate:"required"`
}

type Object struct {
	gorm.Model
	Size         uint
	Bucket       string
	Filename     string
	StorageID    string
	Type         Type
	Downloadable bool
	ExportID     uint
	Export       Export
	SnapshotID   uint
	UploadedAt   time.Time
}
