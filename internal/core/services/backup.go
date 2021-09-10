package services

import (
	"github.com/gelleson/packup/internal/core/models"
	"github.com/gelleson/packup/pkg/database"
)

type keystoreService interface {
	Get(key string) (models.Credential, error)
}

type BackupService struct {
	db              *database.Database
	keystoreService keystoreService
}

func NewBackupService(db *database.Database, keystoreService keystoreService) *BackupService {
	return &BackupService{db: db, keystoreService: keystoreService}
}

func (bs BackupService) Create(b models.Backup) (models.Backup, error) {

	if err := b.Validate(); err != nil {
		return models.Backup{}, err
	}

	b.SetDefaults()

	if b.Keystore != "" {
		if _, err := bs.keystoreService.Get(b.Keystore); err != nil {
			return models.Backup{}, err
		}
	}

	if tx := bs.db.Conn().Create(&b); tx.Error != nil {
		return models.Backup{}, tx.Error
	}

	return b, nil
}

func (bs BackupService) Update(b models.Backup) (models.Backup, error) {

	if err := b.Validate(); err != nil {
		return models.Backup{}, err
	}

	b.SetDefaults()

	if b.Keystore != "" {
		if _, err := bs.keystoreService.Get(b.Keystore); err != nil {
			return models.Backup{}, err
		}
	}

	if tx := bs.db.Conn().Model(&b).Updates(b); tx.Error != nil {
		return models.Backup{}, tx.Error
	}

	return b, nil
}

func (bs BackupService) Find(skip uint, limit uint) ([]models.Backup, error) {

	backups := make([]models.Backup, 0)

	if tx := bs.db.Conn().Limit(int(limit)).Offset(int(skip)).Find(&backups); tx.Error != nil {
		return nil, tx.Error
	}

	return backups, nil
}

func (bs BackupService) FindById(id uint) (models.Backup, error) {

	backup := models.Backup{}

	if tx := bs.db.Conn().First(&backup, "id = ?", id); tx.Error != nil {
		return models.Backup{}, tx.Error
	}

	return backup, nil
}

func (bs BackupService) Delete(id uint) error {

	if tx := bs.db.Conn().Delete(&models.Backup{}, "id = ?", int(id)); tx.Error != nil {
		return tx.Error
	}

	return nil
}
