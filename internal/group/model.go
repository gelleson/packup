package group

import (
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Name string
}

type Permission struct {
	gorm.Model
	Operation string
	GroupID   uint
	Group     Group
}
