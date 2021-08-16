package group

import (
	"github.com/gelleson/packup/pkg/database"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Service struct {
	db *database.Database

	createdDefaultDatabase bool
}

func NewService(db *database.Database) *Service {
	return &Service{db: db}
}

func (s Service) Create(input CreateGroupInput) (Group, error) {

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

func (s Service) CreateDefaultGroup() (Group, error) {

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

func (s Service) FindById(id uint) (Group, error) {

	if !s.HasDefaultGroup() {
		return Group{}, errors.New("first need to create default group")
	}

	group := Group{}

	if trx := s.db.Conn().First(&group, "id = ?", id); trx.Error != nil {
		return Group{}, trx.Error
	}

	return group, nil
}

func (s Service) Exist(id uint) bool {

	_, err := s.FindById(id)

	return err == nil
}

func (s *Service) HasDefaultGroup() bool {

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
