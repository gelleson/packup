package models

import "gorm.io/gorm"

type Operation string

const (
	WriteOps Operation = "WRITE"
	ReadOps  Operation = "READ"
)

type Resource string

const (
	GroupResource Resource = "GROUP"
	UserResource  Resource = "USER"
)

type Rule struct {
	gorm.Model
	Operation Operation
	Resource  Resource
	GroupID   uint
}
