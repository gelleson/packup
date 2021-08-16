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
	Name     string `validate:"required"`
	Keystore string `validate:"required"`
	Type     Type   `validate:"required,oneof=s3 fs"`
	Tag      string `validate:"required"`

	S3Endpoint string `validate:"required_if=Type s3" json:"s3_endpoint"`
	S3Bucket   string `validate:"required_if=Type s3" json:"s3_bucket"`

	FSPath string `validate:"required_if=Type fs" json:"fs_path"`
}

type Snapshot struct {
	gorm.Model
	Size       uint
	Filename   string
	Namespace  string
	UploadID   string
	SnapshotID uint
	UploadedAt time.Time
}
