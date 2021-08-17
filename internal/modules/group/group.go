package group

import (
	"github.com/gelleson/packup/pkg/database"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type GroupService struct {
	db *database.Database

	createdDefaultDatabase bool
}

func NewGroupService(db *database.Database) *GroupService {
	return &GroupService{db: db}
}

func (s GroupService) Create(input CreateGroupInput) (Group, error) {

	if !s.HasDefaultGroup() {
		return Group{}, errors.New("first need to create default group")
	}

	if err := input.Validate(); err != nil {
		return Group{}, err
	}

	group := Group{
		Name: input.Name,
	}

	if trx := s.db.Conn().Create(&group); trx.Error != nil {
		return Group{}, trx.Error
	}

	return group, nil
}

func (s GroupService) CreateDefaultGroup() (Group, error) {

	group := Group{
		Model: gorm.Model{
			ID: DefaultGroupId,
		},
		Name: "default",
	}

	if trx := s.db.Conn().Create(&group); trx.Error != nil {
		return Group{}, trx.Error
	}

	return group, nil
}

func (s GroupService) FindById(id uint) (Group, error) {

	if !s.HasDefaultGroup() {
		return Group{}, errors.New("first need to create default group")
	}

	group := Group{}

	if trx := s.db.Conn().First(&group, "id = ?", id); trx.Error != nil {
		return Group{}, trx.Error
	}

	return group, nil
}

func (s GroupService) Exist(id uint) bool {

	_, err := s.FindById(id)

	return err == nil
}

func (s *GroupService) HasDefaultGroup() bool {

	if s.createdDefaultDatabase {
		return true
	}

	group := Group{}

	if trx := s.db.Conn().First(&group, "id = ?", DefaultGroupId); trx.Error != nil {
		return false
	}

	s.createdDefaultDatabase = true

	return true
}
