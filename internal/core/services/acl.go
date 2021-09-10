package services

import (
	"github.com/gelleson/packup/internal/core/constants"
	"github.com/gelleson/packup/internal/core/models"
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

func (a AclService) Create(r models.Rule) (models.Rule, error) {

	if !a.groupService.Exist(r.GroupID) {
		return models.Rule{}, errors.New("group doesn't exist")
	}

	rule := models.Rule{
		Resource:  r.Resource,
		Operation: r.Operation,
		GroupID:   r.GroupID,
	}

	if tx := a.db.Conn().Create(&rule); tx.Error != nil {
		return models.Rule{}, tx.Error
	}

	return rule, nil
}

func (a AclService) Can(groupId uint, operation models.Operation, resource models.Resource) bool {

	rule := models.Rule{}

	if tx := a.db.Conn().Where("group_id = ? and operation = ? and resource = ?", groupId, operation, resource).First(&rule); tx.Error != nil {
		return false
	}

	return rule.ID != 0
}

func (a AclService) HasDefaultRules() bool {

	rule := models.Rule{}

	if tx := a.db.Conn().Where("group_id = ? ", constants.DefaultGroupId).First(&rule); tx.Error != nil {
		return false
	}

	return rule.ID != 0
}
