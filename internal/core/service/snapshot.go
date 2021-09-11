package service

import (
	"github.com/gelleson/packup/internal/core/dto"
	"github.com/gelleson/packup/internal/core/model"
	"github.com/gelleson/packup/pkg/database"
	"github.com/gelleson/packup/pkg/storage"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"io"
	"sync"
	"time"
)

type bucketService interface {
	FindById(id uint) (model.Backup, error)
}

type SnapshotService struct {
	db            *database.Database
	bucketService bucketService
	uploader      storage.API
}

func NewSnapshotService(db *database.Database, bucketService bucketService, uploader storage.API) *SnapshotService {
	return &SnapshotService{db: db, bucketService: bucketService, uploader: uploader}
}

func (s SnapshotService) Snap(agentId, backupId uint, objName string, size uint, body io.Reader) (model.Snapshot, error) {

	backupObject, err := s.bucketService.FindById(backupId)

	if err != nil {
		return model.Snapshot{}, err
	}

	snapshot := model.Snapshot{
		Status:     model.PendingStatus,
		Size:       size,
		AgentID:    agentId,
		BackupID:   backupObject.ID,
		Tag:        backupObject.Tag,
		ExecutedAt: time.Now(),
	}

	if tx := s.db.Conn().Create(&snapshot); tx.Error != nil {
		return model.Snapshot{}, nil
	}

	objectId, err := s.uploader.Put(backupObject.Namespace, objName, body)

	if err != nil {
		return s.FailedExistSnapshot(snapshot.ID, err)
	}

	snapshot.Message = "ok"
	snapshot.ObjectId = objectId
	snapshot.ObjectName = objName
	snapshot.Size = size
	snapshot.Status = model.OkStatus
	snapshot.ExecutedAt = time.Now()

	if tx := s.db.Conn().Save(&snapshot); tx.Error != nil {
		return model.Snapshot{}, tx.Error
	}

	return snapshot, nil
}

func (s SnapshotService) FailedExistSnapshot(snapshotId uint, errMessage error) (model.Snapshot, error) {

	if errMessage == nil {
		return model.Snapshot{}, errors.New("errMessage should be non nil value")
	}

	snapshot := model.Snapshot{}

	if tx := s.db.Conn().First(&snapshot, "id = ?", snapshotId); tx.Error != nil {
		return model.Snapshot{}, tx.Error
	}

	snapshot.Message = errMessage.Error()
	snapshot.Status = model.FailedStatus

	if tx := s.db.Conn().Save(&snapshot); tx.Error != nil {
		return model.Snapshot{}, nil
	}

	return snapshot, nil
}

func (s SnapshotService) Failed(agentId, backupId uint, errMessage error) (model.Snapshot, error) {

	if errMessage == nil {
		return model.Snapshot{}, errors.New("errMessage should be non nil value")
	}

	backup, err := s.bucketService.FindById(backupId)

	if err != nil {
		return model.Snapshot{}, err
	}

	snapshot := model.Snapshot{
		Status:     model.FailedStatus,
		AgentID:    agentId,
		BackupID:   backup.ID,
		Message:    errMessage.Error(),
		Tag:        backup.Tag,
		ExecutedAt: time.Now(),
	}

	if tx := s.db.Conn().Create(&snapshot); tx.Error != nil {
		return model.Snapshot{}, nil
	}

	return snapshot, nil
}

func (s SnapshotService) Find(input dto.FindSnapshotQuery) (dto.SnapshotWithTotal, error) {

	snapshotWithTotal := dto.SnapshotWithTotal{}

	input.Init()

	var wg sync.WaitGroup

	var errDb error

	conn := s.db.Conn()

	conn = conn.Where(queryConstructor(conn, input))

	wg.Add(2)

	go func() {

		defer wg.Done()

		if tx := conn.Limit(int(input.Limit)).Offset(int(input.Skip)).Find(&snapshotWithTotal.Snapshots); tx.Error != nil {
			errDb = tx.Error
		}
	}()

	go func() {

		defer wg.Done()

		if tx := conn.Count(&snapshotWithTotal.Total); tx.Error != nil {
			errDb = tx.Error
		}
	}()

	wg.Wait()

	if errDb != nil {
		return dto.SnapshotWithTotal{}, errDb
	}

	return snapshotWithTotal, nil
}

func queryConstructor(conn *gorm.DB, input dto.FindSnapshotQuery) *gorm.DB {

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
