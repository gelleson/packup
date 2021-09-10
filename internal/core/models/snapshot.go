package models

import (
	"gorm.io/gorm"
	"time"
)

type Snapshot struct {
	gorm.Model
	Status     Status
	Message    string
	Size       uint
	ObjectName string
	ObjectId   string
	Tag        string
	AgentID    uint
	BackupID   uint
	Backup     Backup
	ExecutedAt time.Time
}

type PendingSnapshot struct {
	gorm.Model
	BackupID      uint
	Backup        Backup
	ExecutionTime time.Time
}
