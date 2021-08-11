package backup

import (
	"github.com/gelleson/packup/pkg/compress"
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

	if et != OnceExecution && et != CronExecution {
		return errors.New("execution type should be equal to cron or once")
	}

	return nil
}

const (
	OnceExecution ExecutionType = "once"
	CronExecution ExecutionType = "cron"
)

type Template string

const DefaultBackupTemplate Template = "{{ .Namespace }}/{{ .Name }}-{{ .Timestamp }}.{{ .Ext }}"

type Backup struct {
	gorm.Model
	Name               string        `validate:"required"`
	Compress           compress.Type `validate:"required"`
	ExecutionType      ExecutionType `validate:"required"`
	ExecutionTime      string        `validate:"required"`
	Tag                string        `validate:"required"`
	Keystore           string
	Bucket             string `validate:"required"`
	BackupNameTemplate Template
	Namespace          string
	Timezone           string
	LastExecutionTime  time.Time
}

type Status string

const (
	FailedStatus Status = "failed"
	OkStatus     Status = "ok"
)

type History struct {
	gorm.Model
	Status     Status
	Message    string
	Size       uint
	Tag        string
	AgentID    uint
	BackupID   uint
	ExecutedAt time.Time
}

type Pending struct {
	gorm.Model
	BackupID      uint
	Backup        Backup
	ExecutionTime time.Time
}
