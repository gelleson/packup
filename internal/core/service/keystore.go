package service

import (
	"github.com/gelleson/packup/internal/core/model"
	"github.com/gelleson/packup/pkg/cipher"
	"github.com/gelleson/packup/pkg/database"
	"github.com/pkg/errors"
)

type Cipher string

type KeystoreService struct {
	db        *database.Database
	cipherKey string
}

func NewKeystoreService(db *database.Database, cipherKey string) *KeystoreService {
	return &KeystoreService{db: db, cipherKey: cipherKey}
}

func (s KeystoreService) Get(key string) (model.Credential, error) {

	cred := model.Credential{}

	if tx := s.db.Conn().Where("key = ?", key).First(&cred); tx.Error != nil {
		return model.Credential{}, tx.Error
	}

	return cred, nil
}

func (s KeystoreService) Create(c model.Credential) (model.Credential, error) {

	if err := c.Validate(); err != nil {
		return model.Credential{}, err
	}

	encryptedCredential, err := s.encrypt(c)

	if err != nil {
		return model.Credential{}, err
	}

	if tx := s.db.Conn().Create(&encryptedCredential); tx.Error != nil {
		return model.Credential{}, tx.Error
	}

	return encryptedCredential, nil
}

func (s KeystoreService) Delete(key string) error {

	if tx := s.db.Conn().Delete(&model.Credential{}, "key = ?", key); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s KeystoreService) checkCipherKey() error {

	testValue := "test"

	encrypted, err := cipher.EncryptString(testValue, s.cipherKey)

	if err != nil {
		return err
	}

	decrypt, err := cipher.DecryptString(encrypted, s.cipherKey)

	if err != nil {
		return err
	}

	if decrypt != testValue {
		return errors.New("cipher is doesn't work")
	}

	return nil
}

func (s KeystoreService) encryptOrGetEmpty(data string) (string, error) {

	if data == "" {
		return "", nil
	}

	encrypted, err := cipher.EncryptString(data, s.cipherKey)

	if err != nil {
		return "", err
	}

	return encrypted, nil
}

func (s KeystoreService) encrypt(c model.Credential) (model.Credential, error) {

	if err := s.checkCipherKey(); err != nil {
		return model.Credential{}, err
	}

	c.Username, _ = s.encryptOrGetEmpty(c.Username)
	c.Password, _ = s.encryptOrGetEmpty(c.Password)
	c.Host, _ = s.encryptOrGetEmpty(c.Host)
	c.Token, _ = s.encryptOrGetEmpty(c.Token)
	c.Database, _ = s.encryptOrGetEmpty(c.Database)
	c.KeyId, _ = s.encryptOrGetEmpty(c.KeyId)
	c.Secret, _ = s.encryptOrGetEmpty(c.Secret)

	return c, nil
}
