package backup

import (
	"github.com/gelleson/packup/pkg/database"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"io"
	"sync"
	"time"
)

type exporterService interface {
	Export(snapshotId uint, namespace, tag, name string, size uint, body io.Reader) error
}

type bucketService interface {
	FindById(id uint) (Backup, error)
}

type SnapshotService struct {
	db              *database.Database
	exporterService exporterService
	bucketService   bucketService
}

func (s SnapshotService) OK(agentId, backupId uint, objName string, size uint, body io.Reader) (Snapshot, error) {

	backup, err := s.bucketService.FindById(backupId)

	if err != nil {
		return Snapshot{}, err
	}

	snapshot := Snapshot{
		Status:     PendingStatus,
		Size:       size,
		AgentID:    agentId,
		BackupID:   backup.ID,
		Tag:        backup.Tag,
		ExecutedAt: time.Now(),
	}

	if tx := s.db.Conn().Create(&snapshot); tx.Error != nil {
		return Snapshot{}, nil
	}

	if err := s.exporterService.Export(snapshot.ID, backup.Namespace, backup.Tag, objName, size, body); err != nil {
		snapshot.Message = err.Error()
		snapshot.Status = ExportFailStatus
		snapshot.ExecutedAt = time.Now()

		if tx := s.db.Conn().Save(&snapshot); tx.Error != nil {
			return Snapshot{}, tx.Error
		}

		return snapshot, err
	}

	snapshot.Message = "ok"
	snapshot.Status = OkStatus
	snapshot.ExecutedAt = time.Now()

	if tx := s.db.Conn().Save(&snapshot); tx.Error != nil {
		return Snapshot{}, tx.Error
	}

	return snapshot, nil
}

func (s SnapshotService) Failed(agentId, backupId uint, errMessage error) (Snapshot, error) {

	if errMessage == nil {
		return Snapshot{}, errors.New("errMessage should be non nil value")
	}

	backup, err := s.bucketService.FindById(backupId)

	if err != nil {
		return Snapshot{}, err
	}

	snapshot := Snapshot{
		Status:     FailedStatus,
		AgentID:    agentId,
		BackupID:   backup.ID,
		Message:    errMessage.Error(),
		Tag:        backup.Tag,
		ExecutedAt: time.Now(),
	}

	if tx := s.db.Conn().Create(&snapshot); tx.Error != nil {
		return Snapshot{}, nil
	}

	return snapshot, nil
}

func (s SnapshotService) Find(input FindSnapshotQuery) (SnapshotWithTotal, error) {

	snapshotWithTotal := SnapshotWithTotal{}

	input.Init()

	var wg sync.WaitGroup

	var errDb error

	conn := s.db.Conn()

	conn = conn.Where(queryConstructor(conn, input))

	go func() {
		wg.Add(1)

		defer wg.Done()

		if tx := conn.Limit(int(input.Limit)).Offset(int(input.Skip)).Find(&snapshotWithTotal.Snapshots); tx.Error != nil {
			errDb = tx.Error
		}
	}()

	go func() {
		wg.Add(1)

		defer wg.Done()

		if tx := conn.Count(&snapshotWithTotal.Total); tx.Error != nil {
			errDb = tx.Error
		}
	}()

	wg.Wait()

	if errDb != nil {
		return SnapshotWithTotal{}, errDb
	}

	return snapshotWithTotal, nil
}

func queryConstructor(conn *gorm.DB, input FindSnapshotQuery) *gorm.DB {

	if input.Agent != 0 {
		conn = conn.Where("agent_id", input.Agent)
	}

	if input.Backup != 0 {
		conn = conn.Where("backup_id", input.Backup)
	}

	if input.Tag != "" {
		conn = conn.Where("tag", input.Tag)
	}

	return conn
}
