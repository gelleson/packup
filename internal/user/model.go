package user

import (
	"github.com/gelleson/packup/internal/group"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email    string
	Password string
	IsActive bool

	GroupID uint
	Group   group.Group

	CreatedAt  time.Time
	LastLogged time.Time
}
