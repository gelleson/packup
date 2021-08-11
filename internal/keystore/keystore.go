package keystore

import "github.com/gelleson/packup/pkg/database"

type Cipher string

type KeystoreService struct {
	db        *database.Database
	cipherKey string
}

func (s KeystoreService) Get(key string) (Credential, error) {

	cred := Credential{}

	if tx := s.db.Conn().Where("key = ?", key).First(&cred); tx.Error != nil {
		return Credential{}, tx.Error
	}

	return cred, nil
}

func (s KeystoreService) encrypt(c Credential) (Credential, error) {

	return Credential{}, nil
}
