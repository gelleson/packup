package service

import (
	"github.com/gelleson/packup/internal/core/model"
	"github.com/gelleson/packup/pkg/database"
)

type BackupService struct {
	db *database.Database
}

func NewBackupService(db *database.Database) *BackupService {
	return &BackupService{db: db}
}

func (bs BackupService) Create(b model.Backup) (model.Backup, error) {

	if err := b.Validate(); err != nil {
		return model.Backup{}, err
	}

	b.SetDefaults()

	if tx := bs.db.Conn().Create(&b); tx.Error != nil {
		return model.Backup{}, tx.Error
	}

	return b, nil
}

func (bs BackupService) Update(b model.Backup) (model.Backup, error) {

	if err := b.Validate(); err != nil {
		return model.Backup{}, err
	}

	b.SetDefaults()

	if tx := bs.db.Conn().Model(&model.Backup{}).Updates(b); tx.Error != nil {
		return model.Backup{}, tx.Error
	}

	return b, nil
}

func (bs BackupService) Find(skip uint, limit uint) ([]model.Backup, error) {

	backups := make([]model.Backup, 0)

	if tx := bs.db.Conn().Limit(int(limit)).Offset(int(skip)).Find(&backups); tx.Error != nil {
		return nil, tx.Error
	}

	return backups, nil
}

func (bs BackupService) FindById(id uint) (model.Backup, error) {

	backup := model.Backup{}

	if tx := bs.db.Conn().First(&backup, "id = ?", id); tx.Error != nil {
		return model.Backup{}, tx.Error
	}

	return backup, nil
}

func (bs BackupService) Delete(id uint) error {

	if tx := bs.db.Conn().Delete(&model.Backup{}, "id = ?", int(id)); tx.Error != nil {
		return tx.Error
	}

	return nil
}
