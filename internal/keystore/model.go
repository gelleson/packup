package keystore

import (
	"github.com/gelleson/packup/pkg/validators"
	"gorm.io/gorm"
)

type Credential struct {
	gorm.DB
	Key      string `validate:"required"`
	Username string
	Password string
	Host     string
	Token    string
	Database string
}

func (c Credential) Validate() error {
	return validators.Struct(c)
}
