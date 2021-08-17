package backup

import (
	"github.com/gelleson/packup/internal/modules/keystore"
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

func (bs BackupService) Create(b Backup) (Backup, error) {

	if err := b.Validate(); err != nil {
		return Backup{}, err
	}

	b.SetDefaults()

	if b.Keystore != "" {
		if _, err := bs.keystoreService.Get(b.Keystore); err != nil {
			return Backup{}, err
		}
	}

	if tx := bs.db.Conn().Create(&b); tx.Error != nil {
		return Backup{}, tx.Error
	}

	return b, nil
}

func (bs BackupService) Update(b Backup) (Backup, error) {

	if err := b.Validate(); err != nil {
		return Backup{}, err
	}

	b.SetDefaults()

	if b.Keystore != "" {
		if _, err := bs.keystoreService.Get(b.Keystore); err != nil {
			return Backup{}, err
		}
	}

	if tx := bs.db.Conn().Model(&b).Updates(b); tx.Error != nil {
		return Backup{}, tx.Error
	}

	return b, nil
}

func (bs BackupService) Find(skip uint, limit uint) ([]Backup, error) {

	backups := make([]Backup, 0)

	if tx := bs.db.Conn().Limit(int(limit)).Offset(int(skip)).Find(&backups); tx.Error != nil {
		return nil, tx.Error
	}

	return backups, nil
}

func (bs BackupService) FindById(id uint) (Backup, error) {

	backup := Backup{}

	if tx := bs.db.Conn().First(&backup, "id = ?", id); tx.Error != nil {
		return Backup{}, tx.Error
	}

	return backup, nil
}

func (bs BackupService) Delete(id uint) error {

	if tx := bs.db.Conn().Delete(&Backup{}, "id = ?", int(id)); tx.Error != nil {
		return tx.Error
	}

	return nil
}
