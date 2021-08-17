package acl

import (
	"github.com/gelleson/packup/internal/modules/group"
	"github.com/gelleson/packup/pkg/database"
	"github.com/pkg/errors"
)

type groupService interface {
	Exist(id uint) bool
}

type ACL struct {
	db           *database.Database
	groupService groupService
}

func NewAcl(db *database.Database, groupService groupService) *ACL {
	return &ACL{db: db, groupService: groupService}
}

func (a ACL) Create(r Rule) (Rule, error) {

	if !a.groupService.Exist(r.GroupID) {
		return Rule{}, errors.New("group doesn't exist")
	}

	rule := Rule{
		Resource:  r.Resource,
		Operation: r.Operation,
		GroupID:   r.GroupID,
	}

	if tx := a.db.Conn().Create(&rule); tx.Error != nil {
		return Rule{}, tx.Error
	}

	return rule, nil
}

func (a ACL) Can(groupId uint, operation Operation, resource Resource) bool {

	rule := Rule{}

	if tx := a.db.Conn().Where("group_id = ? and operation = ? and resource = ?", groupId, operation, resource).First(&rule); tx.Error != nil {
		return false
	}

	return rule.ID != 0
}

func (a ACL) HasDefaultRules() bool {

	rule := Rule{}

	if tx := a.db.Conn().Where("group_id = ? ", group.DefaultGroupId).First(&rule); tx.Error != nil {
		return false
	}

	return rule.ID != 0
}
