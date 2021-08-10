package acl

import (
	"github.com/gelleson/packup/internal/group"
	"github.com/gelleson/packup/pkg/database"
	"github.com/pkg/errors"
)

type groupService interface {
	Exist(id uint) bool
}

type Service struct {
	db           *database.Database
	groupService groupService
}

func NewService(db *database.Database, groupService groupService) *Service {
	return &Service{db: db, groupService: groupService}
}

func (s Service) Create(r Rule) (Rule, error) {

	if !s.groupService.Exist(r.GroupID) {
		return Rule{}, errors.New("group doesn't exist")
	}

	rule := Rule{
		Resource:  r.Resource,
		Operation: r.Operation,
		GroupID:   r.GroupID,
	}

	if tx := s.db.Conn().Create(&rule); tx.Error != nil {
		return Rule{}, tx.Error
	}

	return rule, nil
}

func (s Service) Can(groupId uint, operation Operation, resource Resource) bool {

	rule := Rule{}

	if tx := s.db.Conn().Where("group_id = ? and operation = ? and resource = ?", groupId, operation, resource).First(&rule); tx.Error != nil {
		return false
	}

	return rule.ID != 0
}

func (s Service) HasDefaultRules() bool {

	rule := Rule{}

	if tx := s.db.Conn().Where("group_id = ? ", group.DefaultGroupId).First(&rule); tx.Error != nil {
		return false
	}

	return rule.ID != 0
}
