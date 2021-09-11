package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email    string
	Password string
	IsActive bool

	GroupID uint
	Group   Group

	CreatedAt  time.Time
	LastLogged time.Time
}
