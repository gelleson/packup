package backup

import (
	"github.com/gelleson/packup/internal/keystore"
	"github.com/gelleson/packup/pkg/database"
)

type keystoreService interface {
	Get(key string) (keystore.Credential, error)
}

type BackupService struct {
	db              *database.Database
	keystoreService keystoreService
}

func NewBackupService(db *database.Database, keystoreService keystoreService) *BackupService {
	return &BackupService{db: db, keystoreService: keystoreService}
}

func (s BackupService) Create(b Backup) (Backup, error) {

	if err := b.Validate(); err != nil {
		return Backup{}, err
	}

	b.SetDefaults()

	if b.Keystore != "" {
		if _, err := s.keystoreService.Get(b.Keystore); err != nil {
			return Backup{}, err
		}
	}

	if tx := s.db.Conn().Create(&b); tx.Error != nil {
		return Backup{}, tx.Error
	}

	return b, nil
}

func (s BackupService) Update(b Backup) (Backup, error) {

	if err := b.Validate(); err != nil {
		return Backup{}, err
	}

	b.SetDefaults()

	if b.Keystore != "" {
		if _, err := s.keystoreService.Get(b.Keystore); err != nil {
			return Backup{}, err
		}
	}

	if tx := s.db.Conn().Model(&b).Updates(b); tx.Error != nil {
		return Backup{}, tx.Error
	}

	return b, nil
}

func (s BackupService) Find(skip uint, limit uint) ([]Backup, error) {

	backups := make([]Backup, 0)

	if tx := s.db.Conn().Limit(int(limit)).Offset(int(skip)).Find(&backups); tx.Error != nil {
		return nil, tx.Error
	}

	return backups, nil
}

func (s BackupService) Delete(id uint) error {

	if tx := s.db.Conn().Delete(&Backup{}, "id = ?", int(id)); tx.Error != nil {
		return tx.Error
	}

	return nil
}
