package model

import (
	"github.com/gelleson/packup/pkg/compress"
	"github.com/gelleson/packup/pkg/validators"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type SourceType string

const (
	PostgresType SourceType = "postgres"
	MySQLType    SourceType = "mysql"
	OracleType   SourceType = "oracle"
	SqliteType   SourceType = "sqlite"
	FileType     SourceType = "file"
)

type ExecutionType string

func (et ExecutionType) Validate() error {

	if et != OnceExecution && et != CronExecution && et != RruleExecution {
		return errors.New("execution type should be equal to cron or once or rrule")
	}

	return nil
}

const (
	OnceExecution  ExecutionType = "once"
	CronExecution  ExecutionType = "cron"
	RruleExecution ExecutionType = "rrule"
)

type Template string

const DefaultBackupTemplate Template = "{{ .Name }}-{{ .Timestamp }}.{{ .Ext }}"

type Backup struct {
	gorm.Model
	Name               string        `validate:"required"`
	Compress           compress.Type `validate:"required"`
	Tag                string        `validate:"required"`
	Keystore           string        `validate:"required"`
	Namespace          string        `validate:"required"`
	BackupNameTemplate Template
	Timezone           string
	LastExecutionTime  time.Time

	ExecutionType ExecutionType `validate:"required,oneof=once cron rrule"`
	Cron          string        `validate:"required_if=ExecutionType cron"`
	Rrule         string        `validate:"required_if=ExecutionType rrule"`
	ExecutionTime time.Time     `validate:"required_if=ExecutionType once"`
}

func (b Backup) Validate() error {

	if err := validators.Struct(b); err != nil {
		return err
	}

	if err := validators.IsValidExecutionValue(b); err != nil {
		return err
	}

	if err := validators.IsValidTimezone(b.Timezone); err != nil {
		return err
	}

	return nil
}

func (b *Backup) SetDefaults() {

	if b.BackupNameTemplate == "" {
		b.BackupNameTemplate = DefaultBackupTemplate
	}

	if b.Timezone == "" {
		zone, _ := time.Now().Local().Zone()

		b.Timezone = zone
	}

	if b.Namespace == "" {
		b.Namespace = b.Name
	}
}

type Status string

const (
	FailedStatus     Status = "failed"
	OkStatus         Status = "ok"
	PendingStatus    Status = "pending"
	ExportFailStatus Status = "export_failed"
)
