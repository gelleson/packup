package service

import (
	"github.com/gelleson/packup/internal/core/constants"
	"github.com/gelleson/packup/internal/core/model"
	"github.com/gelleson/packup/pkg/database"
	"github.com/pkg/errors"
)

type groupService interface {
	Exist(id uint) bool
}

type AclService struct {
	db           *database.Database
	groupService groupService
}

func NewAclService(db *database.Database, groupService groupService) *AclService {
	return &AclService{db: db, groupService: groupService}
}

func (a AclService) Create(r model.Rule) (model.Rule, error) {

	if !a.groupService.Exist(r.GroupID) {
		return model.Rule{}, errors.New("group doesn't exist")
	}

	rule := model.Rule{
		Resource:  r.Resource,
		Operation: r.Operation,
		GroupID:   r.GroupID,
	}

	if tx := a.db.Conn().Create(&rule); tx.Error != nil {
		return model.Rule{}, tx.Error
	}

	return rule, nil
}

func (a AclService) Can(groupId uint, operation model.Operation, resource model.Resource) bool {

	rule := model.Rule{}

	if tx := a.db.Conn().Where("group_id = ? and operation = ? and resource = ?", groupId, operation, resource).First(&rule); tx.Error != nil {
		return false
	}

	return rule.ID != 0
}

func (a AclService) HasDefaultRules() bool {

	rule := model.Rule{}

	if tx := a.db.Conn().Where("group_id = ? ", constants.DefaultGroupId).First(&rule); tx.Error != nil {
		return false
	}

	return rule.ID != 0
}
